-- Drop indexes first
DROP INDEX IF EXISTS idx_articles_status;
DROP INDEX IF EXISTS idx_articles_admin_id;
DROP INDEX IF EXISTS idx_articles_slug;
DROP INDEX IF EXISTS idx_devices_device_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_users_admin_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables in reverse order (handle foreign key dependencies)
DROP TABLE IF EXISTS articles CASCADE;
DROP TABLE IF EXISTS devices CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS admins CASCADE;