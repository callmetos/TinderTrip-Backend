-- Seed travel_styles master table
INSERT INTO travel_styles (code, display_name, icon, sort_order, is_active) VALUES
    ('outdoor_activity', 'Outdoor Activity', '🏃', 1, true),
    ('social_activity', 'Social Activity', '👥', 2, true),
    ('cafe_dessert', 'Cafe & Dessert', '🍰', 3, true),
    ('bubble_tea', 'Bubble Tea', '🧋', 4, true),
    ('bakery_cake', 'Bakery / Cake', '🧁', 5, true),
    ('bingsu_ice_cream', 'Bingsu / Ice Cream', '🍧', 6, true),
    ('coffee', 'Coffee', '☕', 7, true),
    ('matcha', 'Matcha', '🍵', 8, true),
    ('pancakes', 'Pancakes', '🥞', 9, true),
    ('movie', 'Movie', '🎬', 10, true),
    ('karaoke', 'Karaoke', '🎤', 11, true),
    ('gaming', 'Gaming', '🎮', 12, true),
    ('board_game', 'Board Game', '🎲', 13, true),
    ('party_celebration', 'Party / Celebration', '🎉', 14, true),
    ('swimming', 'Swimming', '🏊', 15, true),
    ('skateboarding', 'Skateboarding', '🛹', 16, true)
ON CONFLICT (code) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Seed food_categories master table
INSERT INTO food_categories (code, display_name, icon, sort_order, is_active) VALUES
    ('thai_food', 'Thai Food', '🍛', 1, true),
    ('chinese_food', 'Chinese Food', '🥟', 2, true),
    ('japanese_food', 'Japanese Food', '🍣', 3, true),
    ('international_food', 'International Food', '🌍', 4, true),
    ('halal_food', 'Halal Food', '☪️', 5, true),
    ('buffet', 'Buffet', '🍽️', 6, true),
    ('bbq_grill', 'BBQ / Grill', '🔥', 7, true)
ON CONFLICT (code) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

