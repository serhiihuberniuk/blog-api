CREATE TABLE user (
    id uuid PRIMARY KEY,
    name varchar,
    email varchar,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE post (
    id uuid PRIMARY KEY,
    title text,
    description text,
    created_by varchar REFERENCES user(id),
    created_at timestamp,
    tags varchar array
);

CREATE TABLE comment (
    id uuid PRIMARY KEY,
    content text,
    created_by varchar REFERENCES user (id),
    created_at timestamp,
    post_id varchar REFERENCES post (id)
);