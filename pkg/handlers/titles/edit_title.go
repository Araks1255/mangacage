package titles

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"slices"
	"strconv"
	"sync"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
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

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingTitleID uint
	tx.Raw("SELECT id FROM titles WHERE id = ?", desiredTitleID).Scan(&existingTitleID)
	if existingTitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var userRoles []string
	tx.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON r.id = ur.role_id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	tx.Raw("SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)", existingTitleID, claims.ID).Scan(&doesUserTeamTranslatesDesiredTitle)
	if !doesUserTeamTranslatesDesiredTitle && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
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

		filter = bson.M{"title_id": existingTitleID}
		update := bson.M{"$set": bson.M{"cover": data}}
		opts := options.Update().SetUpsert(true)

		if _, err := h.TitlesOnModerationCovers.UpdateOne(context.TODO(), filter, update, opts); err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		var editedTitle models.TitleOnModeration

		editedTitle.ExistingID = sql.NullInt64{Int64: int64(existingTitleID), Valid: true}
		editedTitle.CreatorID = claims.ID // В creator_id будет записываться id того, кто отправил на модерацию (создатель записи на модерации, в целом логично)

		if len(values["name"]) != 0 {
			editedTitle.Name = sql.NullString{String: values["name"][0], Valid: true}
		}
		if len(values["description"]) != 0 {
			editedTitle.Description = values["description"][0]
		}

		if slices.Contains(userRoles, "moder") || slices.Contains(userRoles, "admin") { // Если юзер модератор или админ, то сохраняем его id как id последнего модератора (на всякий случай)
			editedTitle.ModeratorID = sql.NullInt64{Int64: int64(claims.ID), Valid: true}
		}

		if len(values["author"]) != 0 {
			var newAuthorID uint
			tx.Raw("SELECT id FROM authors WHERE name = ?", values["author"][0]).Scan(&newAuthorID)
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

		tx.Raw("SELECT id FROM titles_on_moderation WHERE existing_id = ?", editedTitle.ExistingID).Scan(&editedTitle.ID) // Если тайтл уже находится на модерации, то айди обращения записывается

		if result := tx.Save(&editedTitle); result.Error != nil {
			log.Println(result.Error)
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
			return
		}

		errChan <- nil
	}()

	wg.Wait()

	for i := 0; i < len(errChan); i++ {
		if err = <-errChan; err != nil {
			h.TitlesOnModerationCovers.DeleteOne(context.TODO(), filter)

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

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutTitleOnModeration(context.TODO(), &pb.TitleOnModeration{ID: uint64(existingTitleID), New: false}); err != nil {
		log.Println(err)
	}
}
