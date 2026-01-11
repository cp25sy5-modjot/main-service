BEGIN;

CREATE TABLE IF NOT EXISTS schema_migrations (
    version     VARCHAR(50) PRIMARY KEY,
    applied_at  TIMESTAMP NOT NULL DEFAULT now()
);

COMMIT;
