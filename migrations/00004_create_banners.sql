-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banners (
    id          CHAR(36)     NOT NULL,
    title       VARCHAR(255) NOT NULL,
    subtitle    VARCHAR(255) NULL,
    image_url   TEXT         NOT NULL,
    link_url    TEXT         NULL,
    is_active   TINYINT(1)   NOT NULL DEFAULT 1,
    sort_order  INT          NOT NULL DEFAULT 0,
    start_date  DATETIME(6)  NULL,
    end_date    DATETIME(6)  NULL,
    created_at  DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at  DATETIME(6)  NULL,
    PRIMARY KEY (id),
    KEY idx_banners_is_active (is_active),
    KEY idx_banners_sort_order (sort_order),
    KEY idx_banners_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;
-- +goose StatementEnd
