ALTER TABLE users
ADD CONSTRAINT fk_users_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;