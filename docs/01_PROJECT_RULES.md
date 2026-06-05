# Project Rules

# ChatSphere V1

Version: 1.0

Status: Mandatory

---

# 1. Purpose

This document defines mandatory project rules for all development activities.

Every implementation must follow these rules.

These rules override AI assumptions.

---

# 2. Project Philosophy

ChatSphere is:

* Educational
* Portfolio Quality
* Maintainable
* Scalable
* Docker First

ChatSphere is NOT:

* Enterprise Software
* Production Messaging Platform
* WhatsApp Clone
* Microservice System

---

# 3. Development Approach

Development must follow:

Requirements
→ Design
→ Implementation
→ Verification
→ Testing

Never skip documentation.

Never implement features not described in documentation.

---

# 4. Scope Control

Only implement features defined in:

* PRD.md
* SRS.md
* Task Breakdown

Anything outside scope must be rejected.

---

# 5. AI Agent Rules

Before coding:

AI must read:

* 01_PROJECT_RULES.md
* 02_ARCHITECTURE_RULES.md
* 03_PRD.md
* 04_SRS.md

If documentation conflicts:

Priority order:

1. PROJECT_RULES
2. ARCHITECTURE_RULES
3. SRS
4. PRD

---

# 6. Technology Rules

Frontend:

* React
* TypeScript
* Vite
* Tailwind CSS

Backend:

* Go
* Gin

Database:

* PostgreSQL

Realtime:

* Native WebSocket

Containerization:

* Docker
* Docker Compose

Do not replace technologies without approval.

---

# 7. Forbidden Technologies

Do NOT introduce:

* Firebase
* Supabase
* Pusher
* MongoDB
* Prisma
* GraphQL
* Microservices
* Kubernetes

Without explicit approval.

---

# 8. Authentication Rules

Authentication method:

JWT

Requirements:

* Access Token
* Password Hashing (bcrypt)

Do not implement:

* OAuth
* Google Login
* GitHub Login

For V1.

---

# 9. Realtime Rules

Use:

Native WebSocket

Do not introduce:

* Socket.IO
* Pusher
* Third-party realtime providers

Without approval.

---

# 10. Database Rules

Database:

PostgreSQL

Requirements:

* Foreign Keys
* Proper Indexing
* Explicit Constraints

Avoid:

* Unnecessary JSON columns
* Denormalized structures

Without justification.

---

# 11. Docker Rules

Project must run using:

docker compose up

A fresh clone must be runnable using documented setup steps.

Docker is mandatory.

---

# 12. Security Rules

Must implement:

* Password Hashing
* JWT Authentication
* Authorization Checks
* Input Validation

Must prevent:

* IDOR
* Unauthorized Conversations
* Unauthorized Message Access

---

# 13. Testing Rules

Every major feature must include tests.

Minimum areas:

* Authentication
* Conversations
* Messaging
* Authorization

Code is not complete without tests.

---

# 14. Documentation Rules

All architecture decisions must be documented.

Required docs:

* PRD
* SRS
* Database Design
* System Design
* Architecture Rules
* Task Breakdown

---

# 15. Git Rules

Commit messages must be meaningful.

Bad:

"fix"

"update"

"final"

Good:

"feat(auth): implement JWT authentication"

"feat(chat): add websocket messaging"

"fix(conversation): prevent duplicate private chats"

---

# 16. UI Rules

Design principles:

* Simple
* Clean
* Responsive
* Consistent

Avoid:

* Excessive animations
* Complex dashboards
* Over-designed interfaces

---

# 17. Performance Rules

Avoid:

* N+1 queries
* Unbounded websocket handlers
* Repeated database calls

Prefer:

* Indexed queries
* Efficient joins
* Connection management

---

# 18. Scalability Rules

Architecture should support future additions:

* Typing Indicator
* Read Receipts
* Group Chat

Without major rewrites.

Do not implement these features yet.

---

# 19. Feature Freeze Rule

V1 includes only:

* Authentication
* User Directory
* Private Conversations
* Realtime Messaging
* Online Status

Any other feature requires approval.

---

# 20. Definition of Done

A feature is complete only when:

* Requirements implemented
* Tests pass
* Documentation updated
* Code reviewed
* No critical bugs remain
