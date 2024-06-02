BEGIN;

CREATE TYPE users_role AS ENUM (
    'user',
    'admin'
    );

COMMIT;


CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    email VARCHAR (100) UNIQUE NOT NULL,
    password text NOT NULL,
    name VARCHAR (50) UNIQUE NOT NULL,
    role text not null,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz DEFAULT NULL
);

