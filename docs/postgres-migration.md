# MongoDB → PostgreSQL migration

This documents the actual cutover from MongoDB to the schema designed in
[`postgres-schema.md`](postgres-schema.md). The schema didn't change from
that design; what's below is how the Go code maps onto it, and where the
migration turned into some incidental fixes beyond a straight port.

## What changed, file by file

- **`db/db.go`** — `mongo.Client` replaced with a `pgxpool.Pool`, connected
  via a single `DATABASE_URL` env var instead of `DB_USER`/`DB_PASS`
  string concatenation. `main.go` didn't need any changes — `db.Connect()`
  kept the same signature.
- **`models/user_model.go`, `models/manga_model.go`** — `primitive.ObjectID`
  replaced with `uuid.UUID` (`github.com/google/uuid`) everywhere. `User`
  dropped `SubscribeList`/`OwnerList`/`RateList` (now separate tables).
  `Manga` gained `IsHighlight`/`SubscribersCount`/`TotalVoters` fields
  matching the new columns.
- **`repo/user.go`, `repo/manga.go`, `repo/vol.go`** — fully rewritten
  against parameterized SQL. `SubscribeMangaById`, `UpdateOwnerList`, and
  `AddMangaVol`/`DeleteManyVols` now run inside a `pgx.Tx`, so the counter
  update (`subscribers_count`, `total_owner_count`, `last_vol`) happens
  atomically with the row insert/delete — the read-modify-write races
  flagged during the original review are gone structurally, not patched.
  `ratings`/`manga.score`/`manga.total_voters` are kept correct by the
  `trg_ratings_refresh_stats` trigger in the schema, replacing the
  incorrect running-average formula `UpdateMangaScore` used to compute.
- **`service/user_service.go`** — deleted. `EncryptHexId`/`shieftHex`/
  `meanHex` existed only to turn a LINE `sub` into something that looked
  like a Mongo ObjectID; `repo.GetOrCreateUserByLineID` now does that with
  a single `INSERT ... ON CONFLICT (line_user_id) DO UPDATE ... RETURNING`
  upsert, so there's no custom encoding and no first-login race between
  two concurrent logins from the same account.
- **`service/auth_service.go`** — deleted. `GetUserIdFromJWT` was an unused
  stub that always returned `nil`; it also imported the Mongo driver, so
  it had to go either way.
- **`handlers/*`** — every `primitive.ObjectIDFromHex` call became
  `uuid.Parse`. Two small behavior changes fell out of this naturally:
  - `GetMangaByIdHandler` now returns `404` for a manga that doesn't
    exist, instead of an empty array (`repo.GetMangaById` returns a single
    `*models.Manga`, not a Mongo-style slice-of-0-or-1).
  - `GetMangaHighlight` reads the `manga.is_highlight` flag via
    `repo.GetHighlightManga()` instead of hitting a hardcoded ObjectID
    (`662d5f00d657e10679478b83`) that only ever pointed at one specific
    Mongo document. `repo.SetMangaHighlight(mangaId)` is available to flip
    which manga is featured; nothing calls it yet — wire it up to an admin
    endpoint when there's a UI for it.
- **`scripts/devtoken`** — generates a `uuid.NewString()` instead of a
  Mongo ObjectID hex string for the placeholder `userId` claim.
- **`go.mod`** — `go.mongodb.org/mongo-driver` dropped entirely (confirmed
  via `go mod tidy` — nothing references it anymore).
  `github.com/jackc/pgx/v5` added. Its current release requires Go ≥1.25,
  which bumped the `go` directive from 1.22.1 → 1.25.0 as a side effect —
  not something this migration set out to do, but a welcome one given how
  stale the toolchain version was. **`Dockerfile`'s builder image was
  bumped from `golang:1.22-alpine` to `golang:1.25-alpine` to match** —
  without that the image build would fail outright.

## Incidental fixes that came from the rewrite, not from a deliberate pass

- The `log.Fatalf`-on-ordinary-DB-error pattern flagged during the
  original security review (`repo/user.go`'s old `DeleteUserById` and
  `UpdateUserProfile`, which would `os.Exit(1)` the whole process on one
  failed request) is gone — every new `repo/*.go` function returns a plain
  `error`.
- The unescaped Mongo connection string (`fmt.Sprintf` concatenating
  `DB_PASS` into a URI) is moot — Postgres now connects via a single
  `DATABASE_URL` you construct yourself, so there's no string-building in
  Go to get wrong.
- `DeleteAllVolsByMangaId` used to set Mongo's `lastVol` field to the
  string `"0"` against an `int`-typed field — a real type bug in the old
  code. The Postgres version sets an actual `int` column.
- Volume ownership (`owned_volumes`) is now a real foreign key into
  `vols(id)`, not just a manga-level existence check
  (`isMangaExist`) — you can no longer record ownership of a volume
  number that was never created, which the old code didn't actually
  prevent.

## Environment variables

Removed: `DB_USER`, `DB_PASS`, `ENCRYPTIONKEY` (no longer meaningful —
see `EncryptHexId` above).
Added: `DATABASE_URL`.
Unchanged: `JWT_SIGNED_STRING`, `LINECLIENTID`, `ADMIN_USERNAME`,
`ADMIN_PASSWORD_HASH`. See the README's environment variables table.

**Action required:** `.github/workflows/google.yml` now passes
`DATABASE_URL` to Cloud Run instead of `DB_USER`/`DB_PASS`/`ENCRYPTIONKEY`.
Add a `DATABASE_URL` secret to the GitHub repo (pointing at whatever
Postgres instance production should use — this repo doesn't provision one
for you), and delete the now-unused `ENCRYPTIONKEY`/`DB_USER`/`DB_PASS`
secrets whenever you're ready.

## Deliberately not done

- **No Mongo → Postgres data migration script.** There was no reachable
  production database at the time of this migration. If real user data
  in Mongo turns up later, write a one-off export/import script before
  pointing a live deployment at this schema — don't assume the tables
  start empty.
- **`golang-jwt` is still on v4.** Unrelated to this migration; still an
  open follow-up from `security-fixes.md`.
- **No wiring for `SetMangaHighlight`** into an admin route — it exists in
  `repo/manga.go` but nothing calls it yet.

## Verifying this migration

Everything above was verified with `go build ./...` and `go vet ./...`
after each file was rewritten — both pass clean. What could **not** be
verified in the environment this migration was written in: actually
running the compiled binary against a live Postgres instance (a local
`dyld: missing LC_UUID` sandbox issue blocks running *any* locally-built
Go binary there, unrelated to this change — see the note in the main
README's dev-token section). Before relying on this:

```bash
docker compose up -d postgres
# set DATABASE_URL in .env to point at it, plus the other required vars
make run
make devtoken   # mint a token, then hit /manga, /user/profile, etc.
```
