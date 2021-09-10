CREATE TABLE IF NOT EXISTS comments (
                                        id uuid PRIMARY KEY,
                                        content text NOT NULL,
                                        created_by uuid REFERENCES users (id) NOT NULL,
                                        created_at timestamp NOT NULL,
                                        post_id uuid REFERENCES posts (id) NOT NULL
);