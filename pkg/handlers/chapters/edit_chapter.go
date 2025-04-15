package chapters

import (
	"context"
	"database/sql"
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (h handler) EditChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id главы должен быть числом"})
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Volume      string `json:"volume"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Name == "" && requestBody.Description == "" && requestBody.Volume == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для редактирования главы"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	var titleID uint
	tx.Raw(
		`SELECT t.id FROM titles AS t
		INNER JOIN volumes AS v ON t.id = v.title_id
		INNER JOIN chapters AS c ON v.id = c.volume_id
		WHERE c.id = ?`, chapterID,
	).Scan(&titleID)

	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"}) // Пояснение в delete_chapter
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	tx.Raw("SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)", titleID, claims.ID).Scan(&doesUserTeamTranslatesDesiredTitle)
	if !doesUserTeamTranslatesDesiredTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит тайтл, в котором находится данная глава"})
		return
	}

	editedChapter := models.ChapterOnModeration{
		ExistingID:  sql.NullInt64{Int64: int64(chapterID), Valid: true},
		Name:        requestBody.Name,
		Description: requestBody.Description,
		CreatorID:   claims.ID,
	}

	if requestBody.Volume != "" {
		var volumeID uint
		tx.Raw("SELECT id FROM volumes WHERE title_id = ? AND lower(name) = lower(?)", titleID, requestBody.Volume).Scan(&volumeID)
		if volumeID == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
			return
		}
		editedChapter.VolumeID = sql.NullInt64{Int64: int64(volumeID), Valid: true}
	}

	if slices.Contains(userRoles, "moder") || slices.Contains(userRoles, "admin") {
		editedChapter.ModeratorID = sql.NullInt64{Int64: int64(claims.ID), Valid: true}
	}

	if result := tx.Save(&editedChapter); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"error": "изменения главы успешно отправлены на модерацию"})

	var chapterName string
	h.DB.Raw("SELECT name FROM chapters WHERE id = ?", chapterID).Scan(&chapterName)

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutChapterOnModeration(context.TODO(), &pb.ChapterOnModeration{Name: chapterName, New: false}); err != nil {
		log.Println(err)
	}
}
