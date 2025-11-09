-- Add unique constraint to display_name column
-- First, handle any duplicate display_names by appending a suffix
DO $$
DECLARE
    dup_rec RECORD;
    user_rec RECORD;
    counter INTEGER;
    new_display_name TEXT;
BEGIN
    -- Find and fix duplicate display_names
    FOR dup_rec IN 
        SELECT display_name, COUNT(*) as cnt
        FROM users
        WHERE display_name IS NOT NULL
          AND deleted_at IS NULL
        GROUP BY display_name
        HAVING COUNT(*) > 1
    LOOP
        counter := 1;
        FOR user_rec IN 
            SELECT id, display_name
            FROM users
            WHERE display_name = dup_rec.display_name
              AND deleted_at IS NULL
            ORDER BY created_at
        LOOP
            IF counter > 1 THEN
                -- Append suffix to make it unique
                new_display_name := user_rec.display_name || '_' || counter;
                UPDATE users
                SET display_name = new_display_name
                WHERE id = user_rec.id;
            END IF;
            counter := counter + 1;
        END LOOP;
    END LOOP;
END $$;

-- Create unique index on display_name (excluding NULL and soft-deleted records)
-- Drop index if exists first to avoid errors
DROP INDEX IF EXISTS ux_users_display_name;

CREATE UNIQUE INDEX ux_users_display_name 
ON users(display_name) 
WHERE display_name IS NOT NULL AND deleted_at IS NULL;

-- Add comment
COMMENT ON INDEX ux_users_display_name IS 'Ensures display_name is unique across active users';

