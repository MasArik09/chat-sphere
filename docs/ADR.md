# Architectural Decision Records (ADR)

This document records major technology stack and architectural decisions made during the design and development of ChatSphere V1.

---

## ADR 001: Backend Service Engine — Go + Gin

### Status
Approved

### Context
We needed a backend technology stack capable of handling concurrent network connections (REST API requests and persistent WebSockets) with low memory footprints and fast execution speeds.

### Decision
We chose **Go (Golang)** as the core programming language, paired with the **Gin Gonic** web framework.

### Consequences
- **Concurrency Model**: Go's native goroutines make it highly suited to WebSocket orchestration where thousands of clients maintain active socket connections.
- **Single Static Binary**: Compiling backend code into a single lightweight Docker container (`alpine:3.19` matching ca-certificates and tzdata) results in ultra-fast boot speeds and small deployment foot-prints (~20MB image size).
- **Gin Framework**: Offers an extremely fast router and simple middleware capabilities (CORS, Recover, Logging) without the weight of an enterprise framework.

---

## ADR 002: Relational Database — PostgreSQL

### Status
Approved

### Context
ChatSphere requires strong schema validation, composite indexing, transaction isolation (ACID), and robust aggregation to manage conversations, user memberships, and messages safely.

### Decision
We chose **PostgreSQL** as our primary relational database engine.

### Consequences
- **ACID Guarantees**: Allows us to wrap conversation creations and participant insertions inside single transactions, guaranteeing rollback safety if parts of the operation fail.
- **Relational Integrity**: Foreign key constraints and unique composite indexes (e.g. unique constraint on conversation and participant lists) prevent data duplicates and orphans.
- **JSON Aggregation**: Utilizes PostgreSQL's JSON building capability (`json_build_object` / `json_agg`) to aggregate conversation participant metadata directly at the database query level, preventing client-side N+1 roundtrips.

---

## ADR 003: Real-Time Event Layer — WebSockets

### Status
Approved

### Context
A messaging app requires real-time updates for typing indicators, online/offline presence tracking, read receipts, and message delivery. Poll-based HTTP approaches add high network overhead and server load.

### Decision
We chose **WebSockets** as the bidirectional real-time communications layer.

### Consequences
- **Persistent TCP Socket**: Replaces HTTP polling with a single persistent connection, cutting TCP handshake latency.
- ** gorilla/websocket Integration**: Go's Gorilla WebSocket library is highly standard and allows efficient client register/unregister and connection thread safety.
- **Low Network Overhead**: Instant notifications and state updates are sent as lightweight JSON packets directly over active sockets, maintaining ultra-low latency.

---

## ADR 004: Frontend State Manager — Zustand

### Status
Approved

### Context
The React frontend needs to react immediately to real-time events (incoming messages, updated unread counts, typing states, and online changes). Standard React Context triggers global re-renders and Redux introduces excessive boilerplate.

### Decision
We chose **Zustand** as the frontend state management framework.

### Consequences
- **Boilerplate Reduction**: Zustand stores are created using simple, readable hooks without action builders, actions, or dispatchers.
- **Targeted Re-renders**: Components select only the exact slices of state they need, reducing unnecessary DOM updates.
- **WebSocket Event Binding**: Store actions are directly bound to the WebSocket event handler, allowing incoming WebSocket payloads to update the UI state reactively.
