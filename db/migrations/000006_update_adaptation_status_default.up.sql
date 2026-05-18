ALTER TABLE adaptations ALTER COLUMN status SET DEFAULT 'en_curso';
UPDATE adaptations SET status = 'en_curso' WHERE status = 'active';
