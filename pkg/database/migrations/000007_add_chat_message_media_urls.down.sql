-- Remove image_url and file_url columns from chat_messages table
ALTER TABLE chat_messages
  DROP COLUMN IF EXISTS image_url,
  DROP COLUMN IF EXISTS file_url;

