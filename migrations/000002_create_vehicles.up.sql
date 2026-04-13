CREATE TABLE IF NOT EXISTS vehicles (
    id                       UUID                     NOT NULL DEFAULT gen_random_uuid(),
    car_group_id             UUID                     NOT NULL,
    branch_id                UUID                     NOT NULL,
    vin                      TEXT                     NOT NULL,
    licence_plate            TEXT                     NOT NULL,
    brand                    TEXT                     NOT NULL,
    model                    TEXT                     NOT NULL,
    year                     INTEGER                  NOT NULL,
    colour                   TEXT,
    fuel_type                TEXT                     NOT NULL,
    transmission_type        TEXT                     NOT NULL,
    current_mileage          INTEGER                  NOT NULL DEFAULT 0,
    status                   TEXT                     NOT NULL,
    designation              TEXT                     NOT NULL,
    acquisition_date         DATE                     NOT NULL,
    ownership_type           TEXT                     NOT NULL,
    lease_details            TEXT,
    insurance_policy_number  TEXT,
    insurance_expiry_date    DATE,
    registration_expiry_date DATE,
    last_inspection_date     DATE,
    next_inspection_due_date DATE,
    notes                    TEXT,
    created_at               TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at               TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at               TIMESTAMP WITH TIME ZONE,
    created_by               TEXT                     NOT NULL,
    updated_by               TEXT                     NOT NULL,
    deleted                  BOOL                     NOT NULL DEFAULT false,

    CONSTRAINT vehicles_pkey PRIMARY KEY (id),
    CONSTRAINT vehicles_vin_unique UNIQUE (vin),
    CONSTRAINT vehicles_licence_plate_unique UNIQUE (licence_plate),
    CONSTRAINT vehicles_status_check CHECK (status IN (
        'available',
        'on_rent',
        'needs_cleaning',
        'needs_inspection',
        'under_maintenance',
        'unavailable',
        'decommissioned'
    )),
    CONSTRAINT vehicles_car_group_fkey FOREIGN KEY (car_group_id) REFERENCES car_groups (id),
    CONSTRAINT vehicles_current_mileage_check CHECK (current_mileage >= 0)
);

CREATE INDEX IF NOT EXISTS vehicles_car_group_id_idx ON vehicles (car_group_id);
CREATE INDEX IF NOT EXISTS vehicles_branch_id_idx ON vehicles (branch_id);
CREATE INDEX IF NOT EXISTS vehicles_fuel_type_idx ON vehicles (fuel_type);
CREATE INDEX IF NOT EXISTS vehicles_transmission_type_idx ON vehicles (transmission_type);
CREATE INDEX IF NOT EXISTS vehicles_status_idx ON vehicles (status);
CREATE INDEX IF NOT EXISTS vehicles_designation_idx ON vehicles (designation);
