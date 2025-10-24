ALTER TABLE titles_translate_requests
ADD CONSTRAINT fk_titles_translate_requests_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL