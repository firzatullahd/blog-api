CREATE TABLE IF NOT EXISTS r_post_tag(
        id serial PRIMARY KEY,
        post_id bigint not null references posts(id),
        tag_id bigint not null references tags(id),
        created_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz DEFAULT NULL
);

