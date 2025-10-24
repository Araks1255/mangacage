ALTER TABLE titles
ADD CONSTRAINT fk_titles_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;