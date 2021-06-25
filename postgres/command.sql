INSERT INTO users (id, name, email, created_at) VALUES
(gen_random_uuid(), 'vasia', 'vasia@mail.com', now() ),
(gen_random_uuid(), 'serhii', 'serhii@mail.com', now());
SELECT * FROM users WHERE name = 'vasia';
SELECT name FROM users WHERE id =  'd59c98ba-84af-44ed-919d-cbe875b95af3';

/*returns names of users with titles of their posts*/
SELECT u.name,  p.title
FROM  users u
    LEFT JOIN posts p
        ON u.id = p.created_by
ORDER BY u.name;

/*returns posts with its comments*/
SELECT p.id, c.id
FROM posts p
    LEFT JOIN comments c
        ON p.id = c.post_id
ORDER BY p.id;

/*returns 5 names of users ordered by email */
SELECT name
FROM  users
ORDER BY email
LIMIT 5;

/*returns names of users which have at least 1 comment*/
SELECT u.name
FROM  users u
    LEFT JOIN posts p
        ON u.id = p.created_by
GROUP BY u.name
HAVING COUNT(p.id)>0;


SELECT to_char(created_at, 'YYYY-MM'), count(id)
FROM posts
GROUP BY to_char(created_at, 'YYYY-MM')
ORDER BY count(id) DESC
LIMIT 5;