# Pragati (ಪ್ರಗತಿ)

**Every Child Can Progress.**

Student Assessment, Progress Tracking and Academic Intelligence Platform.

## Quick Start

```bash
# Clone and start all services
docker compose -f docker/docker-compose.yml up -d

# Run database migrations
docker compose exec api ./server          # auto-runs on startup

# Seed admin user
docker compose exec api ./server seed

# Open in browser
open http://localhost:5050
```

**Default admin:** `admin@pragati.edu` / `pragati123`

## Architecture

| Layer | Technology | Purpose |
|---|---|---|
| Frontend | SvelteKit 5 + Tailwind CSS 3 | SSR, SPA, mobile-responsive UI |
| Backend API | Go 1.23 + Chi router | REST API, JWT auth, business logic |
| Database | PostgreSQL 16 | Primary data store |
| Cache | Redis 7 | Session cache, rate limiting, pub/sub |
| Reverse Proxy | Caddy | Auto-HTTPS, static file serving |

## Project Structure

```
backend/          # Go API server
  cmd/server/     # Server entry point (also handles `seed` subcommand)
  internal/
    config/       # Configuration
    auth/         # JWT, password hashing
    database/     # PostgreSQL connection
    handlers/     # HTTP handlers by module
    middleware/   # Auth, logging, recovery
    models/       # Data types
  migrations/     # SQL migration files

frontend/         # SvelteKit application
  src/
    routes/       # Page components
    lib/
      api/        # API client
      stores/     # Auth state
      types/      # TypeScript types
      components/ # Reusable UI components

docker/           # Dockerfiles and compose
```

## Key Features (Phase 1)

- [x] Student CRUD with 9-digit SATS number
- [x] CSV import/export
- [x] Class, section, subject management
- [x] Academic year management
- [x] Role-based authentication (6 roles)
- [x] JWT access + refresh tokens
- [ ] Assessment management (Phase 2)
- [ ] Marks entry spreadsheet (Phase 2)
- [ ] Question bank (Phase 3)
- [ ] Online tests (Phase 3)
- [ ] Foundational assessments (Phase 3)
- [ ] Analytics (Phase 4)
- [ ] Reports & certificates (Phase 5)

## Development

```bash
# Backend (requires Go 1.23+)
cd backend
go run ./cmd/server
go run ./cmd/server seed  # Seed initial data

# Frontend (requires Node 22+)
cd frontend
npm run dev

# Both
docker compose -f docker/docker-compose.yml up --watch
```

## API

All endpoints are prefixed with `/api/v1`. See docs/api.md for full reference.

## License

Proprietary — All rights reserved.
