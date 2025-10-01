-- Create travel_preferences table
CREATE TABLE IF NOT EXISTS travel_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    travel_style VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(user_id, travel_style)
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_travel_preferences_user_id ON travel_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_travel_preferences_style ON travel_preferences(travel_style);

-- Insert default travel styles (no default insert needed for binary preferences)
