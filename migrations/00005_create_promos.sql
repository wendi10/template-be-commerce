-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS promo_codes (
    id              CHAR(36)        NOT NULL,
    code            VARCHAR(50)     NOT NULL,
    name            VARCHAR(255)    NOT NULL DEFAULT '',
    description     TEXT            NULL,
    discount_type   ENUM('percentage','fixed') NOT NULL,
    discount_value  DECIMAL(15,2)   NOT NULL,
    min_purchase    DECIMAL(15,2)   NOT NULL DEFAULT 0.00,
    max_discount    DECIMAL(15,2)   NOT NULL DEFAULT 0.00,
    usage_limit     INT             NOT NULL DEFAULT 0 COMMENT '0 = unlimited',
    used_count      INT             NOT NULL DEFAULT 0,
    is_active       TINYINT(1)      NOT NULL DEFAULT 1,
    start_date      DATETIME(6)     NOT NULL,
    end_date        DATETIME(6)     NOT NULL,
    created_at      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at      DATETIME(6)     NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_promo_codes_code (code),
    KEY idx_promo_codes_is_active (is_active),
    KEY idx_promo_codes_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS promo_codes;
-- +goose StatementEnd
