CREATE TABLE IF NOT EXISTS transaction(
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    customer_id VARCHAR(36) NOT NULL,
    menu VARCHAR(100) NOT NULL,
    price BIGINT NOT NULL,
    qty BIGINT NOT NULL,
    payment VARCHAR(100) NOT NULL,
    total BIGINT NOT NULL,
    created_at BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
    FOREIGN KEY (customer_id) REFERENCES customer(id),
    INDEX menu_idx (menu)
);