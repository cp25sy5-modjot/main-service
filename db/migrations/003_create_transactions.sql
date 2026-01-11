BEGIN;

CREATE TABLE transactions (
    transaction_id VARCHAR PRIMARY KEY,
    user_id        VARCHAR NOT NULL,
    date           TIMESTAMP NOT NULL,
    type           VARCHAR NOT NULL,

    CONSTRAINT fk_transactions_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);

INSERT INTO schema_migrations(version)
VALUES ('003_create_transactions');

COMMIT;
