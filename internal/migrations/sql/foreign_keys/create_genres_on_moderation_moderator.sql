ALTER TABLE genres_on_moderation
ADD CONSTRAINT fk_genres_on_moderation_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;