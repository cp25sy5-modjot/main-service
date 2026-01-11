BEGIN;

-- users
CREATE INDEX idx_users_google_id ON users(google_id);

-- categories
CREATE INDEX idx_categories_user_id ON categories(user_id);

-- transactions
CREATE INDEX idx_transactions_user_date
    ON transactions(user_id, date DESC);

-- transaction_items
CREATE INDEX idx_items_transaction_id
    ON transaction_items(transaction_id);

CREATE INDEX idx_items_category_id
    ON transaction_items(category_id);

INSERT INTO schema_migrations(version)
VALUES ('005_create_indexes');

COMMIT;
