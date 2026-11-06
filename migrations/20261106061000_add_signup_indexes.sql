-- +goose Up
CREATE INDEX idx_signups_class_id ON signups (class_id);

CREATE UNIQUE INDEX uq_signups_class_client
    ON signups (class_id, LOWER(client));

-- +goose Down
DROP INDEX IF EXISTS uq_signups_class_client;
DROP INDEX IF EXISTS idx_signups_class_id;
