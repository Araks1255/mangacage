package chapters

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (h handler) CreateChapter(c *gin.Context) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	pathToDirectory := viper.Get("PATH_TO_CHAPTERS_DIRECTORY").(string)

	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles ON roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if IsUserTeamOwner := slices.Contains(userRoles, "team_leader"); !IsUserTeamOwner {
		c.AbortWithStatusJSON(403, gin.H{"error": "Добавлять главы может только лидер команды"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	title := strings.ToLower(c.Param("title"))

	name := strings.ToLower(form.Value["name"][0])
	description := strings.ToLower(form.Value["description"][0])

	numberOfPages, err := strconv.Atoi(form.Value["numberOfPages"][0])
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var desiredTitleID uint
	h.DB.Raw("SELECT id FROM titles WHERE name = ?", title).Scan(&desiredTitleID)
	if desiredTitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Тайтл, в котором вы хотите выложить главу не найден"})
		return
	}

	var IsUserTeamTranslatesThisTitle bool
	h.DB.Raw("SELECT CAST("+
		"CASE WHEN ? = "+
		"(SELECT titles.id FROM titles INNER JOIN teams ON titles.team_id = teams.id INNER JOIN users ON teams.id = users.team_id WHERE users.id = ?) "+
		"THEN true ELSE false END AS BOOLEAN)", desiredTitleID, claims.ID).Scan(&IsUserTeamTranslatesThisTitle)

	if !IsUserTeamTranslatesThisTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "Ваша команда не переводит данный тайтл"})
		return
	}

	const NUMBER_OF_GORUTINES int = 2
	errChan := make(chan error, NUMBER_OF_GORUTINES)

	pathToChapter := fmt.Sprintf("%s/%s/%s", pathToDirectory, title, name)

	if err := os.MkdirAll(pathToChapter, 0755); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	go func() {
		for i := 0; i < numberOfPages; i++ {
			file, err := c.FormFile(strconv.Itoa(i))
			if err != nil {
				log.Println(err)
				errChan <- err
				return
			}

			path := fmt.Sprintf("%s/%d.jpg", pathToChapter, i)
			page, err := os.Create(path)
			if err != nil {
				log.Println(err)
				errChan <- err
				return
			}

			if err := c.SaveUploadedFile(file, path); err != nil {
				log.Println(err)
				errChan <- err
				return
			}

			page.Close()
		}

		errChan <- nil
	}()

	transaction := h.DB.Begin()
	go func() {
		chapter := models.Chapter{
			Name:          name,
			Description:   description,
			Path:          pathToChapter,
			NumberOfPages: numberOfPages,
			TitleID:       desiredTitleID,
		}

		if IsUserAdmin := slices.Contains(userRoles, "admin"); IsUserAdmin {
			chapter.OnModeration = false
		} else {
			chapter.OnModeration = true
		}

		if result := transaction.Create(&chapter); result.Error != nil {
			log.Println(err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	for i := 0; i < NUMBER_OF_GORUTINES; i++ {
		if <-errChan != nil {
			transaction.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": "Произошла ошибка при создании главы"})
			os.RemoveAll(pathToChapter)
			return
		}
	}

	transaction.Commit()

	c.JSON(201, gin.H{"success": "Глава успешно создана"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutNewChapterOnModeration(context.Background(), &pb.ChapterOnModeration{TitleName: title, ChapterName: name}); err != nil {
		log.Println(err)
	}
}
