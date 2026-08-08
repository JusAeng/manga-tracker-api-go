-- Relational redesign — see docs/postgres-schema.md for the rationale.
-- Domain: a manga work has zero or more official Thai editions (each from
-- a publisher), each edition has volumes. Authors and genres are
-- many-to-many against the manga work. Users follow manga (not editions
-- or volumes). No derived/denormalized values are stored — counts and
-- "latest volume" are computed on read.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    line_user_id  text UNIQUE NOT NULL,
    display_name  text,
    picture_url   text,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE manga (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title_original text,
    title_en       text,
    introduction   text,
    image_url      text,
    first_date_jp  date,
    status         text,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_manga_title_original ON manga (title_original);
CREATE INDEX idx_manga_title_en ON manga (title_en);

CREATE TABLE publishers (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text UNIQUE NOT NULL,
    website_url text,
    logo_url    text,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE thai_editions (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    manga_id      uuid NOT NULL REFERENCES manga (id) ON DELETE CASCADE,
    publisher_id  uuid NOT NULL REFERENCES publishers (id) ON DELETE RESTRICT,
    title_th      text NOT NULL,
    first_date_th date,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (manga_id, publisher_id)
);
CREATE INDEX idx_thai_editions_manga_id ON thai_editions (manga_id);
CREATE INDEX idx_thai_editions_publisher_id ON thai_editions (publisher_id);

CREATE TABLE volumes (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    thai_edition_id uuid NOT NULL REFERENCES thai_editions (id) ON DELETE CASCADE,
    volume_number   int NOT NULL,
    isbn            text UNIQUE,
    publish_date    date,
    price           numeric(10, 2),
    image_url       text,
    status          text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    UNIQUE (thai_edition_id, volume_number)
);
CREATE INDEX idx_volumes_thai_edition_id ON volumes (thai_edition_id);

CREATE TABLE follows (
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    manga_id    uuid NOT NULL REFERENCES manga (id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, manga_id)
);
CREATE INDEX idx_follows_manga_id ON follows (manga_id);

CREATE TABLE authors (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE manga_authors (
    manga_id   uuid NOT NULL REFERENCES manga (id) ON DELETE CASCADE,
    author_id  uuid NOT NULL REFERENCES authors (id) ON DELETE CASCADE,
    role       text NOT NULL CHECK (role IN ('author', 'artist', 'story', 'illustrator')),
    PRIMARY KEY (manga_id, author_id, role)
);
CREATE INDEX idx_manga_authors_author_id ON manga_authors (author_id);

CREATE TABLE genres (
    id    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name  text UNIQUE NOT NULL
);

CREATE TABLE manga_genres (
    manga_id  uuid NOT NULL REFERENCES manga (id) ON DELETE CASCADE,
    genre_id  uuid NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
    PRIMARY KEY (manga_id, genre_id)
);
CREATE INDEX idx_manga_genres_genre_id ON manga_genres (genre_id);

-- Supabase exposes every public-schema table over PostgREST by default,
-- regardless of whether this app uses that API. Enable RLS with no
-- policies so the PostgREST anon/authenticated roles are default-denied;
-- the Go backend connects as the postgres role (BYPASSRLS), so this has
-- no effect on it — access control stays entirely in the Go handlers.
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE manga ENABLE ROW LEVEL SECURITY;
ALTER TABLE publishers ENABLE ROW LEVEL SECURITY;
ALTER TABLE thai_editions ENABLE ROW LEVEL SECURITY;
ALTER TABLE volumes ENABLE ROW LEVEL SECURITY;
ALTER TABLE follows ENABLE ROW LEVEL SECURITY;
ALTER TABLE authors ENABLE ROW LEVEL SECURITY;
ALTER TABLE manga_authors ENABLE ROW LEVEL SECURITY;
ALTER TABLE genres ENABLE ROW LEVEL SECURITY;
ALTER TABLE manga_genres ENABLE ROW LEVEL SECURITY;
