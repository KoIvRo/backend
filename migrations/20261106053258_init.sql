-- +goose Up
CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    trainer TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0)
);

CREATE TABLE signups (
    id SERIAL PRIMARY KEY,
    class_id INTEGER NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    client TEXT NOT NULL,
    phone TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS signups;
DROP TABLE IF EXISTS classes;
