CREATE INDEX idx_messages_conversation_sent_id ON messages(conversation_id, sent_at DESC, id DESC);
