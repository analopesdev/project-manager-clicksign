# Project Manager API

Go + chi + Postgres, wired with uber/fx.

## Running

```bash
docker compose up -d          # Postgres on localhost:5434
echo "DATABASE_URL=postgres://postgres:postgres@localhost:5434/project_manager?sslmode=disable" > .env
go run ./cmd                  # migrations run on startup
```


## Endpoints

| Method | Path | Notes |
|---|---|---|
| GET | `/projects?q=&favorite=&sort=` | `sort`: `name` (default), `end_date` or `start_date` |
| POST | `/projects` | 201 |
| GET | `/projects/{id}` | |
| PUT | `/projects/{id}` | full update |
| PATCH | `/projects/{id}/favorite` | body `{"favorite": true}` |
| DELETE | `/projects/{id}` | 204 |

Validation errors return **422** with per-field messages:
`{"error":"validation failed","fields":{"name":"name is required"}}`.
Unknown JSON fields return 400; bodies are limited to 1 MiB.

## Decisions

- **Total of projects:** the list response is `{"data": [...], "meta": {"total": N}}`.
  `total` is the number of registered projects, ignoring `q` and `favorite`, so the
  front does not need a second request.
- **"Closest to the end date" (`sort=end_date`):** upcoming projects first, soonest
  deadline first; projects already finished come last, most recently finished first.
- **"Most recently started" (`sort=start_date`):** `start_date` descending. Every sort
  ends with `id` as a tie-breaker, so the order is stable between requests.
- **Types:** `UUID`, `DATE`, `BOOLEAN`, `TIMESTAMPTZ`. Dates are exchanged as
  `YYYY-MM-DD`; timestamps come from the database (`RETURNING`), so POST and GET agree.
- **Constraints in the database:** `NOT NULL`, non-blank name/client and
  `CHECK (end_date >= start_date)`, besides validation in the domain.
- **Search:** `ILIKE` with a parameterised pattern; `%`, `_` and `\` typed by the user
  are escaped. `pg_trgm` GIN indexes on `name` and `client` support it.
- **Sort:** validated against an allow-list and mapped to fixed SQL, never concatenated
  from user input.
- **Logging:** request logs and panics are structured JSON through `slog`.
- **ID generation:** the database owns it (`DEFAULT gen_random_uuid()`); the `INSERT` does
  not send an id and reads it back with `RETURNING`. Malformed ids are rejected in the
  service as "not found" (404), so the repository does not validate input format.
- **Constraint violations:** if a `CHECK` fires (`projects_dates_check`, non-blank
  checks), the repository maps the constraint name to a field-level validation error
  (422) instead of a 500. The domain validates first; the constraints are the safety net.
- **Cover:** stored as `NULL` when absent, but the API always returns `cover_url` as a
  string, so "no cover" arrives as `""`.
- **`updated_at`:** means "content last edited". Favoriting does not change it.
- **Time zone:** the connection pins the session to `America/Sao_Paulo`, so `CURRENT_DATE`
  (which decides whether a project is finished) follows Brazilian dates, not the
  database server's zone.
