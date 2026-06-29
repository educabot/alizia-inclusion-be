CREATE TABLE device_videos (
    id         BIGSERIAL PRIMARY KEY,
    device_id  BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    title      VARCHAR(255),
    url        TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_videos_device_id ON device_videos(device_id);
