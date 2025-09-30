-- =============================
-- Extensions
-- =============================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- =============================
-- Enums
-- =============================
DO $$ BEGIN
  CREATE TYPE auth_provider AS ENUM ('password','google');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE gender AS ENUM ('male','female','nonbinary','prefer_not_say');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE smoking AS ENUM ('no','yes','occasionally');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE event_type AS ENUM ('meal','one_day_trip','overnight');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE event_status AS ENUM ('draft','active','closed','cancelled','completed');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE member_role AS ENUM ('creator','participant');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE member_status AS ENUM ('pending','confirmed','declined','kicked','left');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE swipe_direction AS ENUM ('like','pass');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- =============================
-- Users & Auth
-- =============================
CREATE TABLE users (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email          CITEXT UNIQUE,
  provider       auth_provider NOT NULL,
  password_hash  TEXT,
  google_id      TEXT,
  display_name   TEXT,
  last_login_at  TIMESTAMPTZ,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at     TIMESTAMPTZ,
  CONSTRAINT ck_pwd_required CHECK ((provider <> 'password') OR password_hash IS NOT NULL),
  CONSTRAINT ck_google_required CHECK ((provider <> 'google') OR google_id IS NOT NULL)
);

CREATE UNIQUE INDEX ux_users_google_id ON users(google_id)
  WHERE provider='google' AND google_id IS NOT NULL;

CREATE TABLE password_resets (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token      TEXT UNIQUE NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =============================
-- Profiles
-- =============================
CREATE TABLE user_profiles (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  bio            TEXT,
  languages      JSONB,
  date_of_birth  DATE,
  gender         gender,
  job_title      TEXT,
  smoking        smoking,
  interests_note TEXT,
  avatar_url     TEXT,
  home_location  TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at     TIMESTAMPTZ
);

-- =============================
-- Tags & User Tags
-- =============================
CREATE TABLE tags (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       CITEXT UNIQUE NOT NULL,
  kind       TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_tags (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tag_id  UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, tag_id)
);

-- =============================
-- Preferences
-- =============================
CREATE TABLE pref_availability (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  mon         BOOLEAN NOT NULL DEFAULT false,
  tue         BOOLEAN NOT NULL DEFAULT false,
  wed         BOOLEAN NOT NULL DEFAULT false,
  thu         BOOLEAN NOT NULL DEFAULT false,
  fri         BOOLEAN NOT NULL DEFAULT false,
  sat         BOOLEAN NOT NULL DEFAULT false,
  sun         BOOLEAN NOT NULL DEFAULT false,
  all_day     BOOLEAN NOT NULL DEFAULT true,
  morning     BOOLEAN NOT NULL DEFAULT false,
  afternoon   BOOLEAN NOT NULL DEFAULT false,
  time_range  TSTZRANGE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE pref_budget (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  meal_min       INT,
  meal_max       INT,
  daytrip_min    INT,
  daytrip_max    INT,
  overnight_min  INT,
  overnight_max  INT,
  unlimited      BOOLEAN NOT NULL DEFAULT false,
  currency       TEXT NOT NULL DEFAULT 'THB',
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =============================
-- Events
-- =============================
CREATE TABLE events (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  creator_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title            TEXT NOT NULL,
  description      TEXT,
  event_type       event_type NOT NULL DEFAULT 'meal',
  address_text     TEXT,
  lat              DOUBLE PRECISION,
  lng              DOUBLE PRECISION,
  start_at         TIMESTAMPTZ,
  end_at           TIMESTAMPTZ,
  capacity         INT CHECK (capacity IS NULL OR capacity >= 1),
  status           event_status NOT NULL DEFAULT 'draft',
  cover_image_url  TEXT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at       TIMESTAMPTZ
);

CREATE TABLE event_photos (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id   UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  url        TEXT NOT NULL,
  sort_no    INT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE event_categories (
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  tag_id   UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (event_id, tag_id)
);

CREATE TABLE event_tags (
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  tag_id   UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (event_id, tag_id)
);

-- =============================
-- Chat (ต้องมาก่อน event_members เพราะมี FK)
-- =============================
CREATE TABLE chat_rooms (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id   UUID NOT NULL UNIQUE REFERENCES events(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE chat_messages (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id      UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
  sender_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  body         TEXT,
  message_type TEXT,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =============================
-- Event Members
-- =============================
CREATE TABLE event_members (
  event_id     UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role         member_role NOT NULL DEFAULT 'participant',
  status       member_status NOT NULL DEFAULT 'pending',
  joined_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  confirmed_at TIMESTAMPTZ,
  left_at      TIMESTAMPTZ,
  note         TEXT,
  confirmation_message_id UUID REFERENCES chat_messages(id),
  PRIMARY KEY (event_id, user_id)
);

-- =============================
-- Event Swipes
-- =============================
CREATE TABLE event_swipes (
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event_id   UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  direction  swipe_direction NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, event_id)
);

-- =============================
-- History
-- =============================
CREATE TABLE user_event_history (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id     UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  completed    BOOLEAN NOT NULL DEFAULT false,
  completed_at TIMESTAMPTZ,
  UNIQUE (event_id, user_id)
);

-- =============================
-- Logs
-- =============================
CREATE TABLE audit_logs (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  entity_table  TEXT NOT NULL,
  entity_id     UUID,
  action        TEXT NOT NULL,
  before_data   JSONB,
  after_data    JSONB,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE api_logs (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id  TEXT,
  user_id     UUID REFERENCES users(id) ON DELETE SET NULL,
  method      TEXT,
  path        TEXT,
  status      INT,
  duration_ms INT,
  ip_address  INET,
  user_agent  TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =============================
-- Indexes
-- =============================
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_provider ON users(provider);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);
CREATE INDEX idx_password_resets_token ON password_resets(token);
CREATE INDEX idx_password_resets_expires_at ON password_resets(expires_at);
CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX idx_user_profiles_gender ON user_profiles(gender);
CREATE INDEX idx_user_profiles_deleted_at ON user_profiles(deleted_at);
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_kind ON tags(kind);
CREATE INDEX idx_events_creator_id ON events(creator_id);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_start_at ON events(start_at);
CREATE INDEX idx_events_deleted_at ON events(deleted_at);
CREATE INDEX idx_event_photos_event_id ON event_photos(event_id);
CREATE INDEX idx_event_photos_sort_no ON event_photos(sort_no);
CREATE INDEX idx_chat_rooms_event_id ON chat_rooms(event_id);
CREATE INDEX idx_chat_messages_room_id ON chat_messages(room_id);
CREATE INDEX idx_chat_messages_sender_id ON chat_messages(sender_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);
CREATE INDEX idx_event_members_user_id ON event_members(user_id);
CREATE INDEX idx_event_members_status ON event_members(status);
CREATE INDEX idx_event_swipes_user_id ON event_swipes(user_id);
CREATE INDEX idx_event_swipes_event_id ON event_swipes(event_id);
CREATE INDEX idx_user_event_history_user_id ON user_event_history(user_id);
CREATE INDEX idx_user_event_history_event_id ON user_event_history(event_id);
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_entity_table ON audit_logs(entity_table);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at);

-- =============================
-- Triggers
-- =============================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at 
    BEFORE UPDATE ON user_profiles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pref_availability_updated_at 
    BEFORE UPDATE ON pref_availability 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pref_budget_updated_at 
    BEFORE UPDATE ON pref_budget 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
