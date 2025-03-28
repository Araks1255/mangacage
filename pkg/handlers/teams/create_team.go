package teams

import (
	"context"
	"io"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if slices.Contains(userRoles, "team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже являетесь владельцем другой команды"})
		return
	}

	var userTeamID uint // Тут можно было бы получить роли юзера и его id команды одним запросом, но это довольно избыточно, + добавляет опасное место на сканировании ряда
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже состоите в другой команде"})
		return
	}

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

	var existing struct { // Там со сканированием ряда пришлось бы создавать две переменные sql.NullInt64, а это памяти много жрёт. Поэтому так
		TeamID             uint
		TeamOnModerationID uint
	}

	h.DB.Raw(
		`SELECT teams.id AS team_id, teams_on_moderation.id AS team_on_moderation_id FROM teams
		RIGHT JOIN teams_on_moderation ON teams.id = teams_on_moderation.existing_id
		WHERE lower(teams.name) = lower(?)
		OR lower(teams_on_moderation.name) = lower(?)`,
		form.Value["name"][0], form.Value["name"][0], // Тут в запросе приведение к нижнему регистру, чтобы нельзя было создать 2 команды с одинаковым именем но разным регистром
	).Scan(&existing)

	if existing.TeamID != 0 || existing.TeamOnModerationID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "команда с таким названием уже существует"})
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

	name := form.Value["name"][0]
	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	newTeam := models.TeamOnModeration{
		Name:        name,
		Description: description,
		CreatorID:   claims.ID,
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

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
