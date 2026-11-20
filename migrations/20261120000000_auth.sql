-- +goose Up
CREATE TABLE users (
 id SERIAL PRIMARY KEY,
 username TEXT NOT NULL UNIQUE,
 password_hash TEXT NOT NULL,
 role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin'))
);
ALTER TABLE classes ADD COLUMN owner_id INTEGER REFERENCES users(id);
CREATE INDEX idx_classes_owner_id ON classes(owner_id);
-- +goose Down
DROP INDEX idx_classes_owner_id;
ALTER TABLE classes DROP COLUMN owner_id;
DROP TABLE users;
