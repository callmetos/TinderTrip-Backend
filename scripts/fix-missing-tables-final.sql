-- =============================================================================
-- Fix Missing Tables - สร้างตารางที่ Backend ต้องการแต่ยังไม่มี
-- รัน script นี้ใน pgAdmin เพื่อเพิ่มตารางที่ขาด
-- =============================================================================

\c "NavMate-V2";

-- 1. email_verifications
CREATE TABLE IF NOT EXISTS email_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL,
    otp VARCHAR(6) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_email_verifications_email ON email_verifications(email);
CREATE INDEX IF NOT EXISTS idx_email_verifications_expires_at ON email_verifications(expires_at);

-- 2. food_preferences
CREATE TABLE IF NOT EXISTS food_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_category VARCHAR(50) NOT NULL,
    preference_level INTEGER NOT NULL CHECK (preference_level IN (1, 2, 3)),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (user_id, food_category)
);

CREATE INDEX IF NOT EXISTS idx_food_preferences_user_id ON food_preferences(user_id);

-- 3. notifications
CREATE TABLE IF NOT EXISTS notifications (
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

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(user_id, read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);

-- 4. travel_preferences
CREATE TABLE IF NOT EXISTS travel_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    travel_style VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (user_id, travel_style)
);

CREATE INDEX IF NOT EXISTS idx_travel_preferences_user_id ON travel_preferences(user_id);

-- 5. user_tags
CREATE TABLE IF NOT EXISTS user_tags (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_user_tags_user_id ON user_tags(user_id);
CREATE INDEX IF NOT EXISTS idx_user_tags_tag_id ON user_tags(tag_id);

-- =============================================================================
-- Triggers สำหรับ auto-update updated_at
-- =============================================================================

-- สร้าง function ถ้ายังไม่มี
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- เพิ่ม triggers
DROP TRIGGER IF EXISTS trigger_email_verifications_updated_at ON email_verifications;
CREATE TRIGGER trigger_email_verifications_updated_at 
    BEFORE UPDATE ON email_verifications
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_food_preferences_updated_at ON food_preferences;
CREATE TRIGGER trigger_food_preferences_updated_at 
    BEFORE UPDATE ON food_preferences
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_travel_preferences_updated_at ON travel_preferences;
CREATE TRIGGER trigger_travel_preferences_updated_at 
    BEFORE UPDATE ON travel_preferences
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- ตรวจสอบผลลัพธ์
-- =============================================================================

SELECT 'email_verifications' as table_name, COUNT(*) as exists 
FROM information_schema.tables 
WHERE table_name = 'email_verifications'
UNION ALL
SELECT 'food_preferences', COUNT(*) 
FROM information_schema.tables 
WHERE table_name = 'food_preferences'
UNION ALL
SELECT 'notifications', COUNT(*) 
FROM information_schema.tables 
WHERE table_name = 'notifications'
UNION ALL
SELECT 'travel_preferences', COUNT(*) 
FROM information_schema.tables 
WHERE table_name = 'travel_preferences'
UNION ALL
SELECT 'user_tags', COUNT(*) 
FROM information_schema.tables 
WHERE table_name = 'user_tags';

\echo '✅ สร้างตารางที่ขาดเรียบร้อยแล้ว!'
\echo 'ตารางทั้งหมดที่ Backend ต้องการ:'
\echo '  1. email_verifications'
\echo '  2. food_preferences'
\echo '  3. notifications'
\echo '  4. travel_preferences'
\echo '  5. user_tags'

