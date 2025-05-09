package titles

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitle(c *gin.Context) {
	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	var title models.TitleDTO

	h.DB.Raw(
		`SELECT
			t.id, t.created_at, t.name, t.description,
			a.name AS author, a.id AS author_id,
			MAX(teams.name) AS team, MAX(teams.id) AS team_id,
			ARRAY_AGG(DISTINCT g.name)::TEXT[] AS genres,
			COUNT(uvs.chapter_id) AS views
		FROM
			titles AS t
			INNER JOIN authors AS a ON a.id = t.author_id
			LEFT JOIN teams ON t.team_id = teams.id
			INNER JOIN title_genres AS tg ON tg.title_id = t.id
			INNER JOIN genres AS g ON tg.genre_id = g.id
			LEFT JOIN volumes AS v ON v.title_id = t.id
			LEFT JOIN chapters AS c ON c.volume_id = v.id
			LEFT JOIN user_viewed_chapters AS uvs ON uvs.chapter_id = c.id
		WHERE
			t.id = ?
		GROUP BY
			t.id, a.id`, // Тут используется MAX() для команды, потому-что все значения в неагрегатных функциях должны быть сгруппированы, а team_id у тайтла nullable столбец, и если использовать группировку по teams.id (даже с left join`ом), то все записи с team_id is null отбросятся. Поэтому teams.id записан в агрегатной функции MAX, чтобы не было необходимости его группировать
		desiredTitleID,
	).Scan(&title)

	if title.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	c.JSON(200, &title)
}
