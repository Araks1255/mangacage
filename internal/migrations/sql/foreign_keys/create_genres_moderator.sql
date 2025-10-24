ALTER TABLE genres
ADD CONSTRAINT fk_genres_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;