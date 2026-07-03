# Fast API

Fast API is a customized fork of the upstream AI gateway project. This repository keeps the upstream remote for future merges, while local branding and deployment files are maintained in the `my-custom` branch.

## Local Development

Requirements:

- Go 1.25+
- Node.js 22+
- npm or Bun
- Docker is optional for local development, but required for image builds

Initial setup:

```bash
git remote -v
git checkout my-custom
cp .env.example .env
go mod download
```

Run the backend:

```bash
go run main.go --port 3000 --log-dir ./logs
```

Run the default frontend:

```bash
cd web/default
npm install
npm run dev
```

Useful environment variables are documented in `.env.example`. The default local database uses SQLite; production can also use MySQL or PostgreSQL via `SQL_DSN`.

## IntelliJ IDEA

The repository includes shared run configurations:

- `Backend: Fast API` runs `main.go` on port `3000`
- `Frontend: npm run dev` runs the default frontend dev server

Open the project folder in IntelliJ IDEA and make sure the Go plugin uses a local Go SDK compatible with `go.mod`.

## Rebranding After Upstream Merges

After merging upstream changes, rerun:

```bash
bash scripts/rebrand.sh
```

On Windows, use Git Bash:

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/rebrand.sh
```

The script is idempotent and can be run repeatedly.

## Docker Image

Build a production image for the `/fast-api/` public path:

```bash
docker build \
  --build-arg WEB_BASE_PATH=/fast-api/ \
  -t fast-api:latest .
```

Run a single container with persistent data:

```bash
mkdir -p /data/fast-api
cp .env.example /data/fast-api/.env
docker run -d \
  --name fast-api \
  --restart always \
  --env-file /data/fast-api/.env \
  -p 127.0.0.1:3001:3000 \
  -v /data/fast-api:/data \
  fast-api:latest
```

## GitHub Actions Deployment

`.github/workflows/deploy.yml` builds and pushes:

```text
<DOCKER_USERNAME>/fast-api:latest
```

Required GitHub repository secrets:

| Secret | Meaning |
| --- | --- |
| `DOCKER_USERNAME` | Docker Hub username |
| `DOCKER_PASSWORD` | Docker Hub password or access token |
| `SERVER_HOST` | Linux server public IP or hostname |
| `SERVER_USER` | SSH user, for example `root` |
| `SERVER_SSH_KEY` | Complete PEM private key authorized on the server |

Optional repository variables:

| Variable | Default | Meaning |
| --- | --- | --- |
| `FAST_API_HOST_PORT` | `3001` | Host loopback port mapped to container port `3000` |
| `FAST_API_BASE_PATH` | `/fast-api/` | Public frontend base path |

Deployment runs only after the image build and push steps succeed.

## Server Nginx

Use `nginx/fast-api.conf` as the reference reverse proxy. It proxies:

- `https://www.fastapi.cool/fast-api/` to the Fast API web UI
- `/api/`, `/v1/`, `/mj/`, `/pg/`, and `/assets/` to the same container

The config includes WebSocket/SSE-friendly settings such as `proxy_buffering off` and `X-Accel-Buffering no`.

## Manual Server Checklist

```bash
# 1. Install Docker if needed.
curl -fsSL https://get.docker.com | sh
systemctl enable --now docker

# 2. Prepare persistent data and environment.
mkdir -p /data/fast-api
cp .env.example /data/fast-api/.env
vi /data/fast-api/.env

# 3. Pull and run the image built by CI.
docker pull <DOCKER_USERNAME>/fast-api:latest
docker rm -f fast-api 2>/dev/null || true
docker run -d \
  --name fast-api \
  --restart always \
  --env-file /data/fast-api/.env \
  -p 127.0.0.1:3001:3000 \
  -v /data/fast-api:/data \
  <DOCKER_USERNAME>/fast-api:latest

# 4. Install the Nginx config, then reload.
nginx -t
systemctl reload nginx

# 5. Verify.
curl -I http://127.0.0.1:3001/
docker logs --tail 100 fast-api
```

Official deployment reference: https://docs.newapi.pro/zh/docs/installation
