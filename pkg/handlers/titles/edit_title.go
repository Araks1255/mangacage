package titles

// import (
// 	"context"
// 	"errors"
// 	"io"
// 	"log"
// 	"slices"
// 	"sync"

// 	"github.com/Araks1255/mangacage/pkg/common/models"
// 	pb "github.com/Araks1255/mangacage_protos"
// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func (h handler) EditTitle(c *gin.Context) {
// 	claims := c.MustGet("claims").(*models.Claims)

// 	title := c.Param("title")

// 	var titleID, titleCreatorID uint

// 	row := h.DB.Raw("SELECT id, creator_id FROM titles WHERE lower(name) = lower(?) AND NOT on_moderation", title).Row()
// 	if row.Scan(&titleID, &titleCreatorID); titleID == 0 {
// 		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
// 		return
// 	}

// 	var userRoles []string
// 	h.DB.Raw(`SELECT roles.name FROM roles
// 	INNER JOIN user_roles ON roles.id = user_roles.role_id
// 	INNER JOIN users ON user_roles.user_id = users.id
// 	WHERE users.id = ?`, claims.ID).Scan(&userRoles)

// 	if titleCreatorID != claims.ID && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
// 		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь создателем этого тайтла"})
// 		return
// 	}

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		log.Println(err)
// 		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	const NUMBER_OF_GORUTINES = 7
// 	errChan := make(chan error, NUMBER_OF_GORUTINES)

// 	var wg sync.WaitGroup
// 	wg.Add(NUMBER_OF_GORUTINES)

// 	tx := h.DB.Begin()

// 	go func() {
// 		defer wg.Done()

// 		if len(form.Value["name"]) == 0 {
// 			errChan <- nil
// 			return
// 		}

// 		title := form.Value["name"][0]
// 		if result := tx.Exec("UPDATE titles SET name = ? WHERE id = ?", title, titleID); result.Error != nil {
// 			log.Println(result.Error)
// 			errChan <- result.Error
// 			return
// 		}

// 		errChan <- nil
// 	}()

// 	go func() {
// 		defer wg.Done()

// 		if len(form.Value["description"]) == 0 {
// 			errChan <- nil
// 			return
// 		}

// 		newDescription := form.Value["description"][0]
// 		if result := tx.Exec("UPDATE titles SET description = ? WHERE id = ?", newDescription, titleID); result.Error != nil {
// 			log.Println(result.Error)
// 			errChan <- result.Error
// 			return
// 		}

// 		errChan <- nil
// 	}()

// 	go func() {
// 		defer wg.Done()

// 		if len(form.Value["author"]) == 0 {
// 			errChan <- nil
// 			return
// 		}

// 		newAuthor := form.Value["author"][0]

// 		var newAuthorID uint
// 		h.DB.Raw("SELECT id FROM authors WHERE lower(name) = lower(?)", newAuthor).Scan(&newAuthorID)
// 		if newAuthorID == 0 {
// 			errChan <- errors.New("Новый автор не найден")
// 			return
// 		}

// 		if result := tx.Exec("UPDATE titles SET author_id = ?", newAuthorID); result.Error != nil {
// 			log.Println(result.Error)
// 			errChan <- result.Error
// 			return
// 		}

// 		errChan <- nil
// 	}()

// 	go func() {
// 		defer wg.Done()

// 		if len(form.Value["genres"]) == 0 {
// 			errChan <- nil
// 			return
// 		}

// 		newGenres := form.Value["genres"]

// 		if result := tx.Exec("DELETE FROM title_genres WHERE title_id = ?", titleID); result.Error != nil {
// 			log.Println(result.Error)
// 			errChan <- result.Error
// 			return
// 		}

// 		if err = AddGenresToTitle(titleID, newGenres, tx); err != nil {
// 			log.Println(err)
// 			errChan <- err
// 			return
// 		}

// 		errChan <- nil
// 	}()

// 	go func() {
// 		defer wg.Done()

// 		if len(form.File["cover"]) == 0 {
// 			errChan <- nil
// 			return
// 		}

// 		newCover, err := c.FormFile("cover")
// 		if err != nil {
// 			log.Println(err)
// 			errChan <- err
// 			return
// 		}

// 		file, err := newCover.Open()
// 		if err != nil {
// 			log.Println(err)
// 			errChan <- err
// 			return
// 		}
// 		defer file.Close()

// 		data, err := io.ReadAll(file)
// 		if err != nil {
// 			log.Println(err)
// 			errChan <- err
// 			return
// 		}

// 		filter := bson.M{"title_id": titleID}
// 		update := bson.M{"$set": bson.M{"cover": data}}

// 		if _, err := h.Collection.UpdateOne(context.TODO(), filter, update); err != nil {
// 			log.Println(err)
// 			errChan <- err
// 			return
// 		}

// 		errChan <- nil
// 	}()

// 	go func() {
// 		defer wg.Done()

// 		if !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
// 			errChan <- nil
// 			return
// 		}

// 		if result := tx.Exec("UPDATE titles SET moderator_id = ? WHERE id = ?", claims.ID, titleID); result.Error != nil {
// 			log.Println(result.Error)
// 			errChan <- result.Error
// 			return
// 		}

// 		errChan <- nil
// 	}()

// 	go func() {
// 		defer wg.Done()

// 		if result := tx.Exec("UPDATE titles SET on_moderation = true WHERE id = ?", titleID); result.Error != nil {
// 			log.Println(result.Error)
// 			errChan <- result.Error
// 			return
// 		}

// 		errChan <- nil
// 	}()

// 	wg.Wait()

// 	close(errChan)

// 	for i := 0; i < NUMBER_OF_GORUTINES; i++ {
// 		err = <-errChan
// 		if err != nil {
// 			tx.Rollback()
// 			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	tx.Commit()

// 	c.JSON(200, gin.H{"success": "тайтл успешно обновлён"})

// 	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer conn.Close()

// 	client := pb.NewNotificationsClient(conn)

// 	if _, err := client.NotifyAboutTitleOnModeration(context.Background(), &pb.TitleOnModeration{Name: title}); err != nil {
// 		log.Println(err)
// 	}
// }
