-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id          CHAR(36)     NOT NULL,
    parent_id   CHAR(36)     NULL,
    name        VARCHAR(150) NOT NULL,
    slug        VARCHAR(200) NOT NULL,
    description TEXT         NULL,
    is_active   TINYINT(1)   NOT NULL DEFAULT 1,
    created_at  DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at  DATETIME(6)  NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_categories_slug (slug),
    KEY idx_categories_parent_id (parent_id),
    KEY idx_categories_deleted_at (deleted_at),
    CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES categories (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id           CHAR(36)        NOT NULL,
    category_id  CHAR(36)        NOT NULL,
    name         VARCHAR(255)    NOT NULL,
    slug         VARCHAR(300)    NOT NULL,
    description  TEXT            NULL,
    price        DECIMAL(15,2)   NOT NULL DEFAULT 0.00,
    weight       DECIMAL(10,3)   NOT NULL DEFAULT 0.000,
    stock        INT             NOT NULL DEFAULT 0,
    is_active    TINYINT(1)      NOT NULL DEFAULT 1,
    created_at   DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at   DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at   DATETIME(6)     NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_products_slug (slug),
    KEY idx_products_category_id (category_id),
    KEY idx_products_is_active (is_active),
    KEY idx_products_price (price),
    KEY idx_products_deleted_at (deleted_at),
    FULLTEXT KEY ft_products_search (name, description),
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_images (
    id          CHAR(36)     NOT NULL,
    product_id  CHAR(36)     NOT NULL,
    url         TEXT         NOT NULL,
    alt_text    VARCHAR(255) NULL,
    is_primary  TINYINT(1)   NOT NULL DEFAULT 0,
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    KEY idx_product_images_product_id (product_id),
    KEY idx_product_images_is_primary (is_primary),
    CONSTRAINT fk_product_images_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_images;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
