# ChatSphere V1 - Security Hardening Recommendations

This document details critical security recommendations for hardening the ChatSphere V1 application in production environments.

---

## 1. Network & Database Security

- **Database Port Protection**: The `postgres` database container in `docker-compose.prod.yml` does NOT expose its port `5432` to the host machine or public internet. This prevents database brute force attacks. All database interactions should go through the internal Docker bridge network `chatsphere-prod-network`.
- **Backend API Gateway Isolation**: The `backend` service is bound to `127.0.0.1:8080:8080`. This ensures it is accessible only by the Nginx proxy running on localhost, and cannot be reached directly via the public host IP, preventing bypass of the Nginx rate-limiting rules.
- **Firewall Setup (UFW)**: Configure `ufw` on your Ubuntu host server to only expose ports 80, 443, and your secure SSH port:
  ```bash
  sudo ufw default deny incoming
  sudo ufw default allow outgoing
  sudo ufw allow 80/tcp
  sudo ufw allow 443/tcp
  sudo ufw allow 22/tcp  # Or your secure SSH port
  sudo ufw enable
  ```

---

## 2. Nginx Security Hardening

The custom `nginx.conf` served in the frontend container has several security options pre-configured:
- **Auth Endpoint Rate Limiting**: Limit register and login endpoints to `20r/m` per IP with `burst=20 nodelay` using `limit_req_zone` to protect against credential stuffing and brute-force attacks.
- **HTTP Header Security**:
  - `X-Frame-Options: DENY` (prevents clickjacking).
  - `X-Content-Type-Options: nosniff` (mitigates MIME-type sniffing).
  - `X-XSS-Protection: 1; mode=block` (filters script injections).
  - `Content-Security-Policy` (CSP): Defines strict resource loading scopes (`default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self' ws: wss:;`).
  - `Referrer-Policy: strict-origin-when-cross-origin`.
- **Limit Body Sizes**: Set body size limits in your host Nginx to restrict large request payloads:
  ```nginx
  client_max_body_size 10M;
  ```

---

## 3. Docker Container Hardening

- **Non-Root Execution**: By default, official images like `nginx` and `postgres` run as dedicated system users (`nginx` and `postgres`). For the compiled Go backend container, consider modifying the Dockerfile to compile statically and run as a non-privileged user (e.g. `nobody` or a custom `app` user created in the alpine stage).
- **Logging Limits**: Ensure log rotations are active. We configure `max-size: "10m"` and `max-file: "3"` in `docker-compose.prod.yml` to prevent disk exhaustion vectors.

---

## 4. Secret Management & Rotation

- **Environment Secrets**: Never commit `.env.production` files containing actual passwords or secret keys. Use `.env.production.example` to document keys, and supply production values on the server.
- **JWT Key Strength**: Use a 256-bit or 512-bit key for `JWT_SECRET`.
- **Periodic Rotation**: Establish a schedule (e.g., every 90 days) to rotate database passwords and JWT secrets.
- **Postgres Password Strength**: Always use generated alphanumeric passwords containing at least 24 characters.
