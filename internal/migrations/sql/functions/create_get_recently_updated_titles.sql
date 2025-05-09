CREATE OR REPLACE FUNCTION get_recently_updated_titles(limit_count INT) -- Эта функция нужна для поиска тайтлов с сортировкой, где сверху будет тайтл, у которого прошло наименьшее время с момента выхода последней главы. Пришлось для этого создавать отдельную функцию, чтобы избежать дубликатов и иметь возможность задать четкий лимит
RETURNS TABLE ( -- Возвращается тайтл
    id BIGINT,
    name TEXT,
    author_id BIGINT,
    author TEXT,
    team_id BIGINT,
    team TEXT,
    genres TEXT[]
) AS $$
DECLARE
    title_count INT := 0; -- Счётчик найденных тайтлов
    chapter_row RECORD; -- !!!
BEGIN
    CREATE TEMP TABLE temp_titles ( -- Временная таблица для результата
        id BIGINT PRIMARY KEY,
        name TEXT,
        author_id BIGINT,
        author TEXT,
        team_id BIGINT,
        team TEXT,
        genres TEXT[]
    ) ON COMMIT DROP; -- С удалением при коммите транзакции

    FOR chapter_row IN -- цикл для перебора всех глав в порядке выхода (сначала новые)
        SELECT c.id AS chapter_id, v.title_id -- v.title_id - нужный нам id тайтла
        FROM chapters AS c
        INNER JOIN volumes AS v ON v.id = c.volume_id
        ORDER BY c.created_at DESC
    LOOP
        EXIT WHEN title_count >= limit_count; -- Выходим, если лимит достигнут

        IF EXISTS (SELECT 1 FROM temp_titles AS tt WHERE tt.id = chapter_row.title_id) THEN -- Если тайтл уже есть в таблице 
            CONTINUE; -- Пропускаем итерацию
        END IF;

        INSERT INTO temp_titles (id, name, author_id, author, team_id, team, genres) -- Если его всё-таки нет в таблице результатов
        SELECT t.id, t.name, a.id, a.name, tm.id, tm.name, -- Получаем о нём все данные
            (
                SELECT ARRAY_AGG(g.name)
                FROM title_genres tg
                JOIN genres g ON g.id = tg.genre_id
                WHERE tg.title_id = t.id
            )
        FROM titles t
        JOIN authors a ON a.id = t.author_id
        LEFT JOIN teams tm ON tm.id = t.team_id
        WHERE t.id = chapter_row.title_id;

        title_count := title_count + 1; -- Увеличиваем счётчик
    END LOOP;

    RETURN QUERY SELECT * FROM temp_titles;
END;
$$ LANGUAGE plpgsql;
