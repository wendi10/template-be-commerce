-- +goose Up
-- Seed default admin user
-- Password: Admin@123  (bcrypt hash below)
-- +goose StatementBegin
INSERT IGNORE INTO users (id, email, password_hash, role, is_active, created_at, updated_at)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'admin@ecommerce.com',
    '$2a$12$7AK3UeAh6gj.H1HEsNxmFOIAMl8uqL/GkJcTN1x7vFcTsn9I7qJ3e',
    'admin',
    1,
    NOW(),
    NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'admin@ecommerce.com';
-- +goose StatementEnd
