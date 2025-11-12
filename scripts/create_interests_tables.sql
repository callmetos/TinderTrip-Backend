-- ============================================
-- Create Unified Interests System
-- ============================================
-- This script creates the interests master table and related tables
-- for unified user and event interests

-- Create unified interests master table
-- This table stores all interests from UI (Restaurant, Cafe, Activity, Pub & Bar, Sport)
CREATE TABLE IF NOT EXISTS interests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    icon TEXT,
    category VARCHAR(50) NOT NULL, -- 'restaurant', 'cafe', 'activity', 'pub_bar', 'sport'
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for category filtering
CREATE INDEX IF NOT EXISTS idx_interests_category ON interests(category) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_interests_active ON interests(is_active) WHERE is_active = TRUE;

-- Create user_interests table (unified preferences)
-- This replaces separate travel_preferences and food_preferences
CREATE TABLE IF NOT EXISTS user_interests (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    interest_id UUID NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, interest_id)
);

-- Create index for user queries
CREATE INDEX IF NOT EXISTS idx_user_interests_user_id ON user_interests(user_id);
CREATE INDEX IF NOT EXISTS idx_user_interests_interest_id ON user_interests(interest_id);

-- Create event_interests table (unified interests for events)
CREATE TABLE IF NOT EXISTS event_interests (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    interest_id UUID NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id, interest_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_event_interests_event_id ON event_interests(event_id);
CREATE INDEX IF NOT EXISTS idx_event_interests_interest_id ON event_interests(interest_id);

-- Create trigger for updated_at (if function exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE TRIGGER update_interests_updated_at 
            BEFORE UPDATE ON interests 
            FOR EACH ROW 
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- ============================================
-- Seed interests data from UI
-- ============================================

-- Restaurant category (for food-related interests)
INSERT INTO interests (code, display_name, icon, category, sort_order, is_active) VALUES
    ('fast_food', 'Fast Food', '🍟', 'restaurant', 1, true),
    ('noodles', 'Noodles', '🍜', 'restaurant', 2, true),
    ('grill', 'Grill', '♨️', 'restaurant', 3, true),
    ('pasta', 'Pasta', '🍝', 'restaurant', 4, true),
    ('dim_sum', 'Dim Sum', '🥟', 'restaurant', 5, true),
    ('indian_food', 'Indian Food', '🇮🇳', 'restaurant', 6, true),
    ('salads', 'Salads', '🥗', 'restaurant', 7, true),
    ('japanese_food', 'Japanese Food', '🇯🇵', 'restaurant', 8, true),
    ('izakaya', 'Izakaya', '🍺', 'restaurant', 9, true),
    ('muu_kra_ta', 'Muu Kra Ta', '🐷', 'restaurant', 10, true),
    ('street_food', 'Street Food', '🥡', 'restaurant', 11, true),
    ('pork', 'Pork', '🐖', 'restaurant', 12, true),
    ('pizza', 'Pizza', '🍕', 'restaurant', 13, true),
    ('vegan', 'Vegan', '🥬', 'restaurant', 14, true),
    ('chinese_food', 'Chinese Food', '🇨🇳', 'restaurant', 15, true),
    ('sushi', 'Sushi', '🍣', 'restaurant', 16, true),
    ('fine_dining', 'Fine Dining', '🍽️', 'restaurant', 17, true),
    ('halal', 'Halal', '☪️', 'restaurant', 18, true),
    ('burger', 'Burger', '🍔', 'restaurant', 19, true),
    ('korean_food', 'Korean Food', '🇰🇷', 'restaurant', 20, true),
    ('buffet', 'Buffet', '😊', 'restaurant', 21, true),
    ('ramen', 'Ramen', '🍜', 'restaurant', 22, true),
    ('bbq', 'BBQ', '🔥', 'restaurant', 23, true),
    ('meat', 'Meat', '🥩', 'restaurant', 24, true),
    ('healthy_food', 'Healthy Food', '🥑', 'restaurant', 25, true),
    ('shabu_sukiyaki_hot_pot', 'Shabu / Sukiyaki / Hot Pot', '🍲', 'restaurant', 26, true),
    ('omakase', 'Omakase', '🍣', 'restaurant', 27, true),
    ('seafood', 'Seafood', '🦀', 'restaurant', 28, true)
ON CONFLICT (code) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    icon = EXCLUDED.icon,
    category = EXCLUDED.category,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Cafe category
INSERT INTO interests (code, display_name, icon, category, sort_order, is_active) VALUES
    ('bubble_tea', 'Bubble Tea', '🥤', 'cafe', 1, true),
    ('bingsu', 'Bingsu', '🍧', 'cafe', 2, true),
    ('matcha', 'Matcha', '🍵', 'cafe', 3, true),
    ('bakery_cake', 'Bakery / Cake', '🍞', 'cafe', 4, true),
    ('ice_cream', 'Ice Cream', '🍦', 'cafe', 5, true),
    ('pancakes', 'Pancakes', '🥞', 'cafe', 6, true),
    ('coffee', 'Coffee', '☕', 'cafe', 7, true)
ON CONFLICT (code) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    icon = EXCLUDED.icon,
    category = EXCLUDED.category,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Activity category
INSERT INTO interests (code, display_name, icon, category, sort_order, is_active) VALUES
    ('chilling', 'Chilling', '😌', 'activity', 1, true),
    ('painting', 'Painting', '🎨', 'activity', 2, true),
    ('baking', 'Baking', '🍪', 'activity', 3, true),
    ('investing', 'Investing', '📈', 'activity', 4, true),
    ('fan_meet', 'Fan Meet', '👥', 'activity', 5, true),
    ('shopping', 'Shopping', '🛍️', 'activity', 6, true),
    ('wakeboard', 'Wakeboard', '🏄', 'activity', 7, true),
    ('laser_tag', 'Laser Tag', '🔫', 'activity', 8, true),
    ('superstition', 'Superstition / Black Magic / Mutelu', '🔮', 'activity', 9, true),
    ('bb_gun', 'BB Gun', '🔫', 'activity', 10, true),
    ('travel', 'Travel', '🌴', 'activity', 11, true),
    ('photography', 'Photography', '📸', 'activity', 12, true),
    ('temple', 'Temple', '🛕', 'activity', 13, true),
    ('night_market', 'Night Market', '🧺', 'activity', 14, true),
    ('park', 'Park', '🌳', 'activity', 15, true),
    ('amusement_park', 'Amusement Park', '🎢', 'activity', 16, true),
    ('movies', 'Movies', '🎞️', 'activity', 17, true),
    ('karaoke', 'Karaoke', '🎤', 'activity', 18, true),
    ('running', 'Running', '🏃', 'activity', 19, true),
    ('art_gallery', 'Art Gallery', '🖼️', 'activity', 20, true),
    ('archery', 'Archery', '🏹', 'activity', 21, true),
    ('scuba_diving', 'Scuba Diving / Snorkeling', '🤿', 'activity', 22, true),
    ('reading', 'Reading', '📚', 'activity', 23, true),
    ('skateboard', 'Skateboard', '🛹', 'activity', 24, true),
    ('walking', 'Walking', '🚶', 'activity', 25, true),
    ('volunteer', 'Volunteer', '🙋‍♀️', 'activity', 26, true),
    ('boardgame', 'Boardgame', '🎲', 'activity', 27, true),
    ('paintball', 'Paintball', '🔫', 'activity', 28, true),
    ('museum', 'Museum', '🏛️', 'activity', 29, true),
    ('flower_arrangement', 'Flower Arrangement', '💐', 'activity', 30, true),
    ('kpop', 'K-POP', '🕺', 'activity', 31, true),
    ('concert', 'Concert', '🎶', 'activity', 32, true),
    ('aquarium', 'Aquarium', '🐠', 'activity', 33, true)
ON CONFLICT (code) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    icon = EXCLUDED.icon,
    category = EXCLUDED.category,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Pub & Bar category
INSERT INTO interests (code, display_name, icon, category, sort_order, is_active) VALUES
    ('wine', 'Wine', '🍷', 'pub_bar', 1, true),
    ('ratchathewi', 'Ratchathewi', '🥃', 'pub_bar', 2, true),
    ('khaosan_road', 'Khaosan Road', '🍚', 'pub_bar', 3, true),
    ('thonglor', 'Thonglor', '🥃', 'pub_bar', 4, true),
    ('thai_music', 'Thai Music', '🎵', 'pub_bar', 5, true),
    ('edm_music', 'EDM Music', '🎵', 'pub_bar', 6, true),
    ('heartbroken', 'Heartbroken', '💔', 'pub_bar', 7, true),
    ('beer_tower', 'Beer tower', '🍺', 'pub_bar', 8, true),
    ('jazz', 'Jazz', '🎷', 'pub_bar', 9, true),
    ('cocktail_bar', 'Cocktail Bar', '🍸', 'pub_bar', 10, true),
    ('live_music', 'Live Music', '🎵', 'pub_bar', 11, true),
    ('rca_plaza', 'RCA Plaza', '🥃', 'pub_bar', 12, true),
    ('kpop_music', 'K-Pop Music', '🎵', 'pub_bar', 13, true),
    ('pubs_bars', 'Pubs & Bars', '🍻', 'pub_bar', 14, true),
    ('rooftop', 'Rooftop', '🏙️', 'pub_bar', 15, true),
    ('alcohol', 'Alcohol', '🥂', 'pub_bar', 16, true)
ON CONFLICT (code) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    icon = EXCLUDED.icon,
    category = EXCLUDED.category,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Sport category
INSERT INTO interests (code, display_name, icon, category, sort_order, is_active) VALUES
    ('football', 'Football', '⚽', 'sport', 1, true),
    ('sports', 'Sports', '🤸', 'sport', 2, true),
    ('snooker', 'Snooker', '🎱', 'sport', 3, true),
    ('rock_climbing', 'Rock climbing', '🧗', 'sport', 4, true),
    ('golf', 'Golf', '⛳', 'sport', 5, true),
    ('boxing', 'Boxing', '🥊', 'sport', 6, true),
    ('fitness', 'Fitness', '🏋️', 'sport', 7, true),
    ('basketball', 'Basketball', '🏀', 'sport', 8, true),
    ('bowling', 'Bowling', '🎳', 'sport', 9, true),
    ('tennis', 'Tennis', '🎾', 'sport', 10, true),
    ('badminton', 'Badminton', '🏸', 'sport', 11, true),
    ('volleyball', 'Volleyball', '🏐', 'sport', 12, true),
    ('racquet', 'Racquet', '🏏', 'sport', 13, true),
    ('table_tennis', 'Table Tennis', '🏓', 'sport', 14, true),
    ('yoga', 'Yoga', '🧘', 'sport', 15, true)
ON CONFLICT (code) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    icon = EXCLUDED.icon,
    category = EXCLUDED.category,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- ============================================
-- Summary
-- ============================================
-- Total interests: 99 items
-- - Restaurant: 28 items
-- - Cafe: 7 items
-- - Activity: 33 items
-- - Pub & Bar: 16 items
-- - Sport: 15 items

