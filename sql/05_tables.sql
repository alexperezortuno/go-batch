\c app_db
-- This file contains the SQL commands to create the schemas for the application.
CREATE TABLE IF NOT EXISTS batch.users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(50)  NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    email      VARCHAR(100) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    age        INT CHECK (age >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE batch.users ADD CONSTRAINT uni_users_email UNIQUE (email);

--- Comment columns
COMMENT ON COLUMN batch.users.id IS 'Primary key for the users table';
COMMENT ON COLUMN batch.users.username IS 'Unique username for the user';
COMMENT ON COLUMN batch.users.password IS 'Password for the user';
COMMENT ON COLUMN batch.users.email IS 'Unique email address for the user';
COMMENT ON COLUMN batch.users.name IS 'Name of the user';
COMMENT ON COLUMN batch.users.age IS 'Age of the user';
COMMENT ON COLUMN batch.users.created_at IS 'Timestamp when the user was created';
COMMENT ON COLUMN batch.users.updated_at IS 'Timestamp when the user was last updated';
