# Anonymous Image feed

Anonymous image upload site with a live feed. React + TypeScript frontend, Go + Chi backend.

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

Edit `.env` first if you need different ports or Postgres credentials (defaults work out of the box).

### Option B: Run services locally

Prerequisites: Go 1.25+, Node 24, Postgres, and [golang-migrate](https://github.com/golang-migrate/migrate).

1. Start Postgres and create a database (or run `docker compose up db` to use the one from Compose).
2. Apply migrations:
   ```
   migrate -path backend/migrations -database "postgres://imagefeed:imagefeed@localhost:5432/imagefeed?sslmode=disable" up
   ```
3. Run the backend (reads `DATABASE_URL`, defaults to the connection string above):
   ```
   cd backend
   go run ./cmd/server
   ```
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