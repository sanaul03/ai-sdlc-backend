CREATE TABLE IF NOT EXISTS vehicle_types (
    id          BIGSERIAL PRIMARY KEY,
    category_id BIGINT              NOT NULL REFERENCES vehicle_categories (id),
    name        VARCHAR(100)        NOT NULL,
    code        VARCHAR(50)         NOT NULL,
    capacity    INT,
    description TEXT,
    created_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    UNIQUE (category_id, code)
);

CREATE INDEX idx_vehicle_types_category_id ON vehicle_types (category_id);
CREATE INDEX idx_vehicle_types_code ON vehicle_types (code);
