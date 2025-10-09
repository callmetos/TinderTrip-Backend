-- Remove budget fields from events table
DROP INDEX IF EXISTS idx_events_currency;
DROP INDEX IF EXISTS idx_events_budget_max;
DROP INDEX IF EXISTS idx_events_budget_min;

ALTER TABLE events 
DROP COLUMN IF EXISTS currency,
DROP COLUMN IF EXISTS budget_max,
DROP COLUMN IF EXISTS budget_min;
