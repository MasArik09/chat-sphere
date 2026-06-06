# Release Notes - ChatSphere V1.0.0

We are proud to announce the initial official stable release of **ChatSphere V1.0.0**, a robust, real-time messaging application and system setup.

---

## 1. Major Features Delivered

* **Full Authentication Flow**: Secure user registration, password hashing (bcrypt), login validation, and JWT authentication.
* **Instant Direct Messaging**: High-performance persistent WebSocket connectivity supporting bidirectional chat synchronization.
* **Online Presence Tracking**: Real-time status sync (online/offline) and last-seen timestamp updates automatically propagated to partners.
* **Typing Indicator Hooks**: Start/stop typing states showing real-time text input actions.
* **Read Receipts Syncing**: Database-level and WS-level updates mapping user last-read message IDs inside conversation threads.
* **Conversation Search**: Fast filtering capabilities based on participant profiles.
* **Production Docker Orchestration**: Secure Compose configurations featuring GIN production release mode, persistent named database volumes, and strict container logging limits.
* **Security Rate Limiting**: Production Nginx reverse proxy gateway implementing rate throttling on registration and login endpoints to secure against credential stuffing.

---

## 2. Technical and Architecture Highlights

* **N+1 SQL Optimizations**: List conversations query utilizes SQL aggregations and lateral joins to load last messages, user profiles, and unread counts in exactly 1 database roundtrip.
* **Transaction Rollback Guarantees**: Core database creations are wrapped under transactional operations. If any action fails (e.g. participant association), all writes are rolled back immediately to prevent database corruption.
* **Gateway Level Isolation**: Services are unexposed to the host/public IP by default. Backend port `8080` is restricted to `127.0.0.1` (accessible only to local Nginx proxy), and Postgres port `5432` remains unexposed.
* **Split Health Checks**: Native `/health/live` (fast process checks) and `/health/ready` (active-pings database connection) split handlers.

---

## 3. Production Readiness Summary

* **Build Validation**: Backend compiles statically without errors. Frontend compiles production assets successfully using Vite.
* **Tests Pass Rate**: 100% test pass rate across unit tests and integration tests.
* **Repository Hygiene**: Confirmed no secrets are committed to the repository. Added `backups/`, `temp.sql`, and local files to `.gitignore`.
* **Logging Safety**: Capped container logs (`max-size: 10m`, `max-file: 3`) to prevent host system disk depletion.

---

## 4. Known Post-MVP Limitations

* **No Group Chats**: V1.0.0 only supports private direct messaging (1-on-1 chats).
* **Relative Timestamps**: Frontend displays exact timestamp strings; formatting relative times (e.g. "2 minutes ago") is left for frontend localizations.
* **No Media Messages**: Only raw text payloads are supported over WebSocket and REST endpoints.
* **Single Device WebSocket**: WebSocket connections do not implement connection limits or device-multiplexing; multiple connections from the same user ID are broadcasted concurrently.

---

## 5. Future Roadmap

* **Phase 11: Group Conversations**: Implement database and API structures to support conversations with 3+ participants.
* **Phase 12: Rich Media Attachments**: Extend message schema and storage layers to support image and file transmissions.
* **Phase 13: Push Notifications**: Implement standard Web Push API notifications to alert users who have closed their browser connections.
