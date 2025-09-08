-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create database if not exists
SELECT 'CREATE DATABASE dashboard'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'dashboard')\gexec

-- Set default configuration
ALTER DATABASE dashboard SET timezone TO 'UTC';
ALTER DATABASE dashboard SET statement_timeout TO '30s';
ALTER DATABASE dashboard SET lock_timeout TO '10s';
ALTER DATABASE dashboard SET idle_in_transaction_session_timeout TO '60s';

-- Create read-only user for monitoring (optional)
CREATE USER dashboard_reader WITH PASSWORD 'readonly_password';
GRANT CONNECT ON DATABASE dashboard TO dashboard_reader;
GRANT USAGE ON SCHEMA public TO dashboard_reader;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO dashboard_reader;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO dashboard_reader;