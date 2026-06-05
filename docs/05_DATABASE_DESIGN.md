# Database Design

# ChatSphere V1

Version: 1.0

Status: Approved

Database Engine:

PostgreSQL

---

# 1. Database Overview

ChatSphere V1 uses a relational database design.

Core Entities:

* Users
* Conversations
* Conversation Participants
* Messages

Design Goals:

* Prevent duplicate private conversations
* Support future group chat expansion
* Support efficient message retrieval
* Support online presence tracking

---

# 2. Entity Relationship Diagram

Users
↓
Conversation Participants
↓
Conversations
↓
Messages

---

Simplified ERD:

Users
│
├── ConversationParticipants
│       │
│       └── Conversations
│               │
│               └── Messages

---

# 3. Tables

## users

Purpose:

Stores registered users.

Columns:

| Column        | Type         | Constraints   |
| ------------- | ------------ | ------------- |
| id            | BIGSERIAL    | PK            |
| name          | VARCHAR(100) | NOT NULL      |
| email         | VARCHAR(255) | UNIQUE        |
| password_hash | VARCHAR(255) | NOT NULL      |
| is_online     | BOOLEAN      | DEFAULT FALSE |
| last_seen_at  | TIMESTAMP    | NULL          |
| created_at    | TIMESTAMP    | NOT NULL      |
| updated_at    | TIMESTAMP    | NOT NULL      |

Indexes:

* email (unique)
* is_online

---

## conversations

Purpose:

Represents a private conversation.

Columns:

| Column     | Type      | Constraints |
| ---------- | --------- | ----------- |
| id         | BIGSERIAL | PK          |
| created_at | TIMESTAMP | NOT NULL    |
| updated_at | TIMESTAMP | NOT NULL    |

Indexes:

* created_at

Notes:

Conversation metadata intentionally minimal.

Future versions may add:

* conversation_type
* group_name
* avatar

Not required in V1.

---

## conversation_participants

Purpose:

Stores participants of a conversation.

Columns:

| Column          | Type      | Constraints |
| --------------- | --------- | ----------- |
| id              | BIGSERIAL | PK          |
| conversation_id | BIGINT    | FK          |
| user_id         | BIGINT    | FK          |
| created_at      | TIMESTAMP | NOT NULL    |

Foreign Keys:

conversation_id → conversations.id

user_id → users.id

Indexes:

* conversation_id
* user_id
* (conversation_id, user_id) UNIQUE

---

Rules:

A user cannot appear twice in the same conversation.

---

## messages

Purpose:

Stores chat messages.

Columns:

| Column          | Type      | Constraints |
| --------------- | --------- | ----------- |
| id              | BIGSERIAL | PK          |
| conversation_id | BIGINT    | FK          |
| sender_id       | BIGINT    | FK          |
| content         | TEXT      | NOT NULL    |
| sent_at         | TIMESTAMP | NOT NULL    |
| created_at      | TIMESTAMP | NOT NULL    |
| updated_at      | TIMESTAMP | NOT NULL    |

Foreign Keys:

conversation_id → conversations.id

sender_id → users.id

Indexes:

* conversation_id
* sender_id
* sent_at

---

# 4. Relationships

Users
1 → Many ConversationParticipants

Conversations
1 → Many ConversationParticipants

Users
1 → Many Messages

Conversations
1 → Many Messages

---

# 5. Conversation Rules

Private conversations only.

Exactly two users participate.

Example:

Conversation #1

Participants:

* User A
* User B

Allowed:

2 participants

Forbidden:

3+ participants

---

# 6. Duplicate Conversation Prevention

Rule:

Two users must have only one private conversation.

Example:

User A + User B

If conversation exists:

Return existing conversation.

Do not create a new one.

Implementation responsibility:

Service Layer

Not database layer.

---

# 7. Message Rules

Messages belong to:

* One conversation
* One sender

Messages cannot exist without:

* Valid conversation
* Valid sender

---

# 8. Presence Rules

Online status stored in:

users.is_online

Last activity stored in:

users.last_seen_at

Updated when:

* User connects
* User disconnects

---

# 9. Query Requirements

Common Queries:

Get Conversation List

Get Conversation Messages

Search Users

Update Presence

Send Message

Retrieve Message History

Indexes must support these queries efficiently.

---

# 10. Future Expansion Support

Current schema must support future:

V2

* Read Receipts
* Typing Indicator

V3

* Group Chat

Potential future tables:

message_reads

typing_events

groups

group_members

Do not implement yet.

---

# 11. Constraints

Content Length:

Message content maximum:

2000 characters

User Name maximum:

100 characters

Email maximum:

255 characters

---

# 12. Cascade Rules

Conversation Delete

Deletes:

* Participants
* Messages

User Delete

Not supported in V1.

No hard delete workflow required.

---

# 13. Security Rules

Users may only access:

* Their conversations
* Their messages

Conversation ownership validation required before message retrieval.

Never trust client-provided user IDs.

Always use authenticated user context.

---

# 14. Migration Order

1. users

2. conversations

3. conversation_participants

4. messages

Order is mandatory due to foreign key dependencies.

---

# 15. Acceptance Criteria

Database design accepted when:

✓ Users can register

✓ Users can login

✓ Conversations can be created

✓ Duplicate conversations prevented

✓ Messages persist correctly

✓ Presence data stored

✓ Foreign keys valid

✓ Indexes created

✓ Future group chat remains possible without redesign
