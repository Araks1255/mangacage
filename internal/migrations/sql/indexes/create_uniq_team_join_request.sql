CREATE UNIQUE INDEX IF NOT EXISTS uniq_team_join_request
ON team_join_requests (candidate_id, team_id)