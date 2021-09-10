CREATE TABLE IF NOT EXISTS posts (
                                     id uuid PRIMARY KEY,
                                     title text NOT NULL,
                                     description text NOT NULL,
                                     created_by uuid REFERENCES users(id) NOT NULL,
                                     created_at timestamp NOT NULL,
                                     tags json
);