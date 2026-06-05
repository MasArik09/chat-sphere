# Implementation Order

# ChatSphere V1

Version: 1.0

Status: Mandatory

---

# Purpose

This document defines the exact implementation order.

AI must follow this order.

Do NOT skip phases.

Do NOT implement future phases early.

---

# Global Rule

Before every phase:

1. Read PROJECT_RULES.md
2. Read ARCHITECTURE_RULES.md
3. Read SRS.md
4. Read DATABASE_DESIGN.md
5. Read TASK_BREAKDOWN.md

After every phase:

STOP

Generate implementation report.

Wait for approval.

---

# Phase 0

## Setup & Environment

### Backend

Create:

cmd/api

internal/

pkg/

migrations/

tests/

---

### Frontend

Create:

src/

features/

components/

services/

hooks/

types/

---

### Docker

Create:

Dockerfile (backend)

Dockerfile (frontend)

docker-compose.yml

.env.example

---

### Verification

Must verify:

- Backend starts
- Frontend starts
- PostgreSQL starts

STOP

---

# Phase 1

## Database Foundation

### Create Migrations

Order:

1. users

2. conversations

3. conversation_participants

4. messages

---

### Create Entities

User

Conversation

ConversationParticipant

Message

---

### Verification

Run migrations.

Verify:

- Foreign Keys
- Constraints
- Indexes

STOP

---

# Phase 2

## Authentication Module

### Backend

Create:

auth repository

auth service

auth handler

jwt middleware

---

### API

POST /auth/register

POST /auth/login

---

### Frontend

Create:

Register Page

Login Page

Auth Service

Auth Context

Protected Routes

---

### Testing

Register

Login

JWT

Unauthorized Access

---

### Verification

Authentication fully functional.

STOP

---

# Phase 3

## User Directory Module

### Backend

User Repository

User Service

User Handler

---

### API

GET /users

GET /users/search

---

### Frontend

User Directory Page

Search Users

User Cards

---

### Testing

User listing

User search

Authorization

---

### Verification

User directory functional.

STOP

---

# Phase 4

## Conversation Module

### Backend

Conversation Repository

Conversation Service

Conversation Handler

---

### Features

Create Conversation

Get Conversation

List Conversations

Duplicate Prevention

Ownership Validation

---

### API

POST /conversations

GET /conversations

GET /conversations/{id}

---

### Frontend

Conversation List

Conversation Item

Start Chat Flow

---

### Testing

Create conversation

Duplicate prevention

Authorization

---

### Verification

Conversation module functional.

STOP

---

# Phase 5

## Messaging Module

### Backend

Message Repository

Message Service

Message Handler

---

### Features

Store Messages

Load Message History

Ownership Validation

Pagination

---

### API

GET /conversations/{id}/messages

---

### Frontend

Chat Window

Message Composer

Message List

Message Bubble

---

### Testing

Message storage

Message retrieval

Authorization

---

### Verification

Messaging module functional.

STOP

---

# Phase 6

## WebSocket Realtime Module

### Backend

WebSocket Handler

Hub

Client Manager

Event Dispatcher

---

### Events

message.send

message.received

---

### Frontend

WebSocket Service

Connection Management

Message Event Handling

Reconnect Logic

---

### Testing

Connection

Message Delivery

Reconnect

---

### Verification

Realtime messaging works.

STOP

---

# Phase 7

## Presence Module

### Backend

Presence Service

Presence Events

Last Seen Updates

---

### Frontend

Presence Badge

Online Indicator

Last Seen Display

---

### Events

presence.online

presence.offline

---

### Testing

Online status

Offline status

Last seen

---

### Verification

Presence functional.

STOP

---

# Phase 8

## UI Refinement

### Improve

Responsive Layout

Loading States

Error States

Empty States

Accessibility

---

### Verification

Desktop UI

Tablet UI

Mobile UI

STOP

---

# Phase 9

## Docker Integration

### Backend Container

Dockerfile

---

### Frontend Container

Dockerfile

---

### Database Container

PostgreSQL Service

---

### Compose

docker-compose.yml

---

### Verification

docker compose up

All services healthy

STOP

---

# Phase 10

## Testing & Security Audit

### Backend Tests

Authentication

Users

Conversations

Messages

WebSocket

Presence

---

### Security Review

JWT

Authorization

IDOR

Validation

---

### Verification

All tests passing.

STOP

---

# Phase 11

## Documentation & Release

### Create

README.md

LICENSE

Screenshots

Architecture Notes

Setup Guide

---

### Verification

Fresh clone works.

Repository ready for public release.

STOP

---

# Forbidden Actions

Do NOT:

- Skip phases
- Implement future phases
- Create undocumented features
- Introduce new technologies
- Refactor unrelated modules

Without approval.

---

# Final Rule

One phase at a time.

Implementation must remain small, verifiable, and reversible.