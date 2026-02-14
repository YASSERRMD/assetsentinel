import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'

const api = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' }
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const auth = {
  login: (email, password) => api.post('/auth/login', { email, password }),
  register: (data) => api.post('/auth/register', data)
}

export const assets = {
  list: (params) => api.get('/assets', { params }),
  get: (id) => api.get(`/assets/${id}`),
  create: (data) => api.post('/assets', data),
  update: (id, data) => api.put(`/assets/${id}`, data),
  delete: (id) => api.delete(`/assets/${id}`)
}

export const maintenance = {
  list: (params) => api.get('/maintenance-plans', { params }),
  get: (id) => api.get(`/maintenance-plans/${id}`),
  create: (data) => api.post('/maintenance-plans', data),
  update: (id, data) => api.put(`/maintenance-plans/${id}`, data),
  delete: (id) => api.delete(`/maintenance-plans/${id}`)
}

export const workOrders = {
  list: (params) => api.get('/work-orders', { params }),
  get: (id) => api.get(`/work-orders/${id}`),
  create: (data) => api.post('/work-orders', data),
  update: (id, data) => api.put(`/work-orders/${id}`, data),
  delete: (id) => api.delete(`/work-orders/${id}`)
}

export const inventory = {
  list: (params) => api.get('/inventory', { params }),
  get: (id) => api.get(`/inventory/${id}`),
  create: (data) => api.post('/inventory', data),
  update: (id, data) => api.put(`/inventory/${id}`, data),
  delete: (id) => api.delete(`/inventory/${id}`)
}

export const reports = {
  depreciation: (assetId) => api.get(`/reports/depreciation/${assetId}`),
  costs: (assetId) => api.get(`/reports/costs/${assetId}`),
  costsAll: () => api.get('/reports/costs')
}

export const dashboard = {
  stats: () => api.get('/dashboard')
}

export const audit = {
  list: (params) => api.get('/audit', { params })
}

class WebSocketService {
  constructor() {
    this.ws = null
    this.listeners = {}
  }

  connect() {
    const token = localStorage.getItem('token')
    this.ws = new WebSocket(`${WS_URL}?token=${token}`)
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      const type = data.type
      if (this.listeners[type]) {
        this.listeners[type].forEach(cb => cb(data))
      }
    }

    this.ws.onclose = () => setTimeout(() => this.connect(), 5000)
  }

  on(event, callback) {
    if (!this.listeners[event]) this.listeners[event] = []
    this.listeners[event].push(callback)
  }

  off(event, callback) {
    if (this.listeners[event]) {
      this.listeners[event] = this.listeners[event].filter(cb => cb !== callback)
    }
  }
}

export const ws = new WebSocketService()
export default api
