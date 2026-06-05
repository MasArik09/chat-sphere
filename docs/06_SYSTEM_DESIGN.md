# System Design

# ChatSphere V1

Version: 1.0

Status: Approved

---

# 1. System Overview

ChatSphere is a real-time messaging application.

Architecture Style:

Modular Monolith

Communication Types:

* REST API
* WebSocket

Deployment Style:

Docker Compose

---

# 2. High-Level Architecture

Frontend
↓
REST API
↓
Go Backend
↓
PostgreSQL

Realtime Channel

Frontend
↔ WebSocket
↔ Go Backend

---

# 3. Technology Stack

Frontend

* React
* TypeScript
* Vite
* Tailwind CSS

Backend

* Go
* Gin

Database

* PostgreSQL

Authentication

* JWT

Realtime

* Native WebSocket

Containerization

* Docker
* Docker Compose

---

# 4. Container Architecture

Docker Compose Services

frontend

backend

postgres

Network:

chatsphere-network

---

Communication:

frontend
↓
backend

backend
↓
postgres

frontend
↔
backend websocket

---

# 5. Backend Module Design

internal/

auth/

users/

conversations/

messages/

websocket/

middleware/

database/

Each module contains:

handler

service

repository

models

---

# 6. Frontend Module Design

features/

auth/

users/

conversations/

messages/

presence/

Each feature contains:

components

hooks

services

types

---

# 7. Authentication Flow

Registration

Client
↓
POST /auth/register
↓
Backend
↓
Database

Response

Success

---

Login

Client
↓
POST /auth/login
↓
JWT Generated
↓
Client Stores Token

---

Authenticated Requests

Client
↓
Authorization: Bearer Token
↓
Backend Validation
↓
Access Granted

---

# 8. Conversation Flow

Create Conversation

User A
↓
Select User B
↓
POST /conversations
↓
Backend Validation
↓
Check Existing Conversation
↓
Create or Return Existing
↓
Response

---

View Conversations

Client
↓
GET /conversations
↓
Backend
↓
Database
↓
Response

---

# 9. Messaging Flow

Send Message

Client
↓
WebSocket Event
↓
WebSocket Hub
↓
Message Service
↓
Repository
↓
PostgreSQL

After Save

Repository
↓
Service
↓
Hub
↓
Recipient

---

# 10. Message History Flow

Client
↓
GET /conversations/{id}/messages
↓
Authorization Check
↓
Repository
↓
Database
↓
Response

---

# 11. Presence Flow

User Connects

Client
↓
WebSocket Connect
↓
Hub
↓
Update Presence
↓
is_online = true

---

User Disconnects

Client Disconnect
↓
Hub
↓
Update Presence
↓
is_online = false
↓
last_seen_at updated

---

# 12. WebSocket Architecture

Components

WebSocket Handler

Hub

Client

Message Service

Presence Service

---

Responsibilities

Handler

* Upgrade Connection
* Register Client

---

Hub

* Register Connections
* Remove Connections
* Route Events
* Broadcast Events

---

Client

* Send Events
* Receive Events

---

Service Layer

* Business Logic
* Validation
* Persistence

---

# 13. Event Definitions

Client Events

message.send

---

Server Events

message.received

presence.online

presence.offline

---

# 14. API Design

Authentication

POST /auth/register

POST /auth/login

---

Users

GET /users

GET /users/search

---

Conversations

GET /conversations

POST /conversations

GET /conversations/{id}

---

Messages

GET /conversations/{id}/messages

---

# 15. Authorization Strategy

Protected Endpoints

Require JWT

Rules

User must belong to conversation.

If not:

403 Forbidden

---

# 16. Error Handling Strategy

Standard Response

Success:

{
"success": true
}

---

Error:

{
"success": false,
"message": "Error description"
}

---

Validation Error:

{
"success": false,
"errors": {}
}

---

# 17. Database Access Strategy

Only repositories access PostgreSQL.

Forbidden:

Handler → Database

Service → Database

WebSocket Hub → Database

---

Required Path:

Handler
↓
Service
↓
Repository
↓
Database

---

# 18. Security Design

Authentication

JWT

Password Hashing

bcrypt

---

Protection

Authorization Checks

Input Validation

IDOR Prevention

Conversation Ownership Validation

---

# 19. Scalability Preparation

Architecture should support:

V2

Typing Indicator

Read Receipts

Delivery Status

---

V3

Group Chat

Media Upload

---

No implementation required in V1.

---

# 20. Testing Strategy

Backend

Unit Tests

Integration Tests

---

Frontend

Component Tests

Feature Tests

---

Critical Areas

Authentication

Conversations

Messaging

Authorization

WebSocket Events

---

# 21. Logging Strategy

Allowed

Application Logs

Error Logs

Connection Logs

---

Forbidden

Password Logs

JWT Logs

Sensitive User Data

---

# 22. Docker Startup Flow

docker compose up

↓

postgres starts

↓

backend starts

↓

frontend starts

↓

application available

---

# 23. Definition of Done

System design is accepted when:

✓ Authentication architecture defined

✓ Conversation architecture defined

✓ Messaging architecture defined

✓ Presence architecture defined

✓ WebSocket architecture defined

✓ Docker architecture defined

✓ Security architecture defined
