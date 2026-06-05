# Software Requirements Specification (SRS)

# ChatSphere V1

Version: 1.0

Status: Approved

---

# 1. Introduction

## Purpose

This document defines the functional and non-functional requirements for ChatSphere V1.

It serves as the primary reference for implementation, testing, and validation.

---

## Product Scope

ChatSphere is a real-time messaging web application that allows users to:

* Register
* Login
* Discover users
* Create private conversations
* Exchange messages in real time
* View online presence

---

# 2. User Roles

## Standard User

Permissions:

* Register account
* Login
* View own profile
* View user directory
* Create conversations
* Send messages
* Receive messages
* View conversation history

No administrator role exists in V1.

---

# 3. Authentication Requirements

## AUTH-001 Registration

Users must register using:

* Name
* Email
* Password

Rules:

* Name required
* Email required
* Email unique
* Password minimum 8 characters

Expected Result:

User account created successfully.

---

## AUTH-002 Login

Users must login using:

* Email
* Password

Expected Result:

System returns JWT token.

---

## AUTH-003 Logout

Users can logout.

Expected Result:

Frontend removes stored token.

---

## AUTH-004 Unauthorized Access

Unauthenticated users cannot access:

* User Directory
* Conversations
* Messages
* Presence Data

Expected Result:

401 Unauthorized

---

# 4. User Requirements

## USER-001 User Directory

Users can view a list of registered users.

Displayed Fields:

* User Name
* Online Status

---

## USER-002 User Search

Users can search users by:

* Name

Search should be case-insensitive.

---

## USER-003 Own Profile

Users can view:

* Name
* Email
* Last Seen

Editing profile is not included in V1.

---

# 5. Conversation Requirements

## CONV-001 Create Conversation

Users can create private conversations.

Conversation participants:

* User A
* User B

Only two users allowed.

---

## CONV-002 Duplicate Prevention

System must prevent duplicate conversations.

Example:

If User A and User B already have a conversation,

a new conversation must not be created.

Expected Result:

Return existing conversation.

---

## CONV-003 View Conversations

Users can view:

* Conversation List
* Last Message Preview
* Last Activity Timestamp

Ordered By:

Most recent activity first.

---

## CONV-004 Access Control

Users may only access conversations where they are participants.

Violation Result:

403 Forbidden

---

# 6. Messaging Requirements

## MSG-001 Send Message

Users can send text messages.

Required Fields:

* Conversation ID
* Message Content

Rules:

* Content required
* Content maximum 2000 characters

---

## MSG-002 Store Message

Every message must be persisted to PostgreSQL.

Expected Result:

Message remains available after reconnect.

---

## MSG-003 Receive Message

Recipients must receive messages in real time.

Page refresh must not be required.

---

## MSG-004 Message History

Users can load previous messages.

Messages ordered:

Oldest → Newest

---

## MSG-005 Message Ownership

Users can only view messages from conversations they belong to.

Violation Result:

403 Forbidden

---

# 7. Presence Requirements

## PRES-001 Online Status

System must mark users as:

* Online
* Offline

---

## PRES-002 Last Seen

System must store:

Last Active Timestamp

Displayed as:

Last Seen

---

## PRES-003 Presence Update

Presence changes must be reflected without page refresh.

WebSocket communication required.

---

# 8. WebSocket Requirements

## WS-001 Connection

Authenticated users may establish WebSocket connections.

Authentication required before connection.

---

## WS-002 Invalid Authentication

Invalid JWT must reject connection.

Expected Result:

Connection refused.

---

## WS-003 Message Event

Message events must be delivered through WebSocket.

---

## WS-004 Presence Event

Presence updates must be delivered through WebSocket.

---

## WS-005 Reconnection

Client should reconnect automatically when disconnected.

---

# 9. Database Requirements

## DB-001 PostgreSQL

Database engine:

PostgreSQL

---

## DB-002 Foreign Keys

All relationships require foreign keys.

---

## DB-003 Indexes

Indexes required for:

* Users
* Conversations
* Participants
* Messages

---

## DB-004 Timestamps

All entities must contain:

* Created At
* Updated At

Messages additionally require:

* Sent At

---

# 10. Security Requirements

## SEC-001 Password Security

Passwords must be hashed using bcrypt.

Plain-text passwords prohibited.

---

## SEC-002 JWT Security

Protected endpoints require valid JWT.

---

## SEC-003 Authorization

Users may access only:

* Their conversations
* Their messages

---

## SEC-004 Input Validation

All API requests must be validated.

---

## SEC-005 IDOR Prevention

Direct access to other users' data is prohibited.

---

# 11. Performance Requirements

## PERF-001 Message Delivery

Local message delivery target:

< 1 second

---

## PERF-002 Conversation Loading

Conversation list target:

< 2 seconds

---

## PERF-003 User Search

Search target:

< 2 seconds

---

# 12. Reliability Requirements

## REL-001 Persistence

Messages must not be lost after refresh.

---

## REL-002 Recovery

Client should reconnect after temporary disconnect.

---

## REL-003 Error Handling

Unexpected failures must return safe error responses.

---

# 13. UI Requirements

## UI-001 Responsive Layout

Application must work on:

* Desktop
* Tablet
* Mobile

---

## UI-002 Chat Layout

Main layout:

* Sidebar
* Conversation List
* Chat Window

---

## UI-003 Real-Time Updates

Messages appear instantly.

No manual refresh allowed.

---

# 14. Out of Scope

The following are excluded:

* Group Chat
* Message Editing
* Message Deletion
* Typing Indicator
* Read Receipts
* File Upload
* Image Upload
* Voice Notes
* Video Calls
* Push Notifications
* End-to-End Encryption

---

# 15. Acceptance Criteria

ChatSphere V1 is accepted when:

✓ Registration works

✓ Login works

✓ JWT authentication works

✓ User search works

✓ Conversation creation works

✓ Duplicate conversation prevention works

✓ Real-time messaging works

✓ Online status works

✓ Docker environment works

✓ Tests pass

✓ Documentation complete
