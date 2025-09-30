-- Create pref_availability table
CREATE TABLE pref_availability (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mon BOOLEAN DEFAULT false,
    tue BOOLEAN DEFAULT false,
    wed BOOLEAN DEFAULT false,
    thu BOOLEAN DEFAULT false,
    fri BOOLEAN DEFAULT false,
    sat BOOLEAN DEFAULT false,
    sun BOOLEAN DEFAULT false,
    all_day BOOLEAN DEFAULT false,
    morning BOOLEAN DEFAULT false,
    afternoon BOOLEAN DEFAULT false,
    time_range TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create pref_budget table
CREATE TABLE pref_budget (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_min INTEGER DEFAULT 0,
    meal_max INTEGER DEFAULT 0,
    daytrip_min INTEGER DEFAULT 0,
    daytrip_max INTEGER DEFAULT 0,
    overnight_min INTEGER DEFAULT 0,
    overnight_max INTEGER DEFAULT 0,
    unlimited BOOLEAN DEFAULT false,
    currency TEXT DEFAULT 'THB',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_pref_availability_user_id ON pref_availability(user_id);
CREATE INDEX idx_pref_budget_user_id ON pref_budget(user_id);

-- Create triggers
CREATE TRIGGER update_pref_availability_updated_at 
    BEFORE UPDATE ON pref_availability 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pref_budget_updated_at 
    BEFORE UPDATE ON pref_budget 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
