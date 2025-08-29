CREATE TYPE user_role AS ENUM ('admin', 'manager', 'client');
CREATE TYPE rental_status AS ENUM ('pending', 'active', 'completed', 'canceled');

CREATE TABLE users
(
    id            UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    name          VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    role          user_role           NOT NULL DEFAULT 'client'
);

-- Админ (пароль: password)
INSERT INTO users (id, name, email, password_hash, role)
VALUES ('11111111-1111-1111-1111-111111111111',
        'System Admin',
        'admin@example.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'admin'),
-- Менеджер (пароль: password)
       ('22222222-2222-2222-2222-222222222222',
        'Test Manager',
        'manager@example.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'manager'),
-- Клиент (пароль: password)
       ('33333333-3333-3333-3333-333333333333',
        'Test Client',
        'client@example.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'client');

CREATE TABLE cars
(
    id            UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    model         VARCHAR(255)   NOT NULL,
    brand         VARCHAR(255)   NOT NULL,
    year          INTEGER        NOT NULL,
    price_per_day NUMERIC(10, 2) NOT NULL,
    is_available  BOOLEAN        NOT NULL DEFAULT TRUE
);

INSERT INTO cars (id, model, brand, year, price_per_day, is_available)
VALUES ('44444444-4444-4444-4444-444444444444',
        'Model S',
        'Tesla',
        2023,
        100.00,
        true);

CREATE TABLE rentals
(
    id          UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    car_id      UUID           NOT NULL REFERENCES cars (id),
    user_id     UUID           NOT NULL REFERENCES users (id),
    start_date  TIMESTAMP      NOT NULL,
    end_date    TIMESTAMP      NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL,
    status      rental_status  NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMP      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cars_availability ON cars (is_available);
CREATE INDEX idx_rentals_dates ON rentals (start_date, end_date);
CREATE INDEX idx_rentals_created_at ON rentals (created_at);