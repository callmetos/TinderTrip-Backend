-- Create event_interests table (unified interests for events)
-- This replaces event_tags for interests, using the same interests table as users
CREATE TABLE event_interests (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    interest_id UUID NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id, interest_id)
);

-- Create indexes for better performance
CREATE INDEX idx_event_interests_event_id ON event_interests(event_id);
CREATE INDEX idx_event_interests_interest_id ON event_interests(interest_id);

