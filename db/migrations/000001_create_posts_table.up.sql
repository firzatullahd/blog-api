BEGIN;

CREATE TYPE posts_status AS ENUM (
    'draft',
    'publish'
    );

COMMIT;

CREATE TABLE IF NOT EXISTS posts(
    id serial PRIMARY KEY,
    title text NOT NULL,
--     content text NOT NULL,
    tags int[] NOT NULL,
    status posts_status NOT NULL DEFAULT 'draft',
    publish_date timestamptz DEFAULT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz DEFAULT NULL
    );

create unique index if not exists idx_posts_id on posts(id);