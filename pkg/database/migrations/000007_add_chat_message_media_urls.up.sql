-- Add image_url and file_url columns to chat_messages table
ALTER TABLE chat_messages
  ADD COLUMN IF NOT EXISTS image_url TEXT,
  ADD COLUMN IF NOT EXISTS file_url TEXT;

