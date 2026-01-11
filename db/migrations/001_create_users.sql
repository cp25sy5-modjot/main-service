BEGIN;

CREATE TABLE users (
    user_id      VARCHAR PRIMARY KEY,
    google_id    VARCHAR,
    name         VARCHAR NOT NULL,
    status       VARCHAR NOT NULL,
    onboarding   BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO schema_migrations(version)
VALUES ('001_create_users');

COMMIT;
