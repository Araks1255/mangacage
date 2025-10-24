ALTER TABLE teams_on_moderation
ADD CONSTRAINT fk_teams_on_moderation_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;