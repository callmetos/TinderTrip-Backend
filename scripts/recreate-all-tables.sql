-- =============================================================================
-- TinderTrip Database - DROP and RECREATE ALL TABLES
-- WARNING: This will DELETE ALL DATA!
-- =============================================================================

\c "NavMate-V2";

-- =============================================================================
-- DROP EVERYTHING
-- =============================================================================

-- Drop all tables (with CASCADE to handle dependencies)
DROP TABLE IF EXISTS api_logs CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS user_event_histories CASCADE;
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS chat_rooms CASCADE;
DROP TABLE IF EXISTS event_tags CASCADE;
DROP TABLE IF EXISTS event_categories CASCADE;
DROP TABLE IF EXISTS event_swipes CASCADE;
DROP TABLE IF EXISTS event_members CASCADE;
DROP TABLE IF EXISTS event_photos CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS travel_preferences CASCADE;
DROP TABLE IF EXISTS food_preferences CASCADE;
DROP TABLE IF EXISTS pref_budgets CASCADE;
DROP TABLE IF EXISTS pref_availabilities CASCADE;
DROP TABLE IF EXISTS user_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS email_verifications CASCADE;
DROP TABLE IF EXISTS password_resets CASCADE;
DROP TABLE IF EXISTS user_profiles CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop all views
DROP VIEW IF EXISTS v_active_users CASCADE;
DROP VIEW IF EXISTS v_active_events CASCADE;

-- Drop all types
DROP TYPE IF EXISTS auth_provider CASCADE;
DROP TYPE IF EXISTS gender CASCADE;
DROP TYPE IF EXISTS smoking CASCADE;
DROP TYPE IF EXISTS event_type CASCADE;
DROP TYPE IF EXISTS event_status CASCADE;
DROP TYPE IF EXISTS member_role CASCADE;
DROP TYPE IF EXISTS member_status CASCADE;
DROP TYPE IF EXISTS swipe_direction CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
DROP FUNCTION IF EXISTS cleanup_old_logs() CASCADE;

-- =============================================================================
-- CREATE EXTENSIONS
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- =============================================================================
-- CREATE ENUM TYPES
-- =============================================================================

CREATE TYPE auth_provider AS ENUM ('password', 'google', 'apple', 'facebook');
CREATE TYPE gender AS ENUM ('male', 'female', 'nonbinary', 'prefer_not_say');
CREATE TYPE smoking AS ENUM ('no', 'yes', 'occasionally');
CREATE TYPE event_type AS ENUM ('meal', 'daytrip', 'overnight', 'activity', 'other');
CREATE TYPE event_status AS ENUM ('published', 'cancelled', 'completed');
CREATE TYPE member_role AS ENUM ('creator', 'participant');
CREATE TYPE member_status AS ENUM ('pending', 'confirmed', 'declined', 'kicked', 'left');
CREATE TYPE swipe_direction AS ENUM ('like', 'pass');

-- =============================================================================
-- CREATE TABLES
-- =============================================================================

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT UNIQUE,
    provider auth_provider NOT NULL,
    password_hash TEXT,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    google_id TEXT,
    display_name TEXT,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- User Profiles
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    languages TEXT,
    date_of_birth DATE,
    gender gender,
    job_title TEXT,
    smoking smoking,
    interests_note TEXT,
    avatar_url TEXT,
    home_location TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Password Resets
CREATE TABLE password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Email Verifications
CREATE TABLE email_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL,
    otp VARCHAR(6) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Tags
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name CITEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User Tags (many-to-many)
CREATE TABLE user_tags (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tag_id)
);

-- Availability Preferences
CREATE TABLE pref_availabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    mon BOOLEAN NOT NULL DEFAULT FALSE,
    tue BOOLEAN NOT NULL DEFAULT FALSE,
    wed BOOLEAN NOT NULL DEFAULT FALSE,
    thu BOOLEAN NOT NULL DEFAULT FALSE,
    fri BOOLEAN NOT NULL DEFAULT FALSE,
    sat BOOLEAN NOT NULL DEFAULT FALSE,
    sun BOOLEAN NOT NULL DEFAULT FALSE,
    all_day BOOLEAN NOT NULL DEFAULT TRUE,
    morning BOOLEAN NOT NULL DEFAULT FALSE,
    afternoon BOOLEAN NOT NULL DEFAULT FALSE,
    time_range TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Budget Preferences
CREATE TABLE pref_budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    meal_min INTEGER,
    meal_max INTEGER,
    daytrip_min INTEGER,
    daytrip_max INTEGER,
    overnight_min INTEGER,
    overnight_max INTEGER,
    unlimited BOOLEAN NOT NULL DEFAULT FALSE,
    currency TEXT NOT NULL DEFAULT 'THB',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Food Preferences
CREATE TABLE food_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_category VARCHAR(50) NOT NULL,
    preference_level INTEGER NOT NULL CHECK (preference_level IN (1, 2, 3)),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (user_id, food_category)
);

-- Travel Preferences
CREATE TABLE travel_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    travel_style VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (user_id, travel_style)
);

-- Events
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    event_type event_type NOT NULL DEFAULT 'meal',
    address_text TEXT,
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    capacity INTEGER CHECK (capacity IS NULL OR capacity >= 1),
    budget_min INTEGER CHECK (budget_min IS NULL OR budget_min >= 0),
    budget_max INTEGER CHECK (budget_max IS NULL OR budget_max >= 0),
    currency VARCHAR(3) DEFAULT 'THB',
    status event_status NOT NULL DEFAULT 'published',
    cover_image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Event Photos
CREATE TABLE event_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    sort_no INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Event Members
CREATE TABLE event_members (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role member_role NOT NULL DEFAULT 'participant',
    status member_status NOT NULL DEFAULT 'pending',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ,
    left_at TIMESTAMPTZ,
    note TEXT,
    confirmation_message_id UUID,
    PRIMARY KEY (event_id, user_id)
);

-- Event Swipes
CREATE TABLE event_swipes (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    direction swipe_direction NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, event_id)
);

-- Event Categories (many-to-many)
CREATE TABLE event_categories (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, tag_id)
);

-- Event Tags (many-to-many)
CREATE TABLE event_tags (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, tag_id)
);

-- Chat Rooms
CREATE TABLE chat_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL UNIQUE REFERENCES events(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Chat Messages
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT,
    message_type TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add FK for confirmation_message_id
ALTER TABLE event_members 
    ADD CONSTRAINT fk_event_members_confirmation_message 
    FOREIGN KEY (confirmation_message_id) 
    REFERENCES chat_messages(id) ON DELETE SET NULL;

-- User Event History
CREATE TABLE user_event_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    type TEXT NOT NULL,
    data JSONB,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ
);

-- Audit Logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    entity_table TEXT NOT NULL,
    entity_id UUID,
    action TEXT NOT NULL,
    before_data JSONB,
    after_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- API Logs
CREATE TABLE api_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id TEXT,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    method TEXT,
    path TEXT,
    status INTEGER,
    duration_ms INTEGER,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- CREATE INDEXES (Production Scale)
-- =============================================================================

-- Users
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_provider ON users(provider);
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- User Profiles
CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);

-- Password & Email
CREATE INDEX idx_password_resets_token ON password_resets(token);
CREATE INDEX idx_password_resets_expires_at ON password_resets(expires_at);
CREATE INDEX idx_email_verifications_email ON email_verifications(email);
CREATE INDEX idx_email_verifications_expires_at ON email_verifications(expires_at);

-- Preferences
CREATE INDEX idx_pref_availabilities_user_id ON pref_availabilities(user_id);
CREATE INDEX idx_pref_budgets_user_id ON pref_budgets(user_id);
CREATE INDEX idx_food_preferences_user_id ON food_preferences(user_id);
CREATE INDEX idx_travel_preferences_user_id ON travel_preferences(user_id);

-- Tags
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_kind ON tags(kind);
CREATE INDEX idx_user_tags_user_id ON user_tags(user_id);
CREATE INDEX idx_user_tags_tag_id ON user_tags(tag_id);

-- Events
CREATE INDEX idx_events_creator_id ON events(creator_id);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_start_at ON events(start_at);
CREATE INDEX idx_events_created_at ON events(created_at DESC);
CREATE INDEX idx_events_location ON events(lat, lng) WHERE lat IS NOT NULL;

-- Event Relations
CREATE INDEX idx_event_photos_event_id ON event_photos(event_id);
CREATE INDEX idx_event_members_event_id ON event_members(event_id);
CREATE INDEX idx_event_members_user_id ON event_members(user_id);
CREATE INDEX idx_event_members_status ON event_members(status);
CREATE INDEX idx_event_swipes_user_id ON event_swipes(user_id);
CREATE INDEX idx_event_swipes_event_id ON event_swipes(event_id);
CREATE INDEX idx_event_swipes_created_at ON event_swipes(created_at DESC);
CREATE INDEX idx_event_categories_event_id ON event_categories(event_id);
CREATE INDEX idx_event_categories_tag_id ON event_categories(tag_id);
CREATE INDEX idx_event_tags_event_id ON event_tags(event_id);
CREATE INDEX idx_event_tags_tag_id ON event_tags(tag_id);

-- Chat
CREATE INDEX idx_chat_rooms_event_id ON chat_rooms(event_id);
CREATE INDEX idx_chat_messages_room_id ON chat_messages(room_id);
CREATE INDEX idx_chat_messages_sender_id ON chat_messages(sender_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(room_id, created_at DESC);

-- History
CREATE INDEX idx_user_event_histories_event_id ON user_event_histories(event_id);
CREATE INDEX idx_user_event_histories_user_id ON user_event_histories(user_id);

-- Notifications
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(user_id, read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- Logs (for monitoring and cleanup)
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_table, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at DESC);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_status ON api_logs(status);

-- =============================================================================
-- TRIGGERS
-- =============================================================================

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_pref_availabilities_updated_at BEFORE UPDATE ON pref_availabilities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_pref_budgets_updated_at BEFORE UPDATE ON pref_budgets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- SEED DATA (Optional)
-- =============================================================================

INSERT INTO tags (name, kind) VALUES
    ('Adventure', 'interest'),
    ('Food', 'interest'),
    ('Culture', 'interest'),
    ('Nature', 'interest'),
    ('Photography', 'interest'),
    ('Music', 'interest'),
    ('Sports', 'interest'),
    ('Art', 'interest'),
    ('Shopping', 'activity'),
    ('Nightlife', 'activity'),
    ('Beach', 'location'),
    ('Mountain', 'location'),
    ('City', 'location'),
    ('Thai Food', 'food'),
    ('Japanese Food', 'food'),
    ('International', 'food')
ON CONFLICT (name) DO NOTHING;

-- =============================================================================
-- VACUUM
-- =============================================================================

VACUUM ANALYZE;

-- =============================================================================
-- COMPLETION
-- =============================================================================

SELECT 
    'Tables created: ' || COUNT(*)::TEXT 
FROM information_schema.tables 
WHERE table_schema = 'public' AND table_type = 'BASE TABLE';

SELECT 
    'Indexes created: ' || COUNT(*)::TEXT 
FROM pg_indexes 
WHERE schemaname = 'public';

\echo '✅ Database recreated successfully!'

