# DangDe! World Initial Implementation Plan

## Goal

Build the first working version of DangDe! World as a web application for kids, parents, and admins based on the requirements in `requirements.md`, using DDD on the backend and a modern feature-oriented architecture on the frontend.

## Working Assumptions

- This repository starts as a greenfield project.
- The first delivery should prioritize a working vertical slice over full production hardening.
- Authentication can start as a simple role-based demo login instead of full password recovery, OAuth, or email verification.
- Curriculum and activities will be seeded with starter data matching the requirements:
  - Alphabets
  - Math
  - Basic Arabic Reading
  - Basic Clock Reading
- Parents will assign activities to kids in the first version.
- Admin users will manage categories, sub-categories, curriculum items, and activities through CRUD APIs and a basic UI.
- The first UI pass should be colorful, simple, and intuitive with a Kahoot-inspired visual direction.

## Delivery Scope

### Backend

- Use Go, Gin, Gorm, and Postgres.
- Use DDD-style layering:
  - `domain`: entities, repository contracts, domain rules
  - `application`: use cases and DTOs
  - `infrastructure`: Gorm persistence and database setup
  - `interfaces/http`: Gin handlers and route wiring
- Provide REST APIs for:
  - Demo login
  - Users
  - Categories
  - Activities
  - Parent activity assignments
  - Parent statistics
- Seed starter roles, users, categories, and educational activities.
- Keep the code organized so auth and persistence can be expanded later.

### Frontend

- Use Vue for the client application.
- Use a Feature-Sliced style frontend structure because it fits a growing product better than a page-only layout:
  - `app`: app bootstrap, router, providers
  - `pages`: route-level composition
  - `widgets`: dashboard sections and larger UI blocks
  - `features`: login, activity assignment, progress update, curriculum management
  - `entities`: user, activity, assignment, category models and API clients
  - `shared`: reusable UI, utilities, styling, constants
- Start with a login page.
- After login, route users to dashboards based on role:
  - Kids dashboard with colorful animated sidebar and activity cards
  - Parents dashboard with child statistics and assignment management
  - Admin dashboard for curriculum and activity management
- Use a small shared API layer and predictable component structure.

## Proposed Repository Structure

- `backend/`
  - `cmd/server`
  - `internal/shared/config`
  - `internal/domain`
  - `internal/application`
  - `internal/infrastructure`
  - `internal/interfaces/http`
- `frontend/`
  - `src/app`
  - `src/pages`
  - `src/widgets`
  - `src/features`
  - `src/entities`
  - `src/shared`
- `docs/`

## Implementation Phases

### Phase 1

- Scaffold backend service and data model
- Add seed data and basic role-aware endpoints
- Scaffold Vue application and app shell
- Implement login and role-based navigation

### Phase 2

- Build kid learning screens for starter activities
- Build parent assignment and statistics flows
- Build admin CRUD screens for categories and activities

### Phase 3

- Improve validation, auth, tests, and deployment setup
- Add progress tracking and richer activity content

## Data Model Outline

- `users`
  - id
  - name
  - role: `admin`, `parent`, `kid`
  - avatar
- `categories`
  - id
  - name
  - slug
  - type: `learning`, `playing`
  - parent_id
- `activities`
  - id
  - title
  - description
  - language
  - difficulty
  - category_id
  - age_group
- `assignments`
  - id
  - parent_id
  - kid_id
  - activity_id
  - status
  - progress

## Immediate Build Plan

1. Create backend with DDD boundaries, Gorm repositories, and seeded sample data.
2. Expose REST endpoints through Gin handlers and application services.
3. Create Vue app with Feature-Sliced structure and a colorful dashboard experience.
4. Wire the frontend to the backend with demo login and seeded data.
5. Verify both applications build cleanly.

## Risks

- Dependency installation may require network access and local toolchain availability.
- Postgres may not be available locally, so a fallback dev database config may be needed for first-run convenience.
- Requirements do not define detailed auth and multi-child parent flows, so those parts will be implemented with reasonable defaults.
