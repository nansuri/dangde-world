# DangDe! World Custom Activity Runtime Plan

## Objective

Allow admins to create activities by providing only:

- HTML
- CSS
- JavaScript

The kid-facing app should render those activities, and the activity JavaScript should be able to use a platform API for data persistence.

## Product Direction

- Activities become configurable mini-apps.
- The platform owns:
  - activity metadata
  - activity rendering shell
  - assignment context
  - persistence API
- Activity authors own:
  - markup
  - styling
  - interaction logic

## Core Design

### Activity definition

Each activity stores:

- metadata: title, description, category, icon, etc.
- `html_code`
- `css_code`
- `js_code`

### Activity player

- Kids open an activity in a dedicated player surface.
- The player renders the admin-authored activity inside a sandboxed `iframe`.
- The app injects a small runtime bridge into the iframe.
- Activity JavaScript uses the bridge instead of calling backend APIs directly.

### Activity data API

The platform exposes CRUD operations for activity-authored data:

- `save`
- `query`
- `update`
- `delete`

The data model should include:

- activity id
- assignment id
- kid id
- key
- value JSON

This allows each activity to save progress, answers, preferences, and temporary state.

## Backend Changes

- Extend `activity` aggregate to store HTML, CSS, and JS code.
- Add a new `activity data` aggregate and repository.
- Add application services for:
  - CRUD activity definitions
  - CRUD activity runtime data
- Expose HTTP endpoints:
  - `GET /api/activity-data`
  - `POST /api/activity-data`
  - `PUT /api/activity-data/:id`
  - `DELETE /api/activity-data/:id`

## Frontend Changes

- Admin activity form gains HTML, CSS, and JS editors.
- Kids dashboard gets a `Play Activity` action.
- Add an `ActivityPlayer` page with:
  - iframe runtime
  - postMessage bridge
  - SDK methods: `save`, `query`, `update`, `delete`, `setProgress`

## Security Note

This feature intentionally allows admin-authored JavaScript execution. For the first version, the trust model is:

- admin-authored activity code is trusted application content

The runtime will still be isolated inside a sandboxed iframe, but this is not equivalent to running arbitrary untrusted third-party code.

## First Delivery Scope

1. Store HTML, CSS, and JS with each activity.
2. Add runtime data CRUD API.
3. Build the sandboxed activity player.
4. Add one seeded sample activity using the new runtime model.
5. Keep existing role flows working.
