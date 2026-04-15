# DangDe! World Handover Status

## Context

This repository started from `requirements.md` and is now a greenfield implementation with:

- Backend in Go using Gin + Gorm with DDD-style structure
- Frontend in Vue 3 + Vite using a feature-oriented structure
- Seeded demo roles for admin, parent, and kids
- A custom activity runtime where activities can carry admin-authored HTML, CSS, and JavaScript

This document is intended to hand over the current repo state to another AI agent.

## Original Product Direction

The app is an educational web platform for kids with three roles:

- Admin
- Parents
- Kids

Main requirements currently reflected in the implementation:

- Admin manages curriculum and activities
- Parents assign activities to kids and see progress
- Kids can browse categories and play activities
- UI should be colorful, intuitive, and Kahoot-inspired
- Target stack is Vue frontend and Go backend with Gorm and Postgres

## Current Repository State

Top-level directories and important files:

- `backend/`
- `frontend/`
- `docs/plan/`
- `docs/updates/`
- `requirements.md`

Important note:

- The whole repo is still uncommitted from Git's perspective in this environment
- `frontend/node_modules/`, `frontend/dist/`, and `backend/dangde-world.db` currently exist in the working tree
- There is no `.gitignore` cleanup yet in this session

## Planning Documents Already Created

These planning notes already exist and should be read first by the next agent:

- `docs/plan/2026-04-15-initial-webapp-plan.md`
- `docs/plan/2026-04-15-ipad-safari-compatibility.md`
- `docs/plan/2026-04-15-custom-activity-runtime.md`

## Backend Status

### Architecture

The backend is structured in DDD-style layers:

- `backend/cmd/server/`
- `backend/internal/shared/config/`
- `backend/internal/domain/`
- `backend/internal/application/`
- `backend/internal/infrastructure/persistence/`
- `backend/internal/interfaces/http/`

Key backend entry points:

- `backend/cmd/server/main.go`
- `backend/internal/interfaces/http/router.go`
- `backend/internal/infrastructure/persistence/database.go`
- `backend/internal/infrastructure/persistence/seed.go`

### Current Domain Coverage

Implemented aggregates/entities:

- `user`
- `category`
- `activity`
- `assignment`
- `activitydata`

Key points:

- `activity` now includes `htmlCode`, `cssCode`, and `jsCode`
- `activitydata` stores runtime data for activities
- seed data includes demo users, categories, assignments, and starter activities

### API Coverage

Implemented API areas:

- Auth demo login
- Users
- Categories
- Activities
- Assignments
- Parent stats
- Activity runtime data CRUD

Important current endpoints:

- `POST /api/auth/login`
- `GET /api/users`
- `GET /api/categories`
- `GET /api/activities`
- `POST /api/activities`
- `PUT /api/activities/:id`
- `GET /api/assignments`
- `POST /api/assignments`
- `PATCH /api/assignments/:id`
- `GET /api/stats/parent/:parentId`
- `GET /api/activity-data`
- `POST /api/activity-data`
- `PUT /api/activity-data/:id`
- `DELETE /api/activity-data/:id`

### Configuration

Current backend env support:

- `backend/.env`
- `PORT`
- `FRONTEND_URL`
- `DATABASE_URL`

Important behavior:

- If `DATABASE_URL` is set, backend uses Postgres
- If `DATABASE_URL` is empty, backend falls back to SQLite in `backend/dangde-world.db`

This fallback was used to make local first-run easy, but the product direction still expects Postgres.

### Seed / Runtime Notes

The seed layer was updated to upsert starter categories, activities, and assignments when records already exist. This was done so activity code fields (`htmlCode`, `cssCode`, `jsCode`) are populated even on an existing local DB.

## Frontend Status

### Frontend Architecture

The frontend is built with Vue 3 and Vite.

Structure in use:

- `frontend/src/app/`
- `frontend/src/pages/`
- `frontend/src/features/`
- `frontend/src/entities/`
- `frontend/src/widgets/`
- `frontend/src/shared/`

Important frontend boot files:

- `frontend/src/main.js`
- `frontend/src/app/App.vue`
- `frontend/src/app/router.js`
- `frontend/src/app/styles/main.css`

### Current Pages

Implemented route-level pages:

- `frontend/src/pages/login/LoginPage.vue`
- `frontend/src/pages/kids/KidsPage.vue`
- `frontend/src/pages/kids/ActivityPlayerPage.vue`
- `frontend/src/pages/parents/ParentsPage.vue`
- `frontend/src/pages/admin/AdminPage.vue`

### Login / Session

Current auth is a demo login flow based on seeded users:

- User chooses a seeded demo user
- Session is stored in localStorage
- Role-based redirect is handled in router guards

Relevant files:

- `frontend/src/features/auth/session.js`
- `frontend/src/features/auth/LoginPanel.vue`
- `frontend/src/entities/user/api.js`

## Kids Experience Status

Kids dashboard currently supports:

- Category selection from sidebar/top nav
- Activity filtering by selected category
- Featured assignment panel
- Activity cards with progress
- Button to launch activity runtime player

Relevant files:

- `frontend/src/pages/kids/KidsPage.vue`
- `frontend/src/pages/kids/ActivityPlayerPage.vue`

Current player behavior:

- Activities render inside sandboxed `iframe`
- Activity HTML/CSS/JS are injected into `srcdoc`
- Parent window exposes a bridge as `window.DangdeAPI`
- Runtime JS can call:
  - `save`
  - `query`
  - `update`
  - `delete`
  - `setProgress`

Data flow:

- Runtime bridge proxies requests to `/api/activity-data`
- Progress updates go through `/api/assignments/:id`

## Parent Experience Status

Parent dashboard currently supports:

- Overview page with aggregate metrics
- Assignment section with create-assignment flow
- Reports section with per-kid summaries
- Sidebar-driven section switching

Relevant file:

- `frontend/src/pages/parents/ParentsPage.vue`

## Admin Experience Status

Admin Studio has gone through several iterations and is now partly reorganized.

Current behavior:

- Top-level sidebar sections:
  - Curriculum
  - Categories
  - Analytics
- Curriculum now has two sub-pages:
  - Management
  - Library
- Library is paginated
- Management uses a simple IDE-style authoring interface

Relevant file:

- `frontend/src/pages/admin/AdminPage.vue`

Current activity IDE:

- Metadata sidebar for title, description, prompt, language, difficulty, age group, category, and icon
- File-style navigation for:
  - `index.html`
  - `styles.css`
  - `app.js`
- Code editor area with line numbers and active tab stats

Relevant file:

- `frontend/src/features/activity-management/ActivityForm.vue`

### Important Admin Caveat

The current admin IDE is only a visual/simple IDE shell:

- No live preview yet
- No syntax highlighting
- No dirty-state tracking
- No linting/validation of activity JS/CSS/HTML

## iPad / Safari Work

Significant work has already been done for iPad Safari compatibility:

- safe-area handling
- `100svh` / `100dvh` support
- touch-friendly sizing
- tablet-aware navigation
- Safari viewport and keyboard/focus handling
- split-pane behavior on kids and player screens

Relevant files:

- `frontend/index.html`
- `frontend/src/app/providers/useSafariViewport.js`
- `frontend/src/shared/lib/useViewportProfile.js`
- `frontend/src/app/styles/main.css`

## Verification Status

Previously verified successfully during this session:

- backend build with `go build ./...`
- frontend build with `npm run build`

Important warning about the latest state:

- After the latest request to apply the same "fit to screen" rule to other pages, work was started but not completed/verified
- The following files were modified in that last step:
  - `frontend/src/shared/ui/AppShell.vue`
  - `frontend/src/pages/kids/KidsPage.vue`
  - `frontend/src/pages/parents/ParentsPage.vue`
- That last step did not complete the matching CSS updates and did not run a final build afterward

This means the most recent working tree should be treated as potentially inconsistent until the next agent verifies and finishes it.

## Latest In-Progress Change

The user asked:

- "please make sure other pages are using this rule as well"

The intended meaning was:

- apply the same "fit to screen" / viewport-fitted workspace rule used in Admin Management to other pages
- avoid long sliding page layouts
- prefer internal pane/content scrolling

What was started:

- `AppShell` now wraps slot content in `.app-content-inner`
- `KidsPage` was wrapped in `.page-workspace`
- `ParentsPage` was wrapped in `.page-workspace parents-workspace`

What is missing:

- corresponding CSS classes are not yet defined or verified
- no final build was run after those exact changes

This is the highest priority unfinished item.

## Known Functional Gaps

These are still open or only partially addressed:

- Activity bundle storage is still in database fields, not external storage/filesystem/object storage
- No Docker Compose setup yet
- No host-mounted storage strategy yet for activity bundles
- No category CRUD UI yet
- No admin analytics beyond simple summaries
- No production auth/authorization
- No backend tests
- No frontend component/unit tests
- No activity sandbox hardening beyond iframe isolation
- No live activity preview while editing
- No `.gitignore` cleanup for generated assets and local DB

## Storage Direction Already Discussed

The preferred future direction was explicitly discussed:

- Store activity code bundles in dedicated storage instead of database
- Keep runtime activity state in database

Preferred future model:

- DB stores metadata and versioning
- Storage holds actual bundle files
- Docker Compose should expose bundle storage to host

This refactor has not yet been implemented.

## Recommended Next Actions For The Next Agent

Highest priority:

1. Finish the cross-page "fit to screen" workspace refactor
2. Add the missing CSS for `.app-content-inner`, `.page-workspace`, and `.parents-workspace`
3. Run `npm run build` and fix any regressions from the latest incomplete change

Then:

1. Clean generated artifacts from versioned workspace expectations:
   - `frontend/node_modules/`
   - `frontend/dist/`
   - `backend/dangde-world.db`
2. Add `.gitignore`
3. Decide whether to keep SQLite fallback or move to Postgres-only local development

Then, depending on product priority:

1. Implement storage-backed activity bundles
2. Add Docker Compose with host-exposed bundle storage
3. Add live preview to Admin IDE
4. Add Category CRUD in Admin Studio

## Useful Files To Read First

For the next agent, these files give the clearest picture fastest:

- `requirements.md`
- `docs/plan/2026-04-15-initial-webapp-plan.md`
- `docs/plan/2026-04-15-ipad-safari-compatibility.md`
- `docs/plan/2026-04-15-custom-activity-runtime.md`
- `backend/internal/interfaces/http/router.go`
- `backend/internal/infrastructure/persistence/seed.go`
- `frontend/src/app/router.js`
- `frontend/src/app/styles/main.css`
- `frontend/src/pages/admin/AdminPage.vue`
- `frontend/src/features/activity-management/ActivityForm.vue`
- `frontend/src/pages/kids/ActivityPlayerPage.vue`

## Run / Verify Commands

Backend:

```bash
cd backend
go run ./cmd/server
```

Frontend:

```bash
cd frontend
npm run dev
```

Build verification:

```bash
cd backend
GOCACHE=$PWD/.gocache go build ./...

cd ../frontend
npm run build
```

## Final Handover Note

The codebase is beyond scaffold stage and already has meaningful product behavior for all three roles. The backend DDD structure, runtime activity bridge, seeded activity player, and iPad-focused UI work are all in place.

The biggest immediate risk is not architectural uncertainty. It is that the very latest UI refactor for "apply fit-to-screen rule to other pages" was started but not finished and not re-verified. The next agent should stabilize that first before adding new features.
