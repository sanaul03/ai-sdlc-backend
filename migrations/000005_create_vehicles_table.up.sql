CREATE TABLE IF NOT EXISTS vehicles (
    id                  BIGSERIAL PRIMARY KEY,
    company_id          BIGINT              NOT NULL REFERENCES companies (id),
    depot_id            BIGINT              REFERENCES depots (id),
    vehicle_type_id     BIGINT              NOT NULL REFERENCES vehicle_types (id),
    registration_number VARCHAR(50)         NOT NULL UNIQUE,
    chassis_number      VARCHAR(100),
    engine_number       VARCHAR(100),
    manufacture_year    SMALLINT,
    color               VARCHAR(50),
    status              VARCHAR(20)         NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT chk_vehicles_status CHECK (status IN ('active', 'inactive', 'maintenance', 'decommissioned'))
);

CREATE INDEX idx_vehicles_company_id ON vehicles (company_id);
CREATE INDEX idx_vehicles_depot_id ON vehicles (depot_id);
CREATE INDEX idx_vehicles_vehicle_type_id ON vehicles (vehicle_type_id);
CREATE INDEX idx_vehicles_registration_number ON vehicles (registration_number);
CREATE INDEX idx_vehicles_status ON vehicles (status);
CREATE INDEX idx_vehicles_deleted_at ON vehicles (deleted_at);
