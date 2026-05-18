UPDATE adaptations SET status = 'active' WHERE status = 'en_curso';
ALTER TABLE adaptations ALTER COLUMN status SET DEFAULT 'active';
