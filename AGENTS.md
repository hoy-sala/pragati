# Pragati — AI Agent Instructions

Pragati is a school management system aligned with India's NEP 2020 (National Education Policy). It covers academics, assessments (CCE pattern), holistic progress cards (HPC), quizzes, mentor tracking, and certificates. The system serves admin, principal, teacher, special educator, student, and parent roles.

---

## Architecture

```
Browser → http://localhost:5050 (SvelteKit SSR/SPA)
              ↓ client-side fetch
         http://localhost:9090/api/v1 (Go REST API)
              ↓ pgx
         PostgreSQL :5432 + Redis :6379
```

| Layer | Stack |
|-------|-------|
| Frontend | SvelteKit 2 + Svelte 5 (runes), TypeScript, Tailwind CSS 3, Vite 5, adapter-node |
| Backend | Go 1.23, chi router, pgx (PostgreSQL), JWT auth, zerolog |
| Database | PostgreSQL 16 (pgcrypto, uuid-ossp, pg_trgm extensions) |
| Cache | Redis 7 |
| Deploy | Docker Compose, 4 services: web, api, db, redis |

All API routes live under `/api/v1`. The frontend SSR calls the API via Docker internal DNS (`http://api:9090`); browser-side calls use `http://{hostname}:9090`.

---

## Project Structure

```
pragati/
├── backend/
│   ├── cmd/server/main.go          # Entry point, seed admin user
│   ├── internal/
│   │   ├── auth/jwt.go             # JWT gen/validate, bcrypt
│   │   ├── config/config.go        # Env var loading
│   │   ├── cce/cce.go              # CCE grading logic (classes 6-10)
│   │   ├── database/db.go          # pgx pool, auto-migration runner
│   │   ├── handlers/               # HTTP handlers (14 files, ~80 endpoints)
│   │   │   ├── routes.go           # Route definitions + RBAC
│   │   │   ├── auth.go, students.go, classes.go, subjects.go,
│   │   │   │   assessments.go, marks.go, questions.go, quizzes.go,
│   │   │   │   dashboard.go, hpc.go, mentors.go, reports.go,
│   │   │   │   certificates.go, users.go
│   │   ├── middleware/auth.go       # Authenticate, RequireRole, RateLimiter
│   │   ├── models/                 # Go structs (10 files)
│   │   └── services/question_import/
│   │       ├── gift.go             # Math-aware GIFT parser (state machine)
│   │       ├── csv.go              # CSV question importer
│   │       └── *_test.go           # 16+ tests
│   ├── migrations/                 # 11 numbered migration pairs (.up/.down.sql)
│   ├── seeds/staff.sql             # 12 pre-seeded staff users
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── app.html, app.css
│   │   ├── lib/
│   │   │   ├── api/client.svelte.ts   # API client (token mgmt, auto-refresh)
│   │   │   ├── stores/auth.svelte.ts  # Auth state (Svelte 5 runes)
│   │   │   ├── types/index.ts         # All TypeScript interfaces
│   │   │   ├── utils/questionUtils.ts # Question type labels/colors
│   │   │   └── components/
│   │   │       ├── Button.svelte, Select.svelte, Pagination.svelte,
│   │   │       │   SearchFilter.svelte, MathText.svelte
│   │   │       └── layout/Sidebar.svelte, Search.svelte
│   │   └── routes/                # 30+ route pages (see below)
│   ├── static/fonts/, logos/
│   ├── package.json
│   ├── svelte.config.js
│   └── tailwind.config.js
├── docker-compose.yml
├── docker/Dockerfile.api, Dockerfile.web
├── .env, .env.example
└── scripts/, deploy.sh, deploy.mjs
```

---

## Build & Run Commands

### Backend (Go)

```bash
cd backend
go build ./cmd/server          # Build binary
go test ./...                  # Run all tests (16+ test functions)
go vet ./...                   # Lint
```

### Frontend (SvelteKit)

```bash
cd frontend
npm install
npm run dev                    # Dev server (Vite, port 5173)
npm run build                  # Production build → build/
npm run check                  # svelte-check (TypeScript)
npm run lint                   # Prettier
npm run test                   # Vitest
```

**Known svelte-check issues:**
- `Pagination.svelte:54` — pre-existing TS error (string|number not assignable to number)
- `Sidebar.svelte:163` — quoted attribute warning
- `mentors/logs/+page.svelte:89,93` — label association warnings
- `settings/+page.svelte:311,314` — dialog role + label warnings

### Docker

```bash
docker compose build api web   # Build containers
docker compose up -d api web   # Start services
docker compose exec api /server seed  # Seed database
```

### Database (direct access)

```bash
# Pipe SQL to psql via Docker
Get-Content file.sql | docker compose exec -T db psql -U pragati -d pragati
```

---

## Backend Conventions

### Handler Pattern

Every handler file follows this pattern:

```go
type FooHandler struct {
    db *pgxpool.Pool
}

func NewFooHandler(db *pgxpool.Pool) *FooHandler {
    return &FooHandler{db: db}
}

func (h *FooHandler) List(w http.ResponseWriter, r *http.Request) {
    // Parse query params with r.URL.Query().Get("key")
    // Query database with h.db.Query/QueryRow
    // Respond with respond(w, http.StatusOK, data)
}
```

**Response helpers** (in `handlers/respond.go` or similar):
- `respond(w, status, data)` — wraps in `{ data: ..., meta: ..., error: null }`
- `respondError(w, status, code, message)` — wraps in `{ data: null, error: { code, message } }`

### Route Definitions

All routes defined in `handlers/routes.go`. Pattern:

```go
r.Route("/foo", func(r chi.Router) {
    r.Use(roleMw.Authenticate)
    r.Get("/", h.List)                                    // any authenticated
    r.Post("/", roleMw.RequireRole("admin")(http.HandlerFunc(h.Create)))  // role-restricted
})
```

### Middleware

- `roleMw.Authenticate` — extracts JWT, sets user info in context
- `roleMw.RequireRole("admin", "principal")` — RBAC gate
- `middleware.NewRateLimiter(10, time.Minute)` — 10 req/min for login endpoints

### Database

- Auto-migration on startup: reads `*.up.sql` from `migrations/`, tracks in `schema_migrations`
- UUIDs everywhere (pgcrypto `gen_random_uuid()`)
- Soft deletes: `deleted_at TIMESTAMPTZ` with partial indexes `WHERE deleted_at IS NULL`
- JSONB columns for flexible data (HPC configs, question options, assessment chapters)
- Academic year scoping: most tables have `academic_year_id` FK

### Auth

Three login methods:
1. `POST /api/v1/auth/login` — email + password (admin/staff)
2. `POST /api/v1/auth/staff-login` — mobile + password (staff)
3. `POST /api/v1/auth/student-login` — SATS number + DOB (students, passwordless)

JWT claims: `user_id`, `school_id`, `role`, `email`, `sats_number`, `token_type`, `permissions`. Access tokens expire in 15m, refresh tokens in 7 days.

### Question Import (GIFT)

`services/question_import/gift.go` is a state-machine parser:
- Math zones `\(...\)` / `\[...\]` are opaque (no GIFT escaping inside)
- `\ce{}` for chemistry (KaTeX mhchem)
- Metadata parsed: `::title::`, `{#marks}`, `[difficulty:easy|medium|hard]`, `[chapter:Name]`, `[tags:a,b]`
- MCQ: `=correct` / `~wrong` markers
- Tests: `go test ./internal/services/question_import/ -v`

---

## Frontend Conventions

### Svelte 5 Runes

All components use Svelte 5 runes (NOT traditional stores):

```svelte
<script lang="ts">
  let { prop1, prop2 = 'default' }: { prop1: string; prop2?: string } = $props();
  let internal = $state(initial);
  let derived = $derived(expression);
  $effect(() => { /* side effects */ });
</script>
```

State files use `.svelte.ts` extension (e.g., `auth.svelte.ts`, `client.svelte.ts`).

### Component Library

Reusable components in `$lib/components/`:
- **Button** — variants: primary/secondary/danger/ghost, sizes: sm/md/lg, loading state
- **Select** — custom dropdown with portal, searchable, clearable
- **Pagination** — client-side pagination with ellipsis
- **SearchFilter** — debounced (300ms) search input
- **MathText** — KaTeX inline (`\(...\)`) + display (`\[...\]`) + mhchem

### API Client

```typescript
import { api, apiUrl } from '$lib/api/client.svelte';

// Standard call (auto-adds auth header, handles 401 refresh)
const res = await api<T>('GET', '/students');
const res = await api<T>('POST', '/students', { first_name: '...' });

// Direct URL (for file uploads/downloads)
const url = apiUrl('/uploads/signatures/abc.png');
```

Returns `{ data?: T; meta?: { pagination }; error?: { code; message } }`.

### Route Pages

| Route | Purpose |
|-------|---------|
| `/login` | 3-tab login (student/staff/admin) |
| `/dashboard` | Student: quiz stats, report card. Staff: stats, charts |
| `/classes`, `/subjects` | CRUD for admin |
| `/students` | Full CRUD with search, pagination |
| `/assessments` [+ `/create`] | Assessment list + creation |
| `/marks` | Tabulator data grid for marks entry (paste, validate, batch save) |
| `/questions` [+ `/create`, `/import`] | Question bank, create (MCQ/T-F/fill/short), GIFT/CSV import |
| `/quizzes` [+ `/create`, `/available`, `/[id]/edit`, `/take/[id]`, `/results/[attemptId]`] | Quiz lifecycle |
| `/analytics` | 3 tabs: overview, class performance, marks progress |
| `/reports` | Mark sheet, student report card (print), mentor-wise report |
| `/settings` | Admin: users, years, classes, subjects, categories |
| `/timetable` | Static display (class-wise + subject-wise views) |
| `/hpc` [+ `/config`, `/assess`, `/entry/[student_id]`, `/lo-import`] | Holistic progress card (NEP 2020) |
| `/mentors` [+ `/roster`, `/dashboard`, `/attendance`, `/logs`] | Mentor-student management |
| `/certificates` [+ `/print/[id]`] | Certificate events, signatories, print |

### Styling

- **Tailwind CSS 3** utility-first (no component library like shadcn)
- Custom CSS in `app.css` for fonts (Kannada: Noto Sans Kannada, code: JetBrains Mono)
- Print styles: `@media print` with `report-print` / `cert-print` class scoping
- Sidebar active state: longest-prefix-wins (prevents `/quizzes` + `/quizzes/available` both highlighting)

### Sidebar Navigation

Role-filtered via `item.roles.includes(effectiveRole)`. Active state uses longest-prefix matching against all nav hrefs — not naive `startsWith`.

---

## Database Schema (Key Entities)

### Core Relationships

```
School → Academic Year → Class → Section → Student
                     ↓
              Class ↔ Subject (many-to-many)
              Teacher ↔ Subject ↔ Class (many-to-many)

Assessment → Category (FA1-FA4, SA1-SA2) → Marks (per student)
Question ↔ Quiz (many-to-many) → Attempt → Response

Mentor(User) ↔ Student (assignment) → Attendance, Logs

Certificate Event → Certificate (per student), Signatories
```

### Seed Data

- Admin: `admin@pragati.edu` / `pragati123` (auto-created)
- 12 staff: pre-seeded via `seeds/staff.sql` (password: `pragati123`)
- Assessment categories: FA1(10%), FA2(10%), FA3(10%), FA4(10%), SA1(20%), SA2(20%)
- Default school ID: `00000000-0000-0000-0000-000000000001`

### CCE Grading (Classes 6-8)

FA: 4 tests × 30 raw → converted to 40. SA: 2 tests × 50 raw → converted to 60. Total = 100.

### HPC Proficiency Scale

1=Beginning, 2=Developing, 3=Proficient, 4=Advanced

---

## RBAC Summary

| Action | admin | principal | teacher | special_educator | student |
|--------|-------|-----------|---------|------------------|---------|
| User CRUD | ✓ | | | | |
| Student CRUD | ✓ | ✓ | ✓ | | |
| Student Delete | ✓ | | | | |
| Assessment CRUD | ✓ | ✓ | ✓ | | |
| Assessment Publish | ✓ | ✓ | | | |
| Marks Entry | ✓ | ✓ | ✓ | | |
| Question CRUD/Import | ✓ | ✓ | ✓ | | |
| Quiz CRUD | ✓ | ✓ | ✓ | | |
| Quiz Attempt | ✓ | ✓ | ✓ | ✓ | ✓ |
| HPC Config | ✓ | | | | |
| HPC Entry/Assess | ✓ | ✓ | ✓ | | |
| Mentor Assign | ✓ | | | | |
| Mentor Attendance/Logs | ✓ | ✓ | ✓ | ✓ | ✓ |
| Certificate CRUD | ✓ | | | | |
| Reports | ✓ | ✓ | ✓ | ✓ | ✓ |

---

## Deployment

### Production Server

App directory: `/opt/apps/pragati`. SSH access via Tailscale or direct.

```bash
# Build and deploy
docker compose build api web
docker compose up -d api web

# Verify
curl http://localhost:5050/    # web → 200
curl http://localhost:9090/health  # api → {"status":"ok"}
```

### Port Mappings

| Service | Container | Host |
|---------|-----------|------|
| web (SvelteKit) | :3000 | **:5050** |
| api (Go) | :9090 | **:9090** |
| db (PostgreSQL) | :5432 | :5432 |
| redis | :6379 | :6379 |

### CI

GitHub Actions (`.github/workflows/ci.yml`): runs on push/PR to `main`.
- Backend: `go vet` + `go build`
- Frontend: `npm ci` + `npm run build`

No CD pipeline — deployment is manual via scripts or SSH.

---

## Key Files Reference

| File | Purpose |
|------|---------|
| `backend/internal/handlers/routes.go` | All API routes + RBAC |
| `backend/internal/database/db.go` | pgx pool + auto-migration |
| `backend/internal/auth/jwt.go` | JWT + password hashing |
| `backend/internal/cce/cce.go` | CCE grading logic |
| `backend/internal/services/question_import/gift.go` | GIFT parser (math-aware) |
| `backend/internal/models/*.go` | All Go structs |
| `backend/migrations/*.sql` | Database schema |
| `frontend/src/lib/api/client.svelte.ts` | API client |
| `frontend/src/lib/stores/auth.svelte.ts` | Auth state |
| `frontend/src/lib/types/index.ts` | TypeScript interfaces |
| `frontend/src/lib/components/MathText.svelte` | KaTeX rendering |
| `frontend/src/routes/+layout.svelte` | Root layout (auth gate, sidebar) |
| `docker-compose.yml` | Service orchestration |
| `.env.example` | Environment variable template |

---

## Gotchas & Notes

1. **Soft deletes everywhere** — queries must filter `WHERE deleted_at IS NULL` (partial indexes handle this, but raw SQL needs it)
2. **Academic year scoping** — most data is scoped to `academic_year_id`; always filter by current year
3. **UUIDs, not integers** — all primary keys are UUIDs (`gen_random_uuid()`)
4. **JSONB for flexibility** — HPC configs, question options, assessment chapters all use JSONB columns
5. **quiz_attempts.user_id FK dropped** — migration 007 intentionally drops this FK because attempts can reference either `users` or `students`
6. **KaTeX requires `\(...\)` / `\[...\]` delimiters** — NOT `$...$` or `$$...$$`. No escaping of `~` `=` `{` `}` inside math zones
7. **Copy button on HTTP** — clipboard API unavailable on plain HTTP; uses `execCommand('copy')` fallback
8. **svelte-check has 1 pre-existing error** — Pagination.svelte:54 TS issue; all warnings are non-blocking
9. **Git push may fail if GitHub unreachable** — use `git bundle` + plink stdin transfer as fallback
10. **Sidebar active state** — uses longest-prefix-wins, not naive `startsWith`, to prevent parent/child nav items both highlighting
