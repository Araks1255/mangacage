ALTER TABLE titles_on_moderation
ADD CONSTRAINT fk_titles_on_moderation_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;