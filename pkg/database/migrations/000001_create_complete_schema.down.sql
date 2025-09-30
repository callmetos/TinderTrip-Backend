-- Drop all tables in reverse order
DROP TABLE IF EXISTS api_logs;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS user_event_history;
DROP TABLE IF EXISTS event_swipes;
DROP TABLE IF EXISTS event_members;
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_rooms;
DROP TABLE IF EXISTS event_tags;
DROP TABLE IF EXISTS event_categories;
DROP TABLE IF EXISTS event_photos;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS pref_budget;
DROP TABLE IF EXISTS pref_availability;
DROP TABLE IF EXISTS user_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS password_resets;
DROP TABLE IF EXISTS users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop enum types
DROP TYPE IF EXISTS swipe_direction;
DROP TYPE IF EXISTS member_status;
DROP TYPE IF EXISTS member_role;
DROP TYPE IF EXISTS event_status;
DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS smoking;
DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS auth_provider;
