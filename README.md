# AssetSentinel

A complete full-stack facility and infrastructure lifecycle management system built with **opencode** and **Ollama minimax-m2:cloud**.

## Tech Stack

- **Frontend**: Vue 3 + Vite
- **Backend**: Go (Gin framework)
- **Database**: SQLite (Raw SQL - No ORM)
- **Real-time**: WebSocket (Gorilla)
- **Containerization**: Docker

## Features

- Multi-tenant authentication with JWT
- Role-based access control (Admin, MaintenanceManager, Technician, Viewer)
- Asset management with soft delete
- Preventive maintenance scheduling
- Work orders with state machine
- Spare parts inventory with transaction-safe operations
- Depreciation & cost tracking
- Real-time notifications via WebSocket

## Quick Start

### Using Docker

```bash
docker-compose up
```

### Manual Setup

**Backend:**
```bash
cd backend
go run ./cmd/server
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Default Credentials

- Email: admin@acme.com
- Password: password123

## API Endpoints

- `POST /api/auth/login` - Login
- `GET /api/assets` - List assets
- `GET /api/maintenance-plans` - List maintenance plans
- `GET /api/work-orders` - List work orders
- `GET /api/inventory` - List inventory
- `GET /api/reports/costs` - Cost reports

---

Built with **opencode** and **Ollama minimax-m2:cloud** 🤖
