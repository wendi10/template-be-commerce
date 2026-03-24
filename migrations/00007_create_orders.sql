-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id              CHAR(36)        NOT NULL,
    customer_id     CHAR(36)        NOT NULL,
    address_id      CHAR(36)        NOT NULL,
    promo_code_id   CHAR(36)        NULL,
    order_number    VARCHAR(50)     NOT NULL,
    status          ENUM('pending','waiting_payment','paid','processing','shipped','completed','cancelled')
                    NOT NULL DEFAULT 'pending',
    subtotal        DECIMAL(15,2)   NOT NULL DEFAULT 0.00,
    discount_amount DECIMAL(15,2)   NOT NULL DEFAULT 0.00,
    shipping_cost   DECIMAL(15,2)   NOT NULL DEFAULT 0.00,
    total_amount    DECIMAL(15,2)   NOT NULL DEFAULT 0.00,
    notes           TEXT            NULL,
    created_at      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uq_orders_order_number (order_number),
    KEY idx_orders_customer_id (customer_id),
    KEY idx_orders_status (status),
    KEY idx_orders_created_at (created_at),
    CONSTRAINT fk_orders_customer   FOREIGN KEY (customer_id)   REFERENCES customers   (id),
    CONSTRAINT fk_orders_address    FOREIGN KEY (address_id)    REFERENCES addresses   (id),
    CONSTRAINT fk_orders_promo_code FOREIGN KEY (promo_code_id) REFERENCES promo_codes (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items (
    id            CHAR(36)      NOT NULL,
    order_id      CHAR(36)      NOT NULL,
    product_id    CHAR(36)      NOT NULL,
    product_name  VARCHAR(255)  NOT NULL,
    product_slug  VARCHAR(300)  NOT NULL DEFAULT '',
    quantity      INT           NOT NULL DEFAULT 1,
    unit_price    DECIMAL(15,2) NOT NULL,
    subtotal      DECIMAL(15,2) NOT NULL,
    created_at    DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    KEY idx_order_items_order_id (order_id),
    KEY idx_order_items_product_id (product_id),
    CONSTRAINT fk_order_items_order   FOREIGN KEY (order_id)   REFERENCES orders   (id) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
