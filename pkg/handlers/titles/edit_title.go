package titles

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	desiredTitle := c.Param("title")

	tx := h.DB.Begin()
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}
	defer tx.Rollback()

	var titleID uint
	tx.Raw("SELECT id FROM titles WHERE name = ?", desiredTitle).Scan(&titleID) // Здесь нет приведения к нижнему регситру, так-как по задумке запрос будет отправляться с уже отрисованной страницы просмотра тайтла (ну и ещё - таблица не гарантирует уникальности при приведении к нижнему регистру, это обеспечивается только при создании тайтла)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var userRoles []string
	tx.Raw("SELECT roles.name FROM roles INNER JOIN user_roles ON roles.id = user_roles.role_id WHERE user_roles.user_id = ?", claims.ID).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	tx.Raw("SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)", titleID, claims.ID).Scan(&doesUserTeamTranslatesDesiredTitle)
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

	const NUMBER_OF_GORUTINES = 2
	errChan := make(chan error, NUMBER_OF_GORUTINES)

	var wg sync.WaitGroup
	wg.Add(NUMBER_OF_GORUTINES)

	filter := bson.M{}
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

		filter = bson.M{"title_id": titleID}
		update := bson.M{"$set": bson.M{"cover": data}}
		opts := options.Update().SetUpsert(true)

		if _, err := h.Collection.UpdateOne(context.TODO(), filter, update, opts); err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		var editedTitle models.TitleOnModeration

		editedTitle.ExistingID = sql.NullInt64{Int64: int64(titleID), Valid: true}
		editedTitle.CreatorID = claims.ID // В creator_id будет записываться id того, кто отправил на модерацию (создатель записи на модерации, в целом логично)

		if len(values["name"]) != 0 {
			editedTitle.Name = values["name"][0]
		}
		if len(values["description"]) != 0 {
			editedTitle.Description = values["description"][0]
		}

		if slices.Contains(userRoles, "moder") || slices.Contains(userRoles, "admin") { // Если юзер модератор или админ, то сохраняем его id как id последнего модератора (на всякий случай)
			editedTitle.ModeratorID = sql.NullInt64{Int64: int64(claims.ID), Valid: true}
		}

		if len(values["author"]) != 0 {
			var newAuthorID uint
			tx.Raw("SELECT id FROM authors WHERE lower(name) = lower(?)", values["author"][0]).Scan(&newAuthorID)
			if newAuthorID == 0 {
				errChan <- errors.New("автор не найден")
				return
			}
			editedTitle.AuthorID = sql.NullInt64{Int64: int64(newAuthorID), Valid: true}
		}

		if len(values["genres"]) != 0 {
			var genresIDs pq.StringArray
			tx.Raw("SELECT genres.id FROM genres JOIN UNNEST(?::TEXT[]) AS genre_name ON genres.name = genre_name", pq.StringArray(values["genres"])).Scan(&genresIDs)

			if len(values["genres"]) != len(genresIDs) {
				errChan <- errors.New("указан несуществующий жанр")
				return
			}

			editedTitle.Genres = values["genres"]
		}

		tx.Raw("SELECT id FROM titles_on_moderation WHERE existing_id = ?", editedTitle.ExistingID).Scan(&editedTitle.ID) // Если тайтл уже находится на модерации, то айди обращения записывается, чтобы метод Save обновил обращение, а не пытался создать заново

		if editedTitle.ID == 0 { // Пояснение этой свистопляски в edit_volume
			if result := tx.Create(&editedTitle); result.Error != nil {
				log.Println(result.Error)
				errChan <- result.Error
				return
			}
		} else {
			if result := tx.Exec(
				`UPDATE titles_on_moderation SET
				created_at = ?,
				name = ?,
				description = ?,
				creator_id = ?,
				moderator_id = ?,
				author_id = ?,
				genres = ?
				WHERE existing_id = ?`,
				time.Now(), editedTitle.Name, editedTitle.Description, editedTitle.CreatorID,
				editedTitle.ModeratorID, editedTitle.AuthorID, editedTitle.Genres, editedTitle.ExistingID,
			); result.Error != nil {
				log.Println(result.Error)
				errChan <- result.Error
				return
			}
		}

		errChan <- nil
	}()

	wg.Wait()

	for i := 0; i < NUMBER_OF_GORUTINES; i++ {
		if err = <-errChan; err != nil {
			_, _ = h.Collection.DeleteOne(context.TODO(), filter)

			if err.Error() == "автор не найден" || err.Error() == "указан несуществующий жанр" {
				c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
				return
			}

			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения тайтла успешно отправлены на модерацию"})
	// уведомление
}
