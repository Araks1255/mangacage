CREATE OR REPLACE TRIGGER trg_title_rates_after_delete
AFTER DELETE ON title_rates
FOR EACH ROW
EXECUTE FUNCTION fn_title_rates_after_delete()