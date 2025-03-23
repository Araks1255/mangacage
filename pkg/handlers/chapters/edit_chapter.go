package chapters

import (
	"database/sql"
	"log"
	"slices"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) EditChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")
	volume := c.Param("volume")
	desiredChapter := c.Param("chapter")

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Volume      string `json:"volume"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Name == "" && requestBody.Description == "" && requestBody.Volume == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var chapterID, titleID uint
	row := tx.Raw(
		`SELECT chapters.id, titles.id FROM chapters
		INNER JOIN volumes ON volumes.id = chapters.volume_id
		INNER JOIN titles ON titles.id = volumes.title_id
		WHERE titles.name = ?
		AND volumes.name = ?
		AND chapters.name = ?`,
		title, volume, desiredChapter, // Опять же, подразумевается, что запрос будет отправляться на субдомен страницы с уже найденной и отображённой главой, так что приведения к нижнему регистру необязательны
	).Row()

	if err := row.Scan(&chapterID, &titleID); err != nil {
		log.Println(err)
	}

	if chapterID == 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	var userRoles []string
	tx.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		tx.Rollback()
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	tx.Raw("SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)", titleID, claims.ID).Scan(&doesUserTeamTranslatesDesiredTitle) // Сомнительная тема с поиском по имени, может быть поменяю
	if !doesUserTeamTranslatesDesiredTitle {
		tx.Rollback()
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит тайтл, в котором находится данная глава"})
		return
	}

	editedChapter := models.ChapterOnModeration{
		ExistingID:  sql.NullInt64{Int64: int64(chapterID), Valid: true},
		Name:        requestBody.Name, // Если имя менять не надо, то в запросе его не будет, и тут просто присвоится пустая строка, которая при записи превратится в NULL
		Description: requestBody.Description,
		CreatorID:   claims.ID,
	}

	if requestBody.Volume != "" {
		var volumeID uint
		tx.Raw("SELECT id FROM volumes WHERE title_id = ? AND lower(name) = lower(?)", titleID, requestBody.Volume).Scan(&volumeID)
		if volumeID == 0 {
			tx.Rollback()
			c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
			return
		}
		editedChapter.VolumeID = sql.NullInt64{Int64: int64(volumeID), Valid: true}
	}

	if slices.Contains(userRoles, "moder") || slices.Contains(userRoles, "admin") {
		editedChapter.ModeratorID = sql.NullInt64{Int64: int64(claims.ID), Valid: true}
	}

	tx.Raw("SELECT id FROM chapters_on_moderation WHERE existing_id = ?", editedChapter.ExistingID).Scan(&editedChapter.ID)

	if editedChapter.ID == 0 {
		if result := tx.Create(&editedChapter); result.Error != nil {
			log.Println(result.Error)
			tx.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
			return
		}
		tx.Commit()
		c.JSON(200, gin.H{"success": "изменения главы успешно отправлены на модерацию"})
		return
	}

	if result := tx.Exec(
		`UPDATE chapters_on_moderation SET
		created_at = ?,
		name = ?,
		description = ?,
		creator_id = ?,
		moderator_id = ?,
		volume_id = ?`,
		time.Now(), editedChapter.Name, editedChapter.Description,
		editedChapter.CreatorID, editedChapter.ModeratorID, editedChapter.VolumeID,
	); result.Error != nil {
		log.Println(result.Error)
		tx.Rollback()
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"error": "изменения главы успешно изменены"})
	// Уведомление
}
