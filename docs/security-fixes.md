# Security fixes — 2026-08-06

This documents the first round of fixes to `manga-tracker-api-go` after a
security pass. Four issues were fixed; each is a real, exploitable gap, not
just style.

## 1. Secrets baked into the Docker image

**Problem:** `.github/workflows/google.yml` had a "Create .env" step that
wrote `JWT_SIGNED_STRING`, `DB_USER`, `DB_PASS`, `ENCRYPTIONKEY`, and
`LINECLIENTID` into a `.env` file, and the Dockerfile then did `COPY . .` in
a single-stage `golang:latest` build. That `.env` file ended up as a layer
in the pushed image. Anyone able to pull the image from Artifact Registry
could read the file straight out of the image and get every secret,
including the JWT signing key — which is enough to forge an admin token and
call any endpoint, including `/admin/*`.

**Fix:**
- `Dockerfile` is now a multi-stage build: `golang:1.22-alpine` compiles a
  static binary, and only that binary is copied into a minimal `alpine`
  runtime image, running as a non-root user. No source, no `.env`, nothing
  but the compiled binary ships.
- `.dockerignore` now excludes `.env`, `.env.local`, `.git`, `.github`, and
  `.vscode` so a local `docker build .` can't accidentally pick them up
  either.
- `.github/workflows/google.yml` no longer writes a `.env` file at all.
  Secrets are passed straight to the Cloud Run revision via the `env_vars`
  input on `deploy-cloudrun@v2`, so they live in the Cloud Run
  configuration (visible only to project IAM), never inside the image.

**Follow-up worth doing (not done here):** move these from plain Cloud Run
env vars into Google Secret Manager and reference them via the `secrets:`
input of `deploy-cloudrun@v2` instead of `env_vars:` — that gets you
rotation and audit logging on secret access. Also worth adding a container
scan step (`docker scout` or Trivy) to CI; base-image CVEs will keep
resurfacing regardless of which tag is pinned, so this needs to be a
recurring check, not a one-time fix.

## 2. `config.GetEnv` required a `.env` file to exist

**Problem:** `config/env.go` called `godotenv.Load()` on *every* call to
`GetEnv()` and returned an error if no `.env` file was present. This only
worked because CI was writing a physical `.env` file into the image (see
#1). Once that stopped, every single `GetEnv()` call — including the JWT
secret lookup on every authenticated request — would have failed.

**Fix:** `.env` is now loaded once at process start (via `sync.Once`), and
a missing file is treated as normal (expected in any deployment where the
platform injects real environment variables) rather than an error.
`GetEnv()` always falls through to `os.Getenv`.

## 3. Hardcoded admin credentials

**Problem:** `handlers/auth_handlers/auth_handler.go` checked
`admin.Username != "admin1" || admin.Password != "admin"` directly in
source — a plaintext, hardcoded, weak password that's visible to anyone
with repo access.

**Fix:** `AdminLogin` verifies the password with
`bcrypt.CompareHashAndPassword`. `golang.org/x/crypto/bcrypt` was already a
transitive dependency (pulled in by the Mongo driver), so it's now promoted
to a direct one — no new dependency added.

**Update:** admin credentials originally lived in `ADMIN_USERNAME`/
`ADMIN_PASSWORD_HASH` env vars (single admin, no row anywhere). That's
since been replaced by a real `admins` table — `AdminLogin` looks the
username up via `repo.GetAdminByUsername` and compares against the stored
hash, same as before, just DB-backed instead of env-backed. Supports more
than one admin, and the JWT now carries a real `adminId` claim instead of
a hardcoded `"admin"` string. See the README's "Creating an admin account"
section for `scripts/seedadmin`, which replaces the old hashgen-then-set-
two-secrets flow — no GitHub Actions secrets needed for this anymore.

## 4. JWT verification didn't pin the signing algorithm

**Problem:** `JWTMiddleware`'s key function returned the HMAC secret
unconditionally, without checking what algorithm the *token itself* claimed
to use (`token.Method`). This is the standard "alg confusion" gap in JWT
libraries — the safe pattern is to only hand back a key for the algorithm
family you actually sign with.

**Fix:** the key function now checks `token.Method` is
`*jwt.SigningMethodHMAC` before returning the key, rejecting anything else
up front.

**Follow-up worth doing (not done here):** `golang-jwt/jwt` is on v4
(v5 is current, v4 is maintenance-only). Migrating to v5 and using
`jwt.WithValidMethods([]string{"HS256"})` would express the same
restriction more idiomatically, but that's a slightly bigger diff (API
changes across all JWT call sites) so it wasn't bundled into this pass.

## Update: resolved by the Postgres migration

The three items originally listed here as "not fixed yet" are gone as of
the MongoDB → PostgreSQL migration (see
[`postgres-migration.md`](postgres-migration.md)), though none of them were
the point of that migration — they fell out of the rewrite:

- `log.Fatalf` on ordinary DB errors — the old `repo/user.go` doesn't exist
  anymore; every function in the rewritten `repo/*.go` returns a plain
  `error` instead of killing the process.
- The unescaped Mongo connection string — moot, Postgres now connects via
  a single `DATABASE_URL` you construct yourself.
- The homegrown `EncryptHexId`/`shieftHex`/`meanHex` scheme — deleted.
  `service/user_service.go` is gone; users get a real generated `id` plus
  a unique `line_user_id` column instead of a derived-from-scratch ObjectID.

## Not fixed yet (flagged, deliberately left alone)

- `golang-jwt/jwt` is still on v4 (v5 is current, v4 is maintenance-only)
  — noted in section 4 above.
- `AdminLogin`'s error paths (`errors.New("no env for ...")`) return Go
  errors directly from a Fiber handler, which Fiber renders as a generic
  500 — fine functionally, just not a deliberately-designed error response.
