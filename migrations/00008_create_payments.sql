-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments (
    id               CHAR(36)        NOT NULL,
    order_id         CHAR(36)        NOT NULL,
    provider         VARCHAR(50)     NOT NULL     COMMENT 'e.g. doku',
    provider_ref     VARCHAR(255)    NULL         COMMENT 'transaction/reference ID from provider',
    method           VARCHAR(50)     NULL         COMMENT 'e.g. virtual_account, credit_card, qris',
    status           ENUM('pending','success','failed','expired','refunded') NOT NULL DEFAULT 'pending',
    amount           DECIMAL(15,2)   NOT NULL,
    currency         CHAR(3)         NOT NULL DEFAULT 'IDR',
    payment_url      TEXT            NULL         COMMENT 'payment URL to redirect customer',
    payload          JSON            NULL         COMMENT 'raw request payload sent to provider',
    callback_payload JSON            NULL         COMMENT 'raw callback payload from provider',
    paid_at          DATETIME(6)     NULL,
    expired_at       DATETIME(6)     NULL,
    created_at       DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at       DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uq_payments_order_id (order_id),
    KEY idx_payments_provider_ref (provider_ref),
    KEY idx_payments_status (status),
    CONSTRAINT fk_payments_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd
