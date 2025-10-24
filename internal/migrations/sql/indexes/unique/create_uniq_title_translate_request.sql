CREATE UNIQUE INDEX IF NOT EXISTS uniq_title_translate_request
ON titles_translate_requests (title_id, team_id)