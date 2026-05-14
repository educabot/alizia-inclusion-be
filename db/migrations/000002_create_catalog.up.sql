CREATE TABLE ramps (
    id BIGSERIAL PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    short_description VARCHAR(255),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE devices (
    id BIGSERIAL PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    ramp_id BIGINT NOT NULL REFERENCES ramps(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    image_url TEXT,
    qr_code VARCHAR(100),
    how_to_use TEXT,
    recommendations TEXT,
    rationale TEXT,
    classroom_benefit TEXT,
    needs_description TEXT,
    evaluation_criteria TEXT,
    quantity INTEGER NOT NULL DEFAULT 1,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, qr_code)
);

CREATE INDEX idx_ramps_organization_id ON ramps(organization_id);
CREATE INDEX idx_devices_organization_id ON devices(organization_id);
CREATE INDEX idx_devices_ramp_id ON devices(ramp_id);
