# Coding Standards

# ChatSphere V1

Version: 1.0

Status: Mandatory

---

# 1. Purpose

This document defines coding standards for ChatSphere.

Every source file must follow these standards.

Code quality is more important than development speed.

---

# 2. General Principles

Code must be:

- Readable
- Maintainable
- Testable
- Predictable

Always optimize for future maintenance.

---

# 3. File Size Limits

## Go Files

Target:

```text
< 200 lines
```

Hard Limit:

```text
300 lines
```

---

## React Components

Target:

```text
< 150 lines
```

Hard Limit:

```text
250 lines
```

---

## Service Files

Target:

```text
< 200 lines
```

Hard Limit:

```text
300 lines
```

---

If limits are exceeded:

Refactor immediately.

---

# 4. Single Responsibility Rule

Every file must have one responsibility.

Good:

```text
AuthService
ConversationService
MessageRepository
```

Bad:

```text
UserConversationMessageService
```

---

# 5. Naming Conventions

Use clear names.

Good:

```go
CreateConversation

GetUserByID

StoreMessage
```

Bad:

```go
DoThing

Process

HandleStuff
```

---

# 6. Go Naming Rules

Exported:

```go
CreateConversation
MessageService
UserRepository
```

Private:

```go
createConversation
validateMessage
```

---

# 7. React Naming Rules

Components:

```tsx
LoginForm.tsx

ConversationList.tsx

MessageBubble.tsx
```

Hooks:

```tsx
useAuth.ts

useWebSocket.ts

useConversations.ts
```

---

# 8. Folder Rules

Never place everything in one folder.

Use feature-based structure.

Correct:

```text
features/
├── auth/
├── users/
├── conversations/
├── messages/
└── presence/
```

Incorrect:

```text
components/
    200 files
```

---

# 9. Handler Rules

Handlers only:

- Receive requests
- Call services
- Return responses

Handlers must NOT:

- Execute SQL
- Perform business logic
- Manage websocket state

---

# 10. Service Rules

Services contain:

- Business rules
- Validation logic
- Domain logic

Services must NOT:

- Return HTTP responses
- Use Gin Context
- Query database directly

---

# 11. Repository Rules

Repositories:

- Access PostgreSQL
- Execute queries
- Return entities

Repositories must NOT:

- Contain business logic

---

# 12. WebSocket Rules

WebSocket handlers must:

- Authenticate
- Register clients
- Forward events

WebSocket handlers must NOT:

- Query database
- Perform business logic

---

# 13. React Component Rules

Components should:

- Focus on UI
- Receive props
- Remain small

Avoid:

- Massive state logic
- API calls inside many components

Prefer:

- Hooks
- Services

---

# 14. Custom Hook Rules

Business logic belongs in hooks.

Examples:

```tsx
useAuth()

useMessages()

useWebSocket()
```

Hooks should:

- Encapsulate logic
- Be reusable

---

# 15. API Service Rules

Frontend API calls belong in:

```text
src/services/
```

Example:

```ts
authService.ts

conversationService.ts

messageService.ts
```

Do not scatter fetch calls across components.

---

# 16. TypeScript Rules

Avoid:

```ts
any
```

Prefer:

```ts
interfaces

types
```

Example:

```ts
interface User {
  id: number;
  name: string;
}
```

---

# 17. Error Handling Rules

Always handle errors.

Backend:

```go
if err != nil {
    return err
}
```

Frontend:

```ts
try {
}
catch {
}
```

Never silently ignore errors.

---

# 18. Logging Rules

Allowed:

```text
Application Logs

Error Logs

Connection Logs
```

Forbidden:

```text
Passwords

JWT Tokens

Database Credentials
```

---

# 19. Comments

Comment WHY.

Do not comment WHAT.

Bad:

```go
// increment i
i++
```

Good:

```go
// Prevent duplicate conversation creation.
```

---

# 20. Magic Values

Avoid:

```go
if role == 3
```

Prefer:

```go
const AdminRole = 3
```

---

# 21. Validation Rules

All external input must be validated.

Sources:

- HTTP Requests
- WebSocket Events
- Query Parameters

Never trust user input.

---

# 22. Security Rules

Always verify:

- JWT ownership
- Conversation ownership
- Message ownership

Prevent:

- IDOR
- Unauthorized access

---

# 23. Database Rules

Always:

- Use indexes
- Use foreign keys
- Use transactions when required

Avoid:

- N+1 queries
- Full table scans

---

# 24. Testing Rules

Every feature requires tests.

Minimum:

Authentication

Conversations

Messages

Presence

WebSocket

---

# 25. Forbidden Anti-Patterns

Do NOT create:

God Handler

God Service

God Repository

God Component

Massive Utility File

---

Example Bad:

```text
ChatService.go

1200 lines
```

---

# 26. Refactoring Trigger Rules

Refactor when:

- File > 300 lines
- Function > 50 lines
- Component too complex
- Duplicate code appears 3+ times

---

# 27. Pull Request Mentality

Every phase should be:

- Small
- Verifiable
- Reversible

Never implement multiple phases together.

---

# 28. Documentation Rules

If architecture changes:

Update documentation.

Required:

- SRS
- System Design
- Database Design

---

# 29. Code Review Checklist

Before completion verify:

✓ Naming clear

✓ Responsibilities clear

✓ Tests exist

✓ No duplicated logic

✓ No security issue

✓ No large files

✓ Documentation updated

---

# 30. Definition of Good Code

Good code is:

- Easy to read
- Easy to test
- Easy to modify
- Easy to extend

Future developers should understand the code without explanation.