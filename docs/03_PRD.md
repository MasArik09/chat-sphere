# Product Requirements Document (PRD)

# ChatSphere V1

Version: 1.0

Status: Approved

---

# 1. Product Overview

## Product Name

ChatSphere

## Product Type

Real-Time Messaging Web Application

## Product Vision

ChatSphere is a modern real-time messaging platform built to explore scalable communication systems, WebSocket architecture, authentication, and full-stack application development using Go and React.

The application allows users to register, discover other users, start private conversations, exchange messages in real time, and see online presence information.

---

# 2. Problem Statement

Many beginner chat applications focus only on CRUD operations and do not properly implement:

* Real-time communication
* WebSocket architecture
* User presence tracking
* Conversation modeling
* Modern frontend-backend separation

ChatSphere aims to provide a practical learning project that demonstrates how modern chat systems work while maintaining a clean and scalable architecture.

---

# 3. Goals

## Primary Goals

* Learn Go backend development.
* Learn WebSocket communication.
* Learn React + TypeScript frontend architecture.
* Learn PostgreSQL relational design.
* Learn Docker-based development workflow.
* Build a portfolio-quality real-time application.

## Success Criteria

A user can:

* Register an account.
* Login securely.
* View available users.
* Start a private conversation.
* Send messages.
* Receive messages instantly without page refresh.
* See online/offline user status.

---

# 4. Target Users

## Primary Users

* Students
* Developers
* Portfolio builders
* Small teams

## User Characteristics

Users want:

* Fast communication
* Simple interface
* Real-time updates
* Easy navigation

---

# 5. Project Scope

## Included in V1

### Authentication

* Register
* Login
* Logout
* JWT Authentication

### User Management

* User Directory
* User Search
* User Profile

### Private Messaging

* Create Conversation
* View Conversation List
* View Message History
* Send Message
* Receive Message

### Real-Time Features

* WebSocket Connection
* Instant Message Delivery
* Online Status
* Offline Status
* Last Seen

### Dashboard

* Recent Conversations
* Conversation Preview
* Unread Message Counter

---

# 6. Out of Scope (Not Included in V1)

The following features are intentionally excluded:

## Messaging

* Group Chat
* Message Editing
* Message Deletion
* Message Reactions
* Message Forwarding
* Message Pinning

## Media

* Image Upload
* File Upload
* Video Upload
* Voice Notes

## Communication

* Voice Calls
* Video Calls
* Screen Sharing

## Advanced Features

* Push Notifications
* End-to-End Encryption
* AI Assistant
* Multi-Device Synchronization

These features may be considered in future versions.

---

# 7. Core User Stories

## Authentication

As a new user,

I want to register an account,

So that I can access ChatSphere.

---

As a user,

I want to login securely,

So that I can access my conversations.

---

## User Directory

As a user,

I want to browse other users,

So that I can start conversations.

---

As a user,

I want to search users,

So that I can quickly find someone.

---

## Conversation

As a user,

I want to create a private conversation,

So that I can chat with another user.

---

As a user,

I want to see my conversation list,

So that I can continue previous chats.

---

## Messaging

As a user,

I want to send a message,

So that the other user receives it instantly.

---

As a user,

I want to receive messages in real time,

So that communication feels immediate.

---

## Presence

As a user,

I want to see who is online,

So that I know who is available.

---

As a user,

I want to see last seen information,

So that I know when someone was last active.

---

# 8. Functional Requirements

## FR-001 User Registration

Users must be able to create an account using:

* Name
* Email
* Password

---

## FR-002 User Login

Users must authenticate using:

* Email
* Password

Authentication must return JWT tokens.

---

## FR-003 User Search

Users can search other users by:

* Name

---

## FR-004 Create Conversation

Users can create private conversations.

A conversation can only contain:

* User A
* User B

---

## FR-005 View Conversations

Users can view:

* Conversation list
* Last message
* Last activity timestamp

---

## FR-006 Send Message

Users can send text messages.

---

## FR-007 Receive Message

Messages must appear without page refresh.

WebSocket communication is required.

---

## FR-008 Presence Tracking

System must track:

* Online
* Offline
* Last Seen

---

# 9. Non-Functional Requirements

## Performance

* Message delivery under 1 second locally.
* Conversation list loads under 2 seconds.

## Security

* Password hashing using bcrypt.
* JWT authentication.
* Authorization checks for conversations.
* Input validation.

## Reliability

* WebSocket reconnection support.
* Graceful error handling.

## Maintainability

* Modular architecture.
* Feature-based folder structure.
* Clear separation between frontend and backend.

---

# 10. Technology Stack

## Frontend

* React
* TypeScript
* Vite
* Tailwind CSS

## Backend

* Go
* Gin Framework

## Database

* PostgreSQL

## Realtime

* Native WebSocket

## Authentication

* JWT

## Containerization

* Docker
* Docker Compose

---

# 11. Constraints

* No paid services.
* No cloud dependency required.
* Must run completely locally.
* Must be open-source friendly.
* Must be Docker-compatible.

---

# 12. Future Versions

## V2

* Typing Indicator
* Message Read Status
* Message Delivery Status

## V3

* Group Chat
* Image Upload
* File Sharing

## V4

* Voice Notes
* Push Notifications

## V5

* End-to-End Encryption
* Multi-Device Support

---

# 13. Release Definition

ChatSphere V1 is considered complete when:

* Authentication works.
* Private conversations work.
* Real-time messaging works.
* Online status works.
* Docker setup works.
* Automated tests pass.
* Documentation is complete.
