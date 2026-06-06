ALTER TABLE conversation_participants
ADD COLUMN last_read_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL;
