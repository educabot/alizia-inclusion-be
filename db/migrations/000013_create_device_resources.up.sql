CREATE TABLE IF NOT EXISTS device_resources (
    id        BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    title     VARCHAR(255) NOT NULL,
    file_url  TEXT NOT NULL,
    file_type VARCHAR(50) NOT NULL DEFAULT 'pdf',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_device_resources_device ON device_resources(device_id);
