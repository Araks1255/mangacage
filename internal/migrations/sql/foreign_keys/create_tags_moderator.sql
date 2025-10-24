ALTER TABLE tags
ADD CONSTRAINT fk_tags_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;