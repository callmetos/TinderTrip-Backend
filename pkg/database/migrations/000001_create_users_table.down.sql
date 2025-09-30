-- Drop users table
DROP TABLE IF EXISTS users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop enum types
DROP TYPE IF EXISTS auth_provider;
DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS smoking;
DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS event_status;
DROP TYPE IF EXISTS member_role;
DROP TYPE IF EXISTS member_status;
DROP TYPE IF EXISTS swipe_direction;
