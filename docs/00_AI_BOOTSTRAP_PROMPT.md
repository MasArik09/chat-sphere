# AI Bootstrap Prompt

# ChatSphere V1

Version: 1.0

Status: Mandatory

---

You are the implementation agent for ChatSphere V1.

Before writing any code, you MUST read and understand the following documents:

1. docs/01_PROJECT_RULES.md
2. docs/02_ARCHITECTURE_RULES.md
3. docs/03_PRD.md
4. docs/04_SRS.md
5. docs/05_DATABASE_DESIGN.md
6. docs/06_SYSTEM_DESIGN.md
7. docs/07_UI_UX_FLOW.md
8. docs/08_TASK_BREAKDOWN.md
9. docs/09_IMPLEMENTATION_ORDER.md
10. docs/10_DOCKER_RULES.md
11. docs/11_CODING_STANDARDS.md

---

# Mission

Implement ChatSphere V1 according to the documentation.

The implementation must be:

- Maintainable
- Modular
- Testable
- Documented

Follow the documented architecture.

Do not invent features.

Do not introduce new technologies.

Do not change architecture without approval.

---

# Required Reading Order

Always read documents in this order:

1. PROJECT_RULES
2. ARCHITECTURE_RULES
3. PRD
4. SRS
5. DATABASE_DESIGN
6. SYSTEM_DESIGN
7. UI_UX_FLOW
8. TASK_BREAKDOWN
9. IMPLEMENTATION_ORDER
10. DOCKER_RULES
11. CODING_STANDARDS

---

# Implementation Rules

Implement only the requested phase.

Example:

If user requests:

"Implement Phase 2"

You must implement:

Phase 2 only.

Do not implement:

- Phase 3
- Phase 4
- Future phases

---

# Stop Rule

After completing a phase:

STOP.

Generate a report.

Wait for approval.

Never continue automatically.

---

# Report Format

Use the following format:

## Files Created

List every new file.

---

## Files Modified

List every modified file.

---

## Features Implemented

Describe completed functionality.

---

## Features NOT Implemented

List anything intentionally left out.

---

## Verification

Explain how implementation was verified.

---

## Test Results

Provide:

- Tests executed
- Results
- Failures if any

---

## Next Recommended Phase

Suggest the next phase.

Then STOP.

---

# Architecture Enforcement

Respect:

Handler
↓
Service
↓
Repository
↓
Database

Never bypass layers.

---

# Backend Rules

Use:

- Go
- Gin
- PostgreSQL

Do not introduce:

- Fiber
- Echo
- GraphQL
- MongoDB
- Prisma

---

# Frontend Rules

Use:

- React
- TypeScript
- Vite
- Tailwind CSS

Do not introduce:

- Next.js
- Redux
- MobX

Without approval.

---

# Realtime Rules

Use:

Native WebSocket

Do not introduce:

- Socket.IO
- Pusher
- Third-party realtime services

---

# Docker Rules

Project must run using:

docker compose up

Do not create manual-only setup steps.

---

# Security Rules

Always enforce:

- JWT authentication
- Authorization checks
- Conversation ownership validation
- Message ownership validation

Prevent:

- IDOR
- Unauthorized access

---

# Testing Rules

Every major feature requires tests.

Do not mark a phase complete without verification.

---

# Documentation Rules

If implementation changes architecture:

Stop.

Explain the conflict.

Request approval.

Do not proceed automatically.

---

# File Size Rules

Follow CODING_STANDARDS.md.

If a file becomes too large:

Refactor.

Do not create:

- God Services
- God Handlers
- God Components

---

# Final Rule

When uncertain:

Follow documentation.

When documentation conflicts:

Priority order:

1. PROJECT_RULES
2. ARCHITECTURE_RULES
3. SRS
4. PRD

Never assume.

Never invent requirements.

Implement only what is documented.