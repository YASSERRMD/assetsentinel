package handlers

import (
	"net/http"
	"strconv"

	"assetsentinel/internal/middleware"
	"assetsentinel/internal/repository"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService interface {
		Create(part *repository.InventoryPart) error
		Get(id, orgID uint) (*repository.InventoryPart, error)
		List(orgID uint, page, pageSize int) ([]repository.InventoryPart, int, error)
		Update(part *repository.InventoryPart) error
		Delete(id, orgID uint) error
	}
}

func NewInventoryHandler(inventoryService interface {
	Create(part *repository.InventoryPart) error
	Get(id, orgID uint) (*repository.InventoryPart, error)
	List(orgID uint, page, pageSize int) ([]repository.InventoryPart, int, error)
	Update(part *repository.InventoryPart) error
	Delete(id, orgID uint) error
}) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) Create(c *gin.Context) {
	var part repository.InventoryPart
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	part.OrganizationID = middleware.GetOrganizationID(c)
	if err := h.inventoryService.Create(&part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, part)
}

func (h *InventoryHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)
	part, err := h.inventoryService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory part not found"})
		return
	}
	c.JSON(http.StatusOK, part)
}

func (h *InventoryHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	parts, total, err := h.inventoryService.List(orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": parts, "total": total, "page": page, "page_size": pageSize})
}

func (h *InventoryHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)
	var part repository.InventoryPart
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	part.ID = uint(id)
	part.OrganizationID = orgID
	if err := h.inventoryService.Update(&part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, part)
}

func (h *InventoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)
	if err := h.inventoryService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inventory part deleted"})
}

type WorkOrderHandler struct {
	workOrderService interface {
		Create(wo *repository.WorkOrder) error
		Get(id, orgID uint) (*repository.WorkOrder, error)
		List(orgID uint, page, pageSize int, status string) ([]repository.WorkOrder, int, error)
		Update(wo *repository.WorkOrder) error
		Delete(id, orgID uint) error
	}
}

func NewWorkOrderHandler(workOrderService interface {
	Create(wo *repository.WorkOrder) error
	Get(id, orgID uint) (*repository.WorkOrder, error)
	List(orgID uint, page, pageSize int, status string) ([]repository.WorkOrder, int, error)
	Update(wo *repository.WorkOrder) error
	Delete(id, orgID uint) error
}) *WorkOrderHandler {
	return &WorkOrderHandler{workOrderService: workOrderService}
}

func (h *WorkOrderHandler) Create(c *gin.Context) {
	var wo repository.WorkOrder
	if err := c.ShouldBindJSON(&wo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wo.OrganizationID = middleware.GetOrganizationID(c)
	if err := h.workOrderService.Create(&wo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, wo)
}

func (h *WorkOrderHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)
	wo, err := h.workOrderService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
		return
	}
	c.JSON(http.StatusOK, wo)
}

func (h *WorkOrderHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	orders, total, err := h.workOrderService.List(orgID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders, "total": total, "page": page, "page_size": pageSize})
}

func (h *WorkOrderHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)
	var wo repository.WorkOrder
	if err := c.ShouldBindJSON(&wo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wo.ID = uint(id)
	wo.OrganizationID = orgID
	if err := h.workOrderService.Update(&wo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wo)
}

func (h *WorkOrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)
	if err := h.workOrderService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Work order deleted"})
}
