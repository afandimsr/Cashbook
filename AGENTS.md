# AGENTS.md

Guidance for AI coding agents working in this repository. Read this first, then read the detailed guide for the layer you're touching **before writing code**.

## What this is

**CashBook** — a full-stack personal finance manager.

- **Backend** (`/backend`): Go 1.25, Gin, PostgreSQL (`database/sql`), golang-migrate, JWT + bcrypt + TOTP. Clean Architecture.
- **Frontend** (`/frontend`): React 19 + TypeScript + Vite, MUI v7, Zustand, React Router v7, vite-plugin-pwa. Money is Rupiah (IDR).

## Where to read before coding (required)

| Working on… | Read |
|---|---|
| Anything backend | [`docs/agent/backend/`](docs/agent/backend/README.md) |
| Anything frontend | [`docs/agent/frontend/`](docs/agent/frontend/README.md) |
| The exact tech stack & commands | [`docs/agent/stack-profile.md`](docs/agent/stack-profile.md) |
| Stack-agnostic rules | [`docs/agent/principles/`](docs/agent/principles/) |
| How the docs are structured / templating them | [`docs/agent/README.md`](docs/agent/README.md) |

The guides are layered: **`principles/`** (stack-agnostic rules) + **`stack-profile.md`** (the chosen stack) + **`backend/` & `frontend/`** (this-stack implementation). Each layer README maps every principle to its implementation file. These encode real conventions and past bugs — treat them as rules, not suggestions.

**Work like a senior engineer.** Two principles apply to every change regardless of layer:
- [`docs/agent/principles/engineering-standards.md`](docs/agent/principles/engineering-standards.md) — mindset, minimal focused diffs, no over-engineering, edge/failure cases, honesty/verification, and the **Definition of Done**.
- [`docs/agent/principles/production-readiness.md`](docs/agent/principles/production-readiness.md) — observability, performance, reliability, config/secrets, migrations/rollback, API contracts, a11y — with a checklist to apply per change.

## Repository layout

```
backend/     Go API (internal/{domain,usecase,delivery,infrastructure,...}, migrations/, cmd/)
frontend/    React app (src/{app,application,domain,infrastructure,presentation,state})
docs/agent/  Agent guides (backend/ + frontend/)  ← the rules
docker-compose.yaml, deploy.sh, README.md
```

## Commands

Backend:
```bash
cd backend
go build ./... && go vet ./... && go test ./...   # verify
make migrate-up                                    # apply migrations
go run cmd/api/main.go                             # run server
```
Frontend:
```bash
cd frontend
npx tsc --noEmit      # fast type check
npm run build         # full build
npm run dev           # dev server
```
Full stack: `docker compose up -d --build backend frontend` then `docker compose run --rm migrate`.

## Top cross-cutting rules

**Backend**
1. Respect dependency direction: `delivery → usecase → domain`, `infrastructure → domain`. Business logic lives in `usecase`; handlers parse/respond; repos persist.
2. Usecases depend on `domain` interfaces (mockable). Parameterized SQL only. Multi-write ops use `WithTransaction`.
3. Return domain sentinel errors for validation (`errors.Is` → 400 in handler). Every schema change is an up/down migration pair.

**Frontend**
1. `apiClient` already unwraps `{data}` — **never access `.data` again** in stores (double-unwrap = empty UI bug).
2. All HTTP via `apiClient`; server data lives in Zustand stores; components read stores.
3. Every protected parent route needs an `index` element — `/` must redirect to `/dashboard` (missing this = blank screen). Format money with `formatIDR`. Keep the service worker simple.

## Definition of Done (apply to every change)

A change is done only when **all** hold (full checklists in the two principle files above):

- Requirement met; edge and failure cases handled; change is minimal and coherent.
- Backend: `go build ./... && go vet ./... && go test ./...` pass; tests added/updated for changed logic.
- Frontend: `npx tsc --noEmit` passes; async views have loading/error/empty states.
- Schema change ⇒ reversible migration pair + updated repo `SELECT`/`INSERT`, same change.
- No secrets/PII committed or logged; errors handled, not swallowed; client errors sanitized.
- Production-readiness checklist reviewed for the change's scope (observability, pagination/timeouts, atomicity, config, a11y).
- Tradeoffs/assumptions surfaced. Do not commit or push unless asked; branch off the default branch. Confirm before destructive or outward-facing actions.
