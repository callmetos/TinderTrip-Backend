-- Remove seed data from master tables (optional - you might want to keep the data)
-- Uncomment if you want to remove seed data on rollback

-- DELETE FROM travel_styles WHERE code IN (
--     'outdoor_activity', 'social_activity', 'cafe_dessert', 'bubble_tea',
--     'bakery_cake', 'bingsu_ice_cream', 'coffee', 'matcha', 'pancakes',
--     'movie', 'karaoke', 'gaming', 'board_game', 'party_celebration',
--     'swimming', 'skateboarding'
-- );

-- DELETE FROM food_categories WHERE code IN (
--     'thai_food', 'chinese_food', 'japanese_food', 'international_food',
--     'halal_food', 'buffet', 'bbq_grill'
-- );

-- Note: We don't actually delete the data on rollback to preserve master data

