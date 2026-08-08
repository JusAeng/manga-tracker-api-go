# manga-tracker-api-go

Go/Fiber backend for the manga tracker. Handles LINE login, manga/volume
data, and per-user subscribe/owner/rating lists.

## Current state

The app runs on **PostgreSQL** — see
[`docs/postgres-schema.md`](docs/postgres-schema.md) for the schema and the
reasoning behind it, and [`docs/postgres-migration.md`](docs/postgres-migration.md)
for what actually changed in the cutover from the old MongoDB backend.

## Environment variables

Copy `.env.example` to `.env` for local development (never commit `.env` —
it's gitignored, and CI injects secrets separately, see
[`docs/security-fixes.md`](docs/security-fixes.md)).

| Variable              | Used for                                              |
|------------------------|--------------------------------------------------------|
| `DATABASE_URL`         | Postgres connection string. Recommended: a Supabase **Session pooler** URI, e.g. `postgresql://postgres.<project-ref>:<url-encoded-password>@aws-0-<region>.pooler.supabase.com:5432/postgres` (URL-encode special characters in the password — see the note in `.env.example`). Alternatively `postgres://manga_tracker:manga_tracker_dev@localhost:5432/manga_tracker?sslmode=disable` for the local Docker Compose setup below |
| `JWT_SIGNED_STRING`    | HMAC secret used to sign/verify session JWTs            |
| `LINECLIENTID`         | LINE login channel ID, used to verify LINE ID tokens    |
| `ADMIN_USERNAME`       | Admin login username                                    |
| `ADMIN_PASSWORD_HASH`  | bcrypt hash of the admin password — generate with the snippet in `docs/security-fixes.md`, never store the raw password |

## Running locally

```bash
go mod download
make run          # or: go run main.go
```

The server listens on `:8080`.

## Local database

```bash
docker compose up -d postgres
```

This starts Postgres on `localhost:5432` (db `manga_tracker`, user
`manga_tracker`, password `manga_tracker_dev` — local dev only, never used
in any deployed environment) and loads the schema from
`docker/postgres/init/001_schema.sql` automatically on first start. Point
`DATABASE_URL` in your `.env` at it (see the table above).

```bash
docker compose exec postgres psql -U manga_tracker -d manga_tracker
```

To reset it (drops all data and reloads the init script):

```bash
docker compose down -v
docker compose up -d postgres
```

## API endpoints

Every route except the two under **Auth** requires an
`Authorization: Bearer <token>` header. Get a token from `POST /auth`
(real LINE login), `POST /auth/admin` (admin login), or `scripts/devtoken`
for local testing without either (see the section below).

### Auth (public)

| Method | Path           | Body                                | Returns                    |
|--------|----------------|--------------------------------------|-----------------------------|
| POST   | `/auth`        | `{ "token": "<LINE ID token>" }`     | `{ token, profile }` — verifies the token against LINE's API, creates the user on first login |
| POST   | `/auth/admin`  | `{ "username", "password" }`         | `{ token }` — checked against `ADMIN_USERNAME`/`ADMIN_PASSWORD_HASH` |

### Manga (requires a valid token, any role)

| Method | Path                 | Notes |
|--------|----------------------|-------|
| GET    | `/manga`             | List all manga |
| GET    | `/manga/:id`         | Single manga; `404` if it doesn't exist |
| GET    | `/manga/highlight`   | The one manga flagged `is_highlight`; `404` if none is set |
| GET    | `/manga/trending`    | 3 random manga — placeholder logic, not actual trending data |
| GET    | `/manga/new`         | 5 random manga — placeholder, not sorted by release date |
| GET    | `/manga/recommend`   | 8 random manga — placeholder |

### User (requires a valid token, any role)

| Method | Path                        | Body / Params                          | Notes |
|--------|-----------------------------|------------------------------------------|-------|
| GET    | `/user/profile`             | —                                         | Current user's profile; `null` if no matching row |
| PATCH  | `/user/profile`             | `{ "key": "name"\|"image", "value": "..." }` | Only `name`/`image` are allowed keys |
| GET    | `/user/subscribelist`       | —                                         | Manga the current user is subscribed to |
| PUT    | `/user/subscribe/:id`       | `:id` = manga id                          | Toggles subscription |
| PUT    | `/user/ownerlist/:id/:vol`  | `:id` = manga id, `:vol` = volume number  | Toggles ownership of that volume |
| PUT    | `/user/rating/:id/:score`   | `:id` = manga id, `:score` = 0-5          | Rates the manga |

### Admin (requires a token with `role: admin`, from `/auth/admin`)

| Method | Path                | Body / Params                          | Notes |
|--------|---------------------|-------------------------------------------|-------|
| GET    | `/admin/users`       | —                                          | List all users |
| DELETE | `/admin/user/:id`    | —                                          | Delete a user |
| POST   | `/admin/manga`       | `Manga` JSON                               | Create a manga |
| PATCH  | `/admin/manga/`      | `Manga` JSON (must include `id`)           | Update a manga — note: id comes from the body, not the URL |
| DELETE | `/admin/manga/:id`   | —                                          | Delete a manga |
| POST   | `/admin/vol/:id`     | `Vol` JSON, `:id` = manga id               | Add a volume |
| PATCH  | `/admin/vol/:id`     | `Vol` JSON, `:id` = manga id               | Update a volume |
| DELETE | `/admin/vol/:id`     | `{ "volNumbers": [1, 2, 3] }`, `:id` = manga id | Delete one or more volumes |

`GET /hello` and `GET /foo/auth/:id` also exist but are leftover mock/dev
routes, not part of the real API surface.

## Testing endpoints without a real LINE login

The normal user flow (`POST /auth`) requires a real LINE ID token verified
against LINE's API — there's no way around that for the actual login route.
For local testing, `scripts/devtoken` mints a locally-signed JWT instead, so
you can hit `/user/*` and `/manga/*` routes directly. It only needs
`JWT_SIGNED_STRING` from your `.env` — it never calls LINE or touches the
database.

```bash
make devtoken                                                       # role=user, random userId
make devtoken ARGS="-role admin"                                     # role=admin
make devtoken ARGS="-user-id 3fa85f64-5717-4562-b3fc-2c963f66afa6"   # reuse a real user's id
```

or directly: `go run ./scripts/devtoken -role user -user-id <uuid> -hours 6`.

It prints the token plus a ready-to-run `curl` example. If you don't pass
`-user-id`, it generates a fresh one — routes that hit the database (e.g.
`GET /user/profile`) will return `null` unless a matching row in `users`
actually exists with that id. To test with real data, log in once through
the real LINE flow (or insert a row into `users` by hand), copy that
user's `id`, and pass it as `-user-id` from then on.

**This is a dev-only tool.** Whoever holds `JWT_SIGNED_STRING` can use it to
mint a token indistinguishable from a real login for any user or the admin
role — never run it against a production `.env`, and it should never ship
in a deployed image (it isn't: the Dockerfile only builds the root `main.go`
package, so `scripts/` is never compiled into the runtime image).

## Deployment

Pushes to `main` build and deploy to Cloud Run via
`.github/workflows/google.yml`. Secrets are injected as Cloud Run
environment variables at deploy time — see
[`docs/security-fixes.md`](docs/security-fixes.md) for why that changed
from the old "bake `.env` into the image" approach.
