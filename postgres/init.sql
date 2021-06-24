CREATE TABLE users (
    id uuid PRIMARY KEY,
    name varchar,
    email varchar,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE posts (
    id uuid PRIMARY KEY,
    title text,
    description text,
    created_by uuid REFERENCES users(id),
    created_at timestamp,
    tags json
);

CREATE TABLE comments (
    id uuid PRIMARY KEY,
    content text,
    created_by uuid REFERENCES users (id),
    created_at timestamp,
    post_id uuid REFERENCES posts (id)
);