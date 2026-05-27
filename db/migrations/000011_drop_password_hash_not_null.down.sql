-- Reverts 000011_drop_password_hash_not_null: restores the original NOT NULL constraint.
-- Note: only applicable if every users.password_hash row is non-NULL.
ALTER TABLE users ALTER COLUMN password_hash DROP DEFAULT;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
