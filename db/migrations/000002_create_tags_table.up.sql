
CREATE TABLE IF NOT EXISTS tags(
    id serial PRIMARY KEY,
    label text UNIQUE NOT NULL,
--     posts int[] NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz DEFAULT NULL
);

create unique index if not exists idx_tags_label on tags(label);