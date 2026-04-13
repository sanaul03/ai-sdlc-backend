CREATE TABLE IF NOT EXISTS depots (
    id         BIGSERIAL PRIMARY KEY,
    company_id BIGINT              NOT NULL REFERENCES companies (id),
    name       VARCHAR(255)        NOT NULL,
    code       VARCHAR(50)         NOT NULL,
    address    TEXT,
    latitude   NUMERIC(10, 8),
    longitude  NUMERIC(11, 8),
    -- status: 1 = active, 0 = inactive
    status     SMALLINT            NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (company_id, code),
    CONSTRAINT chk_depots_status CHECK (status IN (0, 1))
);

CREATE INDEX idx_depots_company_id ON depots (company_id);
CREATE INDEX idx_depots_status ON depots (status);
CREATE INDEX idx_depots_deleted_at ON depots (deleted_at);
