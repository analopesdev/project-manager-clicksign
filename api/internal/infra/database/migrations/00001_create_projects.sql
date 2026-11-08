-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE projects (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    client      TEXT        NOT NULL,
    start_date  DATE        NOT NULL,
    end_date    DATE        NOT NULL,
    cover_url   TEXT,
    favorite    BOOLEAN     NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT projects_name_not_blank   CHECK (btrim(name) <> ''),
    CONSTRAINT projects_client_not_blank CHECK (btrim(client) <> ''),
    CONSTRAINT projects_cover_not_blank  CHECK (btrim(cover_url) <> ''),
    CONSTRAINT projects_dates_check      CHECK (end_date >= start_date)
);

CREATE INDEX projects_name_trgm_idx   ON projects USING gin (name gin_trgm_ops);
CREATE INDEX projects_client_trgm_idx ON projects USING gin (client gin_trgm_ops);
CREATE INDEX projects_start_date_idx  ON projects (start_date);
CREATE INDEX projects_end_date_idx    ON projects (end_date);

-- +goose Down
DROP TABLE projects;