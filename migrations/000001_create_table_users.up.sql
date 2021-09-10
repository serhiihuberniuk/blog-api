CREATE TABLE IF NOT EXISTS users (
                                     id uuid PRIMARY KEY,
                                     name varchar NOT NULL,
                                     email varchar UNIQUE NOT NULL,
                                     created_at timestamp NOT NULL,
                                     updated_at timestamp NOT NULL,
                                     password varchar(60) UNIQUE NOT NULL
);