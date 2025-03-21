package users

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	const NUMBER_OF_GORUTINES = 4
	errChan := make(chan error, NUMBER_OF_GORUTINES)

	var wg sync.WaitGroup
	wg.Add(NUMBER_OF_GORUTINES)

	tx := h.DB.Begin()

	go func() {
		defer wg.Done()

		if len(form.Value["userName"]) == 0 {
			errChan <- nil
			return
		}

		newUserName := form.Value["userName"][0]

		if result := tx.Exec("UPDATE users SET user_name = ? WHERE id = ?", newUserName, claims.ID); result.Error != nil {
			log.Println(result.Error)
			errChan <- result.Error
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		if len(form.Value["aboutYourself"]) == 0 {
			errChan <- nil
			return
		}

		newAboutYourself := form.Value["aboutYourself"][0]

		if result := tx.Exec("UPDATE users SET about_yourself = ? WHERE id = ?", newAboutYourself, claims.ID); result.Error != nil {
			log.Println(result.Error)
			errChan <- result.Error
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		if len(form.File["profilePicture"]) == 0 {
			errChan <- nil
			return
		}

		newProfilePicture, err := c.FormFile("profilePicture")
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		file, err := newProfilePicture.Open()
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

		filter := bson.M{"user_id": claims.ID}
		update := bson.M{"$set": bson.M{"profile_picture": data}}

		if _, err = h.Collection.UpdateOne(context.TODO(), filter, update); err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		if result := tx.Exec("UPDATE users SET on_moderation = true WHERE id = ?", claims.ID); result.Error != nil {
			log.Println(result.Error)
			errChan <- result.Error
			return
		}

		errChan <- nil
	}()

	wg.Wait()

	for i := 0; i < NUMBER_OF_GORUTINES; i++ {
		err = <-errChan
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "профиль успешно обновлён"})

	var userName string
	h.DB.Raw("SELECT user_name FROM users WHERE id = ?", claims.ID).Scan(&userName)
	if userName == "" {
		return
	}

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutUser(context.Background(), &pb.User{Name: userName}); err != nil {
		log.Println(err)
	}
}
