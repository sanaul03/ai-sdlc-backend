CREATE TABLE IF NOT EXISTS vehicle_categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100)        NOT NULL,
    code        VARCHAR(50)         NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vehicle_categories_code ON vehicle_categories (code);
