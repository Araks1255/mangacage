package teams

import (
	"context"
	"database/sql"
	"io"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) EditTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь владельцем команды перевода"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 && len(form.Value["description"]) == 0 && len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "запрос должен содержать хотя-бы один изменяемый параметр"})
		return
	}

	const NUMBER_OF_GORUTINES = 2 // Это вынесено наверх, чтобы не занимать время транзакции
	errChan := make(chan error, NUMBER_OF_GORUTINES)

	var wg sync.WaitGroup
	wg.Add(NUMBER_OF_GORUTINES)

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	var userTeamID uint // Поиск команды юзера на всякий случай производится в транзакции. Мало ли кто-то попробует одновременно выйти из команды и отредактировать её. И получится трындец
	tx.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID == 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}

	go func() {
		defer wg.Done()

		if len(form.File["cover"]) == 0 {
			errChan <- nil
			return
		}

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

		filter := bson.M{"team_id": userTeamID}
		update := bson.M{"$set": bson.M{"cover": data}}
		opts := options.Update().SetUpsert(true)

		if _, err := h.TeamsOnModerationCovers.UpdateOne(context.TODO(), filter, update, opts); err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		if len(form.Value["name"]) == 0 && len(form.Value["description"]) == 0 {
			errChan <- nil
			return
		}

		editedTeam := models.TeamOnModeration{
			ExistingID: sql.NullInt64{Int64: int64(userTeamID), Valid: true},
			CreatorID:  claims.ID,
		}

		if len(form.Value["name"]) != 0 {
			editedTeam.Name = sql.NullString{String: form.Value["name"][0], Valid: true}
		}
		if len(form.Value["description"]) != 0 {
			editedTeam.Description = form.Value["description"][0]
		}

		tx.Raw("SELECT id FROM teams_on_moderation WHERE existing_id = ?", userTeamID).Scan(&editedTeam.ID)

		editedTeam.CreatedAt = time.Now() // Попытка обновить уже существующее обращение считается новым обращением, поэтому CreatedAt меняется

		if result := tx.Save(&editedTeam); result.Error != nil {
			log.Println(result.Error)
			errChan <- result.Error
			return
		}

		errChan <- nil
	}()

	wg.Wait()

	for i := 0; i < NUMBER_OF_GORUTINES; i++ {
		if err := <-errChan; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return // любой return в хэндлере так или иначе означает tx.Rollback() (так-как этот метод вызван в defer)
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения команды успешно отправлены на модерацию"})
}
