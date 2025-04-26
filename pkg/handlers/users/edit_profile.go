package users

import (
	"context"
	"database/sql"
	"io"
	"log"
	"sync"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (h handler) EditProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 && len(form.Value["aboutYourself"]) == 0 && len(form.File["profilePicture"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	const NUMBER_OF_GORUTINES = 2
	errChan := make(chan error, NUMBER_OF_GORUTINES)

	var wg sync.WaitGroup
	wg.Add(NUMBER_OF_GORUTINES)

	tx := h.DB.Begin()
	utils.RollbackOnPanic(tx)

	var filter bson.M
	go func() {
		defer wg.Done()

		if len(form.File["profilePicture"]) == 0 {
			errChan <- nil
			return
		}

		file, err := form.File["profilePicture"][0].Open()
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

		filter = bson.M{"user_id": claims.ID}
		update := bson.M{"$set": bson.M{"profile_picture": data}}
		opts := options.Update().SetUpsert(true)

		if _, err := h.Collection.UpdateOne(context.TODO(), filter, update, opts); err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	var editedUser models.UserOnModeration // Это потом поменяю, не дошёл ещё, а под изменения grpc адаптировать надо
	go func() {
		defer wg.Done()

		if len(form.Value["userName"]) == 0 && len(form.Value["aboutYourself"]) == 0 {
			errChan <- nil
			return
		}

		editedUser.ExistingID = sql.NullInt64{Int64: int64(claims.ID), Valid: true}

		if len(form.Value["userName"]) != 0 {
			editedUser.UserName = sql.NullString{String: form.Value["userName"][0], Valid: true}
		}
		if len(form.Value["aboutYourself"]) != 0 {
			editedUser.AboutYourself = form.Value["aboutYourself"][0]
		}

		tx.Raw("SELECT id FROM users_on_moderation WHERE existing_id = ?", claims.ID).Scan(&editedUser.ID)

		if result := tx.Save(&editedUser); result.Error != nil {
			log.Println(result.Error)
			errChan <- result.Error
			return
		}

		errChan <- nil
	}()

	wg.Wait()

	for i := 0; i < NUMBER_OF_GORUTINES; i++ {
		if err = <-errChan; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			tx.Rollback()
			_, _ = h.Collection.DeleteOne(context.TODO(), filter)
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения профиля успешно отправлены на модерацию"})

	var userName string
	h.DB.Raw("SELECT user_name FROM users WHERE id = ?", claims.ID).Scan(&userName)

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutUserOnModeration(context.Background(), &pb.User{ID: uint64(editedUser.ExistingID.Int64), New: false}); err != nil {
		log.Println(err)
	}
}
