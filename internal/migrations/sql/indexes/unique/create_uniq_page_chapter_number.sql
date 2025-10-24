CREATE UNIQUE INDEX IF NOT EXISTS uniq_page_chapter_number
ON pages (chapter_id, number)