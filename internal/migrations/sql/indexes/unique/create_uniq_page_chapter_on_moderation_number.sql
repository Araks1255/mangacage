CREATE UNIQUE INDEX IF NOT EXISTS uniq_page_chapter_on_moderation_number
ON pages (chapter_on_moderation_id, number)