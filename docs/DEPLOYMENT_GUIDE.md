# Deployment Guide

AntiquaOS Backend is designed to run inside Linux environments (like Hostinger VPS, AWS EC2, or DigitalOcean) via a compiled Go binary or Docker container.

## Standard Deployment
1. Set up a PostgreSQL 15+ database instance.
2. Clone the repository and install dependencies (`go mod tidy`).
3. Set your environment variables (DB credentials, external port, secret keys) matching what's in your local `.env`.
4. Compile the application into an executable binary:
   ```bash
   go build -o antiquaos-server main.go
   ```
5. Configure your domain and reverse proxy (e.g., NGINX acting as a proxy pass to your local go port)
6. Run using `systemd` or a screen session: `./antiquaos-server`

## Note on WebSockets
If routing via NGINX, make sure proxy settings allow upgrade requests `Connection "upgrade"` to keep WebSocket lines open.