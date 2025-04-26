package teams

import (
	"context"
	"database/sql"
	"io"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 || len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе нет имени команды или её обложки"})
		return
	}

	errChan := make(chan error)

	var teamCover struct {
		TeamOnModerationID uint   `bson:"team_on_moderation_id"`
		Cover              []byte `bson:"cover"`
	}

	go func() {
		file, err := form.File["cover"][0].Open()
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		teamCover.Cover = data

		errChan <- nil
	}()

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var userTeamID uint // Тут можно было бы получить роли юзера и его id команды одним запросом, но это довольно избыточно, + добавляет опасное место на сканировании ряда
	tx.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже состоите в другой команде"})
		return
	}

	var teamCreatedByUserID uint // Команда созданная пользователем
	tx.Raw("SELECT id FROM teams_on_moderation WHERE creator_id = ? LIMIT 1", claims.ID).Scan(&teamCreatedByUserID)
	if teamCreatedByUserID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "команда, созданная вами, уже ожидает модерации"})
		return
	}

	name := form.Value["name"][0]
	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	var existing struct {
		TeamID             uint
		TeamOnModerationID uint
	}

	tx.Raw(
		`SELECT t.id AS team_id, tom.id AS team_on_moderation_id FROM teams AS t
		RIGHT JOIN teams_on_moderation AS tom ON t.id = tom.existing_id
		WHERE lower(t.name) = lower(?)
		OR lower(tom.name) = lower(?)`,
		name, name,
	).Scan(&existing)

	if existing.TeamID != 0 || existing.TeamOnModerationID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "команда с таким названием уже существует"})
		return
	}

	newTeam := models.TeamOnModeration{
		Name:        sql.NullString{String: name, Valid: true},
		Description: description,
		CreatorID:   claims.ID,
	}

	if result := tx.Create(&newTeam); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if err = <-errChan; err != nil { // Канал небуферизированный, так что здесь хэндлер заблокируется, пока не выполнится горутина для обработки фото
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	teamCover.TeamOnModerationID = newTeam.ID

	if _, err := h.TeamsOnModerationCovers.InsertOne(context.TODO(), teamCover); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "команда успешно отправлена на модерацию"})
	// Уведомление
}
