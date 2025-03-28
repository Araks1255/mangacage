package teams

import (
	"context"
	"io"
	"log"
	"slices"
	"sync"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) EditTeam(c *gin.Context) { // Это старая функция, сейчас буду переделывать, как с транзакциями закончу
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "редактировать команду может только её лидер"})
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

	var teamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&teamID)
	if teamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "произошла ошибка. команда не найдена"}) // логически это невозможно, но мало ли
		return
	}

	const NUMBER_OF_GORUTINES = 3
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

	go func() {
		defer wg.Done()

		if len(form.Value["name"]) == 0 {
			errChan <- nil
			return
		}

		newName := form.Value["name"][0]

		if result := tx.Exec("UPDATE teams SET name = ? WHERE id = ?", newName, teamID); result.Error != nil {
			log.Println(result.Error)
			errChan <- result.Error
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		if len(form.Value["description"]) == 0 {
			errChan <- nil
			return
		}

		newDescription := form.Value["description"][0]

		if result := tx.Exec("UPDATE teams SET description = ? WHERE id = ?", newDescription, teamID); result.Error != nil {
			log.Println(result.Error)
			errChan <- result.Error
			return
		}

		errChan <- nil
	}()

	go func() {
		defer wg.Done()

		if len(form.File["cover"]) == 0 {
			errChan <- nil
			return
		}

		newCover := form.File["cover"][0]

		file, err := newCover.Open()
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

		filter := bson.M{"team_id": teamID}
		update := bson.M{"$set": bson.M{"cover": data}}

		if _, err = h.Collection.UpdateOne(context.TODO(), filter, update); err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	wg.Wait()
	close(errChan)

	for i := 0; i < NUMBER_OF_GORUTINES; i++ {
		err = <-errChan
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "команда успешно обновлена"})
}
