CREATE TABLE IF NOT EXISTS companies (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255)        NOT NULL,
    code       VARCHAR(50)         NOT NULL UNIQUE,
    address    TEXT,
    phone      VARCHAR(50),
    email      VARCHAR(255),
    -- status: 1 = active, 0 = inactive
    status     SMALLINT            NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_companies_status CHECK (status IN (0, 1))
);

CREATE INDEX idx_companies_code ON companies (code);
CREATE INDEX idx_companies_status ON companies (status);
CREATE INDEX idx_companies_deleted_at ON companies (deleted_at);
