ALTER TABLE teams
ADD CONSTRAINT fk_teams_moderator
FOREIGN KEY (moderator_id)
REFERENCES users (id) ON DELETE SET NULL;