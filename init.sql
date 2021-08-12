CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    name varchar NOT NULL,
    email varchar UNIQUE NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    password varchar(60) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    description text NOT NULL,
    created_by uuid REFERENCES users(id) NOT NULL,
    created_at timestamp NOT NULL,
    tags json
);

CREATE TABLE IF NOT EXISTS comments (
    id uuid PRIMARY KEY,
    content text NOT NULL,
    created_by uuid REFERENCES users (id) NOT NULL,
    created_at timestamp NOT NULL,
    post_id uuid REFERENCES posts (id) NOT NULL
);