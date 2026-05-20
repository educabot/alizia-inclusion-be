CREATE TABLE IF NOT EXISTS adaptation_devices (
    adaptation_id BIGINT NOT NULL REFERENCES adaptations(id) ON DELETE CASCADE,
    device_id     BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    PRIMARY KEY (adaptation_id, device_id)
);

CREATE INDEX IF NOT EXISTS idx_adaptation_devices_adaptation ON adaptation_devices(adaptation_id);
CREATE INDEX IF NOT EXISTS idx_adaptation_devices_device ON adaptation_devices(device_id);

ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS adaptation_type VARCHAR(50) DEFAULT '';
