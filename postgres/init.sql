CREATE SCHEMA models
    CREATE TABLE models.user (ID varchar, name varchar, email varchar, createdAt date, updatedAt date)
    CREATE TABLE models.post (ID varchar, title text, description text, createdBy varchar, createdAt date, tags varchar array)
    CREATE TABLE models.comment (ID varchar, content text, createdBy varchar, createdAt date, postId varchar);