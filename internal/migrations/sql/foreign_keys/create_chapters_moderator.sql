ALTER TABLE chapters
ADD CONSTRAINT fk_chapters_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;