# ChatSphere V1

> A high-performance, production-ready real-time collaborative messaging platform built with Go, React, and WebSockets.

[![Go Report Card](https://goreportcard.com/badge/github.com/MasArik09/chat-sphere)](https://goreportcard.com/report/github.com/MasArik09/chat-sphere)
[![License: MIT](https://img.shields.99/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![ChatSphere CI](https://github.com/MasArik09/chat-sphere/actions/workflows/ci.yml/badge.svg)](https://github.com/MasArik09/chat-sphere/actions/workflows/ci.yml)

ChatSphere is a lightweight, secure chat platform designed for seamless real-time communication. It features a compiled Go API backend, a reactive TypeScript React frontend, persistent PostgreSQL storage, and Nginx gateway reverse-proxying with built-in security rate-limiting.

---

## 🚀 Key Features

* **Real-time Bidirectional Messaging**: Single-connection WebSocket event bus for instant message propagation.
* **Typing Indicators**: Live typing state start/stop broadcasts (`typing.start` / `typing.stop`).
* **Online Presence UI**: Real-time online/offline presence tracking and partner last-seen timestamps.
* **Read Receipts**: Visual receipt sync showing user last-read message tracking.
* **Unread Message Counters**: Dynamically calculated counters isolated per user conversation.
* **Optimized Conversation Search**: Fast database-level filters with composite indexes preventing N+1 queries.
* **DevOps & Container Orchestration**: Multi-container setups utilizing Docker Compose, split health checks (`/health/live`, `/health/ready`), and JSON log file rotation.
* **Gateway Rate Limiting**: Built-in Nginx limit zones protecting register and login routes against brute force attacks.
* **Transaction Rollback Hardening**: Database Transaction Manager wrapping updates, guaranteeing no partial writes on participant or message creation failures.

---

## 🏗️ System Architecture

ChatSphere uses a reverse-proxy gateway architecture to securely routing HTTP/WS traffic to the client application and API services.

```mermaid
graph TD
    Client[Web Browser Client] -->|HTTP / WebSockets| Nginx[Nginx Gateway Proxy <br/>Port 80 / 443]
    subgraph Docker Network [chatsphere-prod-network]
        Nginx -->|Serves Static Files| FE[Frontend Static SPA <br/>React + TypeScript]
        Nginx -->|Proxies /api/v1/| BE[Go API Server <br/>Port 8080]
        Nginx -->|Proxies /ws| BE
        BE -->|ACID SQL Queries| DB[(PostgreSQL <br/>Port 5432)]
    end
```

---

## 🛠️ Technology Stack

### Backend Services
- **Language**: Go (Golang) v1.22+
- **HTTP Web Framework**: Gin Gonic (Release Mode)
- **Database Driver**: `lib/pq` (native PostgreSQL driver)
- **WebSockets**: Standard Gorilla WebSockets implementation

### Frontend Client
- **Framework**: React 19 + TypeScript
- **Bundler**: Vite
- **State Management**: Zustand
- **Styling**: Tailwind CSS
- **Routing**: React Router DOM v7

### DevOps & Storage
- **Database**: PostgreSQL 16
- **Reverse Proxy**: Nginx (Alpine)
- **Containerization**: Docker & Docker Compose
- **CI Pipeline**: GitHub Actions

---

## ⚙️ Environment Variables Reference

Create a `.env` file for development or a `.env.production` file for production. Here are the required parameters:

| Variable | Description | Default | Environment |
| :--- | :--- | :--- | :--- |
| `DB_HOST` | Database host container | `postgres` | Dev / Prod |
| `DB_PORT` | Database connection port | `5432` | Dev / Prod |
| `DB_USER` | Database administrator username | `postgres` | Dev / Prod |
| `DB_PASSWORD` | Database user password | `postgres` | Dev / Prod |
| `DB_NAME` | Database schema name | `chatsphere` | Dev / Prod |
| `JWT_SECRET` | Secret key used to sign tokens | `super_secret_jwt_key` | Dev / Prod |
| `PORT` | Go backend service listening port | `8080` | Dev / Prod |
| `VITE_API_URL` | Frontend client target REST API | `http://localhost:8080/api/v1` | Build Time |
| `VITE_WS_URL` | Frontend client target WebSocket | `ws://localhost:8080/ws` | Build Time |


---

## 📦 Local Quickstart (Docker Development Setup)

To spin up the entire application stack in development mode:
```bash
# Clone the repository
git clone https://github.com/MasArik09/chat-sphere.git
cd chat-sphere

# Boot the containers
docker compose up -d --build
```
Once booted:
* **Frontend SPA**: Access at `http://localhost:5173`
* **Go Backend API**: Access at `http://localhost:8080/health/live`
* **PostgreSQL DB**: Bound to port `5432`

To shut down:
```bash
docker compose down
```

---

## 🛠️ Local Manual Development Setup

If you prefer to run services locally without Docker:

### 1. Database Configuration
Run a local PostgreSQL server and create a database named `chatsphere`. Run migrations (see [Database Migrations](#-database-migrations)).

### 2. Backend Setup
```bash
cd backend
# Create configuration env file
cp .env.example .env
# Start the Go server
go run cmd/api/main.go
```

### 3. Frontend Setup
```bash
cd frontend
# Create build env file
cp .env.example .env
# Install dependencies and start development server
npm install
npm run dev
```

---

## 🚢 Production Deployment

For production environments, use the optimized production Compose file:

1. Copy the production environment example file:
   ```bash
   cp .env.production.example .env.production
   ```
2. Open `.env.production` and configure secure credentials and production endpoints (e.g. `https://yourdomain.com`).
3. Deploy the production stack:
   ```bash
   docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build
   ```

*Note: In production Compose, Nginx exposes port 80 to the public, while ports 8080 (Go API) and 5432 (Postgres) remain closed and protected internally. Refer to the [Deployment Guide](docs/DEPLOYMENT_GUIDE.md) to set up SSL/TLS using host-level Certbot.*

---

## 🗄️ Database Migrations

Database schemas are managed using SQL files. You can apply migrations to the production container easily using the official `migrate/migrate` Docker image without having to install any host-level binary:

```bash
docker run --rm -v "$(pwd)/backend/migrations:/migrations" \
  --network chat-sphere_chatsphere-prod-network \
  migrate/migrate -path=/migrations \
  -database "postgres://DB_USER:DB_PASSWORD@postgres:5432/DB_NAME?sslmode=disable" up
```
*(Replace `DB_USER`, `DB_PASSWORD`, and `DB_NAME` with your production variables).*

---

## 🧪 Testing Suites

Run backend tests inside a container network to bypass host permissions and ensure clean sandbox executions:

```bash
# Run Unit Tests (Mocked DB layers)
docker run --rm -v "$(pwd)/backend:/app" -w /app golang:alpine go test -v ./internal/...

# Run Integration Tests (Requires database container network)
docker run --rm --network chat-sphere_chatsphere-network \
  -e DB_HOST=chatsphere-postgres -v "$(pwd)/backend:/app" -w /app \
  golang:alpine go test -v ./tests/...
```

---

## 📂 Project Structure

```text
├── .github/workflows/   # CI GitHub Actions pipelines
├── backend/             # Go REST & WebSocket API codebase
│   ├── cmd/api/         # Application entry point (main.go)
│   ├── internal/        # Domain core packages (auth, conversations, database, etc.)
│   ├── migrations/      # Versioned DB SQL schema files
│   └── tests/           # Integration tests
├── docs/                # Architecture, reference guides, and release notes
├── frontend/            # React + TypeScript Vite frontend
│   ├── src/             # Component, hooks, store, and service layers
│   └── Dockerfile.prod  # Nginx multi-stage build Dockerfile
├── scripts/             # Production database backup and restore scripts
├── docker-compose.yml   # Local development compose
└── docker-compose.prod.yml # Production gateway compose
```

---

## 🗺️ Roadmap

- [ ] **Group Messaging**: Extend single chat conversations to support multiple group participants.
- [ ] **Media Attachments**: Support image, video, and file transfers over WebSocket payloads.
- [ ] **Push Notifications**: Integrate web push notifications for offline users.
- [ ] **Message Editing/Deletion**: Implement soft deletion and editing tracking.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
