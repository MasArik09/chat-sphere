# Architecture Rules

# ChatSphere V1

Version: 1.0

Status: Mandatory

---

# 1. Purpose

This document defines architecture constraints.

Every implementation must follow these rules.

Maintainability is more important than speed.

---

# 2. Architecture Overview

System Architecture:

Frontend
↓
REST API
↓
Go Backend
↓
PostgreSQL

Realtime Channel:

Frontend
↔ WebSocket
↔ Go Backend

Container Layer:

Docker Compose

---

# 3. Monolith Rule

ChatSphere V1 must be a modular monolith.

Do NOT implement:

* Microservices
* Service Mesh
* Event Bus
* Kafka

V1 complexity must remain manageable.

---

# 4. Repository Structure

Root:

chat-sphere/
│
├── frontend/
├── backend/
├── docs/
├── docker-compose.yml
├── README.md
└── LICENSE

---

# 5. Frontend Architecture

Frontend:

React
TypeScript
Vite

Structure:

frontend/
│
├── src/
│   ├── api/
│   ├── components/
│   ├── features/
│   ├── hooks/
│   ├── layouts/
│   ├── pages/
│   ├── routes/
│   ├── services/
│   ├── types/
│   └── utils/

---

# 6. React Rules

Components must remain small.

Target:

* Under 200 lines

Maximum:

* 300 lines

If larger:

Split components.

---

# 7. Feature-Based Organization

Use feature folders.

Example:

features/
│
├── auth/
├── users/
├── conversations/
├── messages/
└── presence/

Do NOT create giant shared folders.

---

# 8. Backend Architecture

Backend:

Go
Gin

Structure:

backend/
│
├── cmd/
│   └── api/
│
├── internal/
│   ├── auth/
│   ├── users/
│   ├── conversations/
│   ├── messages/
│   ├── websocket/
│   ├── middleware/
│   └── database/
│
├── pkg/
│
├── migrations/
│
└── tests/

---

# 9. Layer Rules

Every feature follows:

Handler
↓
Service
↓
Repository
↓
Database

Never skip layers.

---

# 10. Handler Rules

Handlers:

* Parse requests
* Return responses

Handlers must NOT:

* Contain business logic
* Contain SQL
* Manage websocket state

---

# 11. Service Rules

Services:

* Business logic only

Services must NOT:

* Return HTTP responses
* Use Gin Context

---

# 12. Repository Rules

Repositories:

* Database access only

Repositories must:

* Contain SQL queries
* Hide persistence details

---

# 13. Database Rules

Repositories are the only layer allowed to query PostgreSQL.

Forbidden:

Database access from:

* Handler
* Service
* WebSocket Hub

---

# 14. WebSocket Architecture

Structure:

Client
↔ WebSocket
↔ Hub
↔ Service
↔ Repository

WebSocket layer must remain separate from HTTP layer.

---

# 15. WebSocket Hub Rules

Hub responsibilities:

* Register clients
* Unregister clients
* Broadcast events
* Manage connections

Hub must NOT:

* Query database
* Execute business logic

---

# 16. Message Flow

Message Flow:

Client
↓
WebSocket
↓
Hub
↓
Message Service
↓
Repository
↓
Database

Response:

Database
↓
Service
↓
Hub
↓
Recipient

---

# 17. State Management Rules

Frontend state:

React Context

Allowed:

* Context API

Forbidden:

* Redux
* MobX

Without approval.

V1 should remain simple.

---

# 18. API Rules

REST API only.

Examples:

POST /auth/register

POST /auth/login

GET /users

GET /conversations

POST /conversations

GET /messages

---

# 19. Validation Rules

All inputs must be validated.

Validation required for:

* Register
* Login
* Conversation Creation
* Messaging

Never trust frontend input.

---

# 20. Authorization Rules

Users must access only:

* Their conversations
* Their messages

Unauthorized access must return:

403 Forbidden

---

# 21. Database Transaction Rules

Use transactions when:

* Creating conversations
* Creating messages

Never leave partial writes.

---

# 22. Error Handling Rules

Never expose:

* SQL errors
* Stack traces
* Internal server details

Return safe messages.

---

# 23. Logging Rules

Allowed:

* Application logs
* Error logs

Forbidden:

* Password logging
* JWT logging

---

# 24. Docker Rules

Services:

frontend
backend
postgres

Target:

docker compose up

must start entire application.

---

# 25. Docker Networking

Containers communicate through:

Docker Network

Never use localhost between containers.

Use service names.

Example:

backend → postgres

NOT:

localhost

---

# 26. Testing Rules

Backend:

Go tests

Frontend:

Vitest

Critical areas:

* Authentication
* Conversations
* Messaging
* Authorization

---

# 27. Security Rules

Must prevent:

* IDOR
* Unauthorized messaging
* Conversation spoofing
* JWT forgery

---

# 28. Scalability Preparation

Architecture should support future:

* Read Receipts
* Typing Indicators
* Group Chat

Without major rewrites.

Do not implement these features yet.

---

# 29. Forbidden Anti-Patterns

Do NOT create:

* God Service
* God Handler
* God Component
* Massive Utility File

If a file grows too large:

Refactor immediately.

---

# 30. Definition of Good Architecture

Good architecture means:

* Small files
* Clear ownership
* Clear responsibility
* Easy testing
* Easy future extension
