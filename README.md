# Anonymous Image feed

Anonymous image upload site with a live feed. React + TypeScript frontend, Go + Chi backend.

## Features
- Anonymous image uploads (jpeg, png, webp), normalized server side to a standard size and JPEG quality
- Live feed updates over WebSocket when a new image is uploaded
- Cursor based pagination and tag filtering on the image feed
- Redis backed response caching for the images and tags endpoints

## Structure
- `backend/` — Go API (Chi router, REST + WebSocket)
- `frontend/` — React + TypeScript (Vite)

## Setup

### Option A: Docker Compose (recommended)

Prerequisite: Docker.

```
cp .env.example .env
docker compose up
```

This starts Postgres, runs migrations, and starts the backend and frontend. Once it's up:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- API docs (Swagger UI): http://localhost:8080/docs

Edit `.env` first if you need different ports or Postgres credentials (defaults work out of the box).

### Option B: Run services locally

Prerequisites: Go 1.25+, Node 24, Postgres, Redis, and [golang-migrate](https://github.com/golang-migrate/migrate).

1. Start Postgres and Redis (or run `docker compose up db redis` to use the ones from Compose).
2. Apply migrations:
   ```
   migrate -path backend/migrations -database "postgres://imagefeed:imagefeed@localhost:5432/imagefeed?sslmode=disable" up
   ```
3. Run the backend (reads `DATABASE_URL` and `REDIS_URL`, both default to localhost as shown below):
   ```
   cd backend
   DATABASE_URL="postgres://imagefeed:imagefeed@localhost:5432/imagefeed?sslmode=disable" \
   REDIS_URL="redis://localhost:6379" \
   go run ./cmd/server
   ```
   API docs (Swagger UI) are then available at http://localhost:8080/docs.
4. Run the frontend:
   ```
   cd frontend
   npm install
   npm run dev
   ```

## Testing

**Backend**
```
cd backend
go test ./...
```

**Frontend**

```
cd frontend
npm test
```