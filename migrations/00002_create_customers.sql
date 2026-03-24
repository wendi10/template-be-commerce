-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS customers (
    id           CHAR(36)     NOT NULL,
    user_id      CHAR(36)     NOT NULL,
    first_name   VARCHAR(100) NOT NULL DEFAULT '',
    last_name    VARCHAR(100) NOT NULL DEFAULT '',
    phone        VARCHAR(20)  NULL,
    avatar       TEXT         NULL,
    created_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at   DATETIME(6)  NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_customers_user_id (user_id),
    KEY idx_customers_deleted_at (deleted_at),
    CONSTRAINT fk_customers_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS addresses (
    id             CHAR(36)     NOT NULL,
    customer_id    CHAR(36)     NOT NULL,
    label          VARCHAR(100) NOT NULL DEFAULT 'Home',
    recipient_name VARCHAR(150) NOT NULL,
    phone          VARCHAR(20)  NOT NULL,
    address_line1  VARCHAR(255) NOT NULL,
    address_line2  VARCHAR(255) NULL,
    city           VARCHAR(100) NOT NULL,
    province       VARCHAR(100) NOT NULL,
    postal_code    VARCHAR(10)  NOT NULL,
    is_default     TINYINT(1)   NOT NULL DEFAULT 0,
    created_at     DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at     DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at     DATETIME(6)  NULL,
    PRIMARY KEY (id),
    KEY idx_addresses_customer_id (customer_id),
    KEY idx_addresses_deleted_at (deleted_at),
    CONSTRAINT fk_addresses_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS addresses;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS customers;
-- +goose StatementEnd
