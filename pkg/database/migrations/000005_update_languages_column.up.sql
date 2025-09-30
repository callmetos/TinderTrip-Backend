-- Update languages column from JSON to TEXT
ALTER TABLE user_profiles 
ALTER COLUMN languages TYPE TEXT USING languages::TEXT;
