# Docker Rules

# ChatSphere V1

Version: 1.0

Status: Mandatory

---

# 1. Purpose

This document defines all Docker-related rules.

Every containerized component must follow these standards.

---

# 2. Docker Philosophy

The application must run using:

```bash
docker compose up
```

A fresh clone should be runnable with minimal setup.

---

# 3. Required Services

Docker Compose must contain:

```text
frontend
backend
postgres
```

No additional services are required in V1.

---

# 4. Service Naming

Required service names:

```text
frontend
backend
postgres
```

Do not rename without approval.

---

# 5. Container Communication

Containers communicate through Docker Network.

Examples:

Correct:

backend → postgres

frontend → backend

Incorrect:

backend → localhost

frontend → localhost

Never use localhost between containers.

---

# 6. Environment Variables

Backend must use:

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chatsphere
```

Frontend must use:

```env
VITE_API_URL=http://backend:8080
VITE_WS_URL=ws://backend:8080/ws
```

---

# 7. Database Container

Database:

PostgreSQL

Required:

- Persistent volume
- Health check
- Environment variables

---

# 8. Backend Container

Responsibilities:

- Run Go API
- Connect PostgreSQL
- Serve WebSocket

Default Port:

```text
8080
```

---

# 9. Frontend Container

Responsibilities:

- Run React Application
- Communicate with Backend

Default Port:

```text
5173
```

---

# 10. Dockerfile Rules

Each service must have its own Dockerfile.

Required:

```text
frontend/Dockerfile

backend/Dockerfile
```

---

# 11. Build Rules

Docker images must be reproducible.

Avoid:

- Hardcoded paths
- Local machine dependencies

---

# 12. Volumes

Database must persist data.

Example:

```yaml
volumes:
  postgres_data:
```

Required for PostgreSQL.

---

# 13. Health Checks

PostgreSQL:

Must expose healthy state.

Backend:

Must expose health endpoint.

Example:

```text
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

---

# 14. Startup Order

Required sequence:

```text
PostgreSQL

↓

Backend

↓

Frontend
```

Use:

depends_on

where appropriate.

---

# 15. Development Workflow

Typical workflow:

```bash
docker compose up

docker compose down

docker compose logs

docker compose build
```

---

# 16. Forbidden Practices

Do NOT:

- Use localhost inside containers
- Store secrets in Dockerfiles
- Hardcode credentials
- Depend on host machine paths

---

# 17. Verification Checklist

Verify:

✓ frontend running

✓ backend running

✓ postgres running

✓ database connection successful

✓ websocket connection successful

✓ API accessible

✓ data persisted after restart

---

# 18. Definition of Done

Docker setup is complete when:

✓ docker compose up works

✓ Fresh clone works

✓ Database persists data

✓ WebSocket works

✓ No manual container fixes required