-- Revert languages column from TEXT to JSON
ALTER TABLE user_profiles 
ALTER COLUMN languages TYPE JSON USING languages::JSON;
