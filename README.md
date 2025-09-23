## golang-realtime-chat

Simple realtime chat built with a Go WebSocket backend and a React frontend. No database — all state is kept in memory for teaching/demo purposes.

### What’s inside

- Go backend (`apps/backend`): HTTP server + WebSocket chat hub (join/leave, broadcast, user IDs)
- Web frontend (`apps/web`): Minimal chat UI (connect, send, users list)

### Prerequisites

- Go (1.21+ recommended)
- Bun or Node.js (for the web app)

### Run the backend

```bash
cd apps/backend
go mod tidy
go run .
# backend runs on http://localhost:8080
```

### Configure and run the frontend

```bash
cd apps/web
cp .env.example .env
# edit .env and set the backend URL (defaults shown)
# VITE_BACKEND_URL=http://localhost:8080

bun install   # or npm install / pnpm install / yarn
bun dev       # or npm run dev / pnpm dev / yarn dev
# frontend runs on http://localhost:3001 (per your dev setup)
```

### How to use

1. Start the backend, then the frontend
2. Open the web app in multiple tabs
3. Enter a username and Connect
4. Send messages; see join/leave and broadcasts in realtime

### Backend endpoints

- `GET /healthz` — health check
- `GET /api/users` — list of connected users and IDs
- `GET /api/stats` — basic server stats
- `WS /ws?username=YourName` — WebSocket chat endpoint

### Concepts demonstrated (for workshops)

- Explicit error handling in Go (`value, err`)
- Goroutines + channels for concurrency
- WebSocket lifecycle (upgrade, read/write pumps, heartbeats)
- In‑memory hub pattern for broadcast and membership

### Notes

- IDs are simple incrementing numbers for clarity in demos
- All state is ephemeral; restarting the backend resets the room
