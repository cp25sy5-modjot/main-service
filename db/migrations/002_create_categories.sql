BEGIN;

CREATE TABLE categories (
    category_id   VARCHAR PRIMARY KEY,
    user_id       VARCHAR NOT NULL,
    category_name VARCHAR NOT NULL,
    budget        NUMERIC(12,2),
    color_code    VARCHAR,
    created_at    TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT fk_categories_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);

INSERT INTO schema_migrations(version)
VALUES ('002_create_categories');

COMMIT;
