BEGIN;

CREATE TABLE transaction_items (
    transaction_id VARCHAR NOT NULL,
    item_id        VARCHAR NOT NULL,
    title          VARCHAR NOT NULL,
    price          NUMERIC(12,2) NOT NULL,
    category_id    VARCHAR,

    PRIMARY KEY (transaction_id, item_id),

    CONSTRAINT fk_items_transaction
        FOREIGN KEY (transaction_id)
        REFERENCES transactions(transaction_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_items_category
        FOREIGN KEY (category_id)
        REFERENCES categories(category_id)
        ON DELETE SET NULL
);

INSERT INTO schema_migrations(version)
VALUES ('004_create_transaction_items');

COMMIT;
