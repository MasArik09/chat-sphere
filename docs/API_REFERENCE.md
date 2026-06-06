# ChatSphere V1 - API Reference Specification

This document details all HTTP REST API endpoints and WebSocket socket event payloads used by the ChatSphere platform.

---

## 1. Authentication Endpoints

All request/response payloads utilize JSON format.

### Register User
* **Endpoint**: `POST /api/v1/auth/register`
* **Headers**: `Content-Type: application/json`
* **Request Body**:
  ```json
  {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "secretpassword123"
  }
  ```
* **Success Response (201 Created)**:
  ```json
  {
    "message": "User registered successfully",
    "user": {
      "id": 1,
      "name": "Jane Doe",
      "email": "jane@example.com",
      "is_online": false,
      "last_seen_at": null,
      "created_at": "2026-06-06T17:15:00Z",
      "updated_at": "2026-06-06T17:15:00Z"
    }
  }
  ```
* **Error Response (409 Conflict)**:
  ```json
  {
    "error": "Email already registered"
  }
  ```

---

### Login User
* **Endpoint**: `POST /api/v1/auth/login`
* **Request Body**:
  ```json
  {
    "email": "jane@example.com",
    "password": "secretpassword123"
  }
  ```
* **Success Response (200 OK)**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NzA5ODM0MDB9...",
    "user": {
      "id": 1,
      "name": "Jane Doe",
      "email": "jane@example.com"
    }
  }
  ```
* **Error Response (401 Unauthorized)**:
  ```json
  {
    "error": "Invalid email or password"
  }
  ```

---

### Get Profile (Me)
* **Endpoint**: `GET /api/v1/auth/me`
* **Headers**: `Authorization: Bearer <token>`
* **Success Response (200 OK)**:
  ```json
  {
    "user": {
      "id": 1,
      "name": "Jane Doe",
      "email": "jane@example.com",
      "is_online": true,
      "last_seen_at": "2026-06-06T17:20:00Z",
      "created_at": "2026-06-06T17:15:00Z",
      "updated_at": "2026-06-06T17:20:00Z"
    }
  }
  ```

---

## 2. Conversation Endpoints

All endpoints below require authentication headers: `Authorization: Bearer <token>`.

### Create Conversation
* **Endpoint**: `POST /api/v1/conversations`
* **Request Body**:
  ```json
  {
    "participant_ids": [2]
  }
  ```
* **Success Response (201 Created)**:
  ```json
  {
    "conversation": {
      "id": 1,
      "created_at": "2026-06-06T17:30:00Z",
      "updated_at": "2026-06-06T17:30:00Z"
    }
  }
  ```

---

### List / Search Conversations
* **Endpoint**: `GET /api/v1/conversations`
* **Query Parameters**:
  - `search` (Optional): Filter conversations by participant name (e.g. `?search=Bob`).
* **Success Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "created_at": "2026-06-06T17:30:00Z",
      "updated_at": "2026-06-06T17:45:00Z",
      "participants": [
        {
          "id": 1,
          "conversation_id": 1,
          "user_id": 1,
          "name": "Jane Doe",
          "email": "jane@example.com",
          "is_online": true,
          "last_seen_at": "2026-06-06T17:20:00Z"
        },
        {
          "id": 2,
          "conversation_id": 1,
          "user_id": 2,
          "name": "Bob Smith",
          "email": "bob@example.com",
          "is_online": false,
          "last_seen_at": "2026-06-06T17:40:00Z"
        }
      ],
      "unread_count": 1,
      "last_message": {
        "id": 15,
        "conversation_id": 1,
        "sender_id": 2,
        "content": "Hey Jane!",
        "sent_at": "2026-06-06T17:45:00Z"
      }
    }
  ]
  ```

---

### Add Participant
* **Endpoint**: `POST /api/v1/conversations/:id/participants`
* **Request Body**:
  ```json
  {
    "user_id": 3
  }
  ```
* **Success Response (200 OK)**:
  ```json
  {
    "message": "Participant added successfully"
  }
  ```

---

### Remove Participant
* **Endpoint**: `DELETE /api/v1/conversations/:id/participants/:userId`
* **Success Response (200 OK)**:
  ```json
  {
    "message": "Participant removed successfully"
  }
  ```

---

### Update Read Receipt
* **Endpoint**: `POST /api/v1/conversations/:id/read`
* **Request Body**:
  ```json
  {
    "last_read_message_id": 15
  }
  ```
* **Success Response (200 OK)**:
  ```json
  {
    "message": "Read receipt updated successfully"
  }
  ```
* **Error Response (400 Bad Request - e.g. invalid message ID)**:
  ```json
  {
    "error": "Message ID does not belong to this conversation"
  }
  ```

---

## 3. Message Endpoints

### Send Message
* **Endpoint**: `POST /api/v1/conversations/:id/messages`
* **Request Body**:
  ```json
  {
    "content": "Hi there Bob!"
  }
  ```
* **Success Response (201 Created)**:
  ```json
  {
    "message": {
      "id": 16,
      "conversation_id": 1,
      "sender_id": 1,
      "content": "Hi there Bob!",
      "sent_at": "2026-06-06T17:50:00Z",
      "created_at": "2026-06-06T17:50:00Z",
      "updated_at": "2026-06-06T17:50:00Z"
    }
  }
  ```

---

### List Messages
* **Endpoint**: `GET /api/v1/conversations/:id/messages`
* **Success Response (200 OK)**:
  ```json
  [
    {
      "id": 15,
      "conversation_id": 1,
      "sender_id": 2,
      "content": "Hey Jane!",
      "sent_at": "2026-06-06T17:45:00Z",
      "created_at": "2026-06-06T17:45:00Z",
      "updated_at": "2026-06-06T17:45:00Z"
    },
    {
      "id": 16,
      "conversation_id": 1,
      "sender_id": 1,
      "content": "Hi there Bob!",
      "sent_at": "2026-06-06T17:50:00Z",
      "created_at": "2026-06-06T17:50:00Z",
      "updated_at": "2026-06-06T17:50:00Z"
    }
  ]
  ```

---

## 4. WebSocket Event Specification

WebSockets connections are initiated via:
`GET /ws?token=<jwt_token>`

All socket data is sent and received wrapped in a standard JSON wrapper:
```json
{
  "event": "event_name",
  "data": { ... }
}
```

### Incoming Client Events (Client -> Server)

#### Start Typing
* **Event**: `typing.start`
* **Data**:
  ```json
  {
    "conversation_id": 1
  }
  ```

#### Stop Typing
* **Event**: `typing.stop`
* **Data**:
  ```json
  {
    "conversation_id": 1
  }
  ```

---

### Outgoing Server Events (Server -> Client)

#### Partner Online Status Changed
* **Event**: `presence.online`
* **Data**:
  ```json
  {
    "user_id": 2
  }
  ```

#### Partner Offline Status Changed
* **Event**: `presence.offline`
* **Data**:
  ```json
  {
    "user_id": 2,
    "last_seen_at": "2026-06-06T17:55:00Z"
  }
  ```

#### Partner Started Typing
* **Event**: `typing.start`
* **Data**:
  ```json
  {
    "conversation_id": 1,
    "user_id": 2
  }
  ```

#### Partner Stopped Typing
* **Event**: `typing.stop`
* **Data**:
  ```json
  {
    "conversation_id": 1,
    "user_id": 2
  }
  ```

#### Message Received
* **Event**: `message.received`
* **Data**: Message Object
  ```json
  {
    "id": 16,
    "conversation_id": 1,
    "sender_id": 1,
    "content": "Hi there Bob!",
    "sent_at": "2026-06-06T17:50:00Z"
  }
  ```

#### Message Read Receipt Update
* **Event**: `message.read`
* **Data**:
  ```json
  {
    "conversation_id": 1,
    "user_id": 1,
    "last_read_message_id": 15
  }
  ```
