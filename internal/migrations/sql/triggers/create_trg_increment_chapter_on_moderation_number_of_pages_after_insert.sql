CREATE OR REPLACE TRIGGER trg_increment_chapter_on_moderation_number_of_pages_after_insert
AFTER INSERT ON pages
FOR EACH ROW
WHEN (NEW.chapter_on_moderation_id IS NOT NULL)
EXECUTE FUNCTION increment_chapter_on_moderation_number_of_pages()
-- На данный момент такой сценарий не нужен, ведь страницы вставляются без chapter_on_moderation_id, и он добавляется к ним позже,
-- но триггер с when на всякий случай есть