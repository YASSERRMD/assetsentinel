package websocket

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	orgID := uint(c.GetFloat64("organization_id"))
	userID := uint(c.GetFloat64("user_id"))

	if orgID == 0 {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(h, conn, orgID, userID)
	h.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

type Hub struct {
	clients    map[uint]map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	OrgID   uint
	Message map[string]interface{}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		broadcast:  make(chan BroadcastMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.orgID] == nil {
				h.clients[client.orgID] = make(map[*Client]bool)
			}
			h.clients[client.orgID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.orgID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.orgID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.OrgID]
			h.mu.RUnlock()

			data, _ := json.Marshal(message.Message)
			for client := range clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastToOrg(orgID uint, message map[string]interface{}) {
	h.broadcast <- BroadcastMessage{OrgID: orgID, Message: message}
}

type Client struct {
	hub  *Hub
	conn interface {
		ReadMessage() (messageType int, p []byte, err error)
		WriteMessage(messageType int, p []byte) error
		Close() error
	}
	send   chan []byte
	orgID  uint
	userID uint
}

type ClientConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, p []byte) error
	Close() error
}

func NewClient(hub *Hub, conn ClientConn, orgID, userID uint) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		orgID:  orgID,
		userID: userID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(1, []byte{})
			break
		}
		c.conn.WriteMessage(1, message)
	}
}
