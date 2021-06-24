INSERT INTO users (id, name, email, created_at) VALUES
(gen_random_uuid(), 'vasia', 'vasia@mail.com', now() ),
(gen_random_uuid(), 'serhii', 'serhii@mail.com', now());
SELECT * FROM users WHERE name = 'vasia';
SELECT name FROM users WHERE id =  'd59c98ba-84af-44ed-919d-cbe875b95af3';