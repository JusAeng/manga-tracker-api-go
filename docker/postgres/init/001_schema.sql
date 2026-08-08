-- Proposed schema from docs/postgres-schema.md.
-- This is not wired into the Go app yet — it's here so the schema can be
-- reviewed/tested against a real database before the repo/models rewrite.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    line_user_id  text UNIQUE NOT NULL,
    name          text NOT NULL,
    image         text,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE manga (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title              text NOT NULL,
    author             text,
    other_titles       text[] NOT NULL DEFAULT '{}',
    other_participate  text[] NOT NULL DEFAULT '{}',
    genre              text,
    other_genres       text[] NOT NULL DEFAULT '{}',
    image              text,
    introduction       text,
    publisher          text,
    first_date_jp      text,
    first_date_th      text,
    last_vol           int NOT NULL DEFAULT 0,
    is_highlight       boolean NOT NULL DEFAULT false,
    subscribers_count  int NOT NULL DEFAULT 0,
    score              numeric(3, 2) NOT NULL DEFAULT 0,
    total_voters       int NOT NULL DEFAULT 0,
    created_at         timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_manga_title ON manga (title);
CREATE INDEX idx_manga_is_highlight ON manga (is_highlight) WHERE is_highlight;

CREATE TABLE vols (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    manga_id           uuid NOT NULL REFERENCES manga (id) ON DELETE CASCADE,
    vol_number         int NOT NULL,
    image              text,
    publish_date       text,
    total_owner_count  int NOT NULL DEFAULT 0,
    UNIQUE (manga_id, vol_number)
);
CREATE INDEX idx_vols_manga_id ON vols (manga_id);

CREATE TABLE subscriptions (
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    manga_id    uuid NOT NULL REFERENCES manga (id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, manga_id)
);
CREATE INDEX idx_subscriptions_manga_id ON subscriptions (manga_id);

CREATE TABLE owned_volumes (
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    vol_id      uuid NOT NULL REFERENCES vols (id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, vol_id)
);
CREATE INDEX idx_owned_volumes_vol_id ON owned_volumes (vol_id);

CREATE TABLE ratings (
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    manga_id    uuid NOT NULL REFERENCES manga (id) ON DELETE CASCADE,
    score       smallint NOT NULL CHECK (score BETWEEN 1 AND 5),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, manga_id)
);
CREATE INDEX idx_ratings_manga_id ON ratings (manga_id);

-- Keep manga.score / manga.total_voters correct automatically, instead of
-- the hand-rolled running-average update the Mongo code did.
CREATE OR REPLACE FUNCTION refresh_manga_rating_stats() RETURNS trigger AS $$
DECLARE
    target_manga_id uuid := COALESCE(NEW.manga_id, OLD.manga_id);
BEGIN
    UPDATE manga
    SET score = COALESCE((SELECT AVG(score) FROM ratings WHERE manga_id = target_manga_id), 0),
        total_voters = (SELECT COUNT(*) FROM ratings WHERE manga_id = target_manga_id)
    WHERE id = target_manga_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_ratings_refresh_stats
AFTER INSERT OR UPDATE OR DELETE ON ratings
FOR EACH ROW EXECUTE FUNCTION refresh_manga_rating_stats();

-- Supabase exposes every public-schema table over PostgREST by default,
-- regardless of whether this app uses that API. Enable RLS with no
-- policies so the PostgREST anon/authenticated roles are default-denied;
-- the Go backend connects as the postgres role (BYPASSRLS), so this has
-- no effect on it — access control stays entirely in the Go handlers.
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE manga ENABLE ROW LEVEL SECURITY;
ALTER TABLE vols ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE owned_volumes ENABLE ROW LEVEL SECURITY;
ALTER TABLE ratings ENABLE ROW LEVEL SECURITY;
