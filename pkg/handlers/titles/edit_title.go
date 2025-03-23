package titles

import (
	"context"
	"database/sql"
	"io"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	desiredTitle := c.Param("title")

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE name = ?", desiredTitle).Scan(&titleID) // Здесь нет приведения к нижнему регситру, так-как по задумке запрос будет отправляться с уже отрисованной страницы просмотра тайтла (ну и ещё - таблица не гарантирует уникальности при приведении к нижнему регистру, это обеспечивается только при создании тайтла)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var userRoles []string
	h.DB.Raw("SELECT roles.name FROM roles INNER JOIN user_roles ON roles.id = user_roles.role_id WHERE user_roles.user_id = ?", claims.ID).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	h.DB.Raw("SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)", titleID, claims.ID).Scan(&doesUserTeamTranslatesDesiredTitle)
	if !doesUserTeamTranslatesDesiredTitle && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	values := form.Value

	if len(values["name"]) == 0 && len(values["description"]) == 0 && len(values["author"]) == 0 && len(values["genres"]) == 0 && len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	var newTitleCover struct {
		TitleID uint   `bson:"title_id"`
		Cover   []byte `bson:"cover"`
	}

	errChan := make(chan error)

	go func() { // Новая обложка обрабатывается в отдельной горутине
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
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		newTitleCover.Cover = data

		errChan <- nil
	}()

	var editedTitle models.TitleOnModeration
	editedTitle.ExistingID = sql.NullInt64{Int64: int64(titleID), Valid: true}
	editedTitle.CreatorID = claims.ID // В creator_id будет записываться id того, кто отправил на модерацию (создатель записи на модерации, в целом логично)

	tx := h.DB.Begin() // Транзакция (так рано) нужна для гарантии того, что данные не изменятся по ходу выполнения хэндлера

	if len(values["name"]) != 0 {
		editedTitle.Name = values["name"][0]
	}
	if len(values["description"]) != 0 {
		editedTitle.Description = values["description"][0]
	}

	if len(values["author"]) != 0 {
		var newAuthorID uint
		tx.Raw("SELECT id FROM authors WHERE lower(name) = lower(?)", values["author"][0]).Scan(&newAuthorID)
		if newAuthorID == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "автор не найден"})
			return
		}
		editedTitle.AuthorID = sql.NullInt64{Int64: int64(newAuthorID), Valid: true}
	}

	if len(values["genres"]) != 0 {
		var genresIDs pq.StringArray
		tx.Raw("SELECT genres.id FROM genres JOIN UNNEST(?::TEXT) AS genre_name ON genres.name = genre_name", pq.StringArray(values["genres"])).Scan(&genresIDs)

		if len(values["genres"]) != len(genresIDs) {
			c.AbortWithStatusJSON(404, gin.H{"error": "указан несуществующий жанр"})
			return
		}

		editedTitle.Genres = values["genres"]
	}

	if err = <-errChan; err != nil { // Канал небуферизированный, поэтому горутина хэндлера заблокируется здесь при попытке считать пустое значение, а разблокируется только тогда, когда значение появится (то есть, после завершения выполения горутины для обработки обложкм)
		tx.Commit() // Пока что было только чтение
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	close(errChan)

	h.DB.Raw("SELECT id FROM titles_on_moderation WHERE existing_id = ?", editedTitle.ExistingID).Scan(&editedTitle.ID) // Если тайтл уже находится на модерации, то айди обращения записывается, чтобы метод Save обновил обращение, а не пытался создать заново

	if result := tx.Save(&editedTitle); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if len(form.File["cover"]) != 0 {
		newTitleCover.TitleID = editedTitle.ID

		filter := bson.M{"title_id": editedTitle.ID}
		update := bson.M{"$set": bson.M{"cover": newTitleCover.Cover}}

		if err := h.Collection.FindOneAndUpdate(context.TODO(), filter, update); err != nil {
			if _, err := h.Collection.InsertOne(context.TODO(), newTitleCover); err != nil {
				tx.Rollback()
				log.Println(err)
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения тайтла успешно отправлены на модерацию"})
	// уведомление
}
