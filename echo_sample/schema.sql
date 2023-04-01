-- psql "postgres://username:password@localhost:5432/database_name" < schema.sql
-- Create authors table
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL
);

-- Create books table
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    isbn VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    author_id INTEGER REFERENCES authors (id) ON DELETE CASCADE
);

