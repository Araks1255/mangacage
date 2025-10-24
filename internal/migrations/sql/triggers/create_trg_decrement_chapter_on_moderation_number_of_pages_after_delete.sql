CREATE OR REPLACE TRIGGER trg_decrement_chapter_on_moderation_number_of_pages_after_delete
AFTER DELETE ON pages
FOR EACH ROW
EXECUTE FUNCTION increment_chapter_on_moderation_number_of_pages()
-- Бизнес-логикой не подразумевается удаление страниц, но чисто на всякий случай, для консистентности данных пусть будет