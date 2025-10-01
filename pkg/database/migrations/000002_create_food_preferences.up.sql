-- Create food_preferences table
CREATE TABLE IF NOT EXISTS food_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_category VARCHAR(50) NOT NULL,
    preference_level INTEGER NOT NULL CHECK (preference_level IN (1, 2, 3)), -- 1=dislike, 2=neutral, 3=love
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, food_category)
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_food_preferences_user_id ON food_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_food_preferences_category ON food_preferences(food_category);

-- Insert default food categories
INSERT INTO food_preferences (user_id, food_category, preference_level) 
SELECT 
    u.id,
    'thai_food',
    2 -- neutral as default
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM food_preferences fp 
    WHERE fp.user_id = u.id AND fp.food_category = 'thai_food'
);

INSERT INTO food_preferences (user_id, food_category, preference_level) 
SELECT 
    u.id,
    'japanese_food',
    2 -- neutral as default
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM food_preferences fp 
    WHERE fp.user_id = u.id AND fp.food_category = 'japanese_food'
);

INSERT INTO food_preferences (user_id, food_category, preference_level) 
SELECT 
    u.id,
    'chinese_food',
    2 -- neutral as default
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM food_preferences fp 
    WHERE fp.user_id = u.id AND fp.food_category = 'chinese_food'
);

INSERT INTO food_preferences (user_id, food_category, preference_level) 
SELECT 
    u.id,
    'international_food',
    2 -- neutral as default
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM food_preferences fp 
    WHERE fp.user_id = u.id AND fp.food_category = 'international_food'
);

INSERT INTO food_preferences (user_id, food_category, preference_level) 
SELECT 
    u.id,
    'halal_food',
    2 -- neutral as default
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM food_preferences fp 
    WHERE fp.user_id = u.id AND fp.food_category = 'halal_food'
);

INSERT INTO food_preferences (user_id, food_category, preference_level) 
SELECT 
    u.id,
    'buffet',
    2 -- neutral as default
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM food_preferences fp 
    WHERE fp.user_id = u.id AND fp.food_category = 'buffet'
);

INSERT INTO food_preferences (user_id, food_category, preference_level) 
SELECT 
    u.id,
    'bbq_grill',
    2 -- neutral as default
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM food_preferences fp 
    WHERE fp.user_id = u.id AND fp.food_category = 'bbq_grill'
);
