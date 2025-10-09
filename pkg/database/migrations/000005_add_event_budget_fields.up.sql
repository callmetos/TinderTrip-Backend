-- Add budget fields to events table
ALTER TABLE events 
ADD COLUMN budget_min INT CHECK (budget_min IS NULL OR budget_min >= 0),
ADD COLUMN budget_max INT CHECK (budget_max IS NULL OR budget_max >= 0),
ADD COLUMN currency VARCHAR(3) DEFAULT 'THB';

-- Add indexes for budget filtering
CREATE INDEX idx_events_budget_min ON events(budget_min) WHERE budget_min IS NOT NULL;
CREATE INDEX idx_events_budget_max ON events(budget_max) WHERE budget_max IS NOT NULL;
CREATE INDEX idx_events_currency ON events(currency) WHERE currency IS NOT NULL;
