CREATE UNIQUE INDEX IF NOT EXISTS uniq_chapter_view
ON user_viewed_chapters (user_id, chapter_id)