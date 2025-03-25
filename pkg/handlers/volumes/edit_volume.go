package volumes

import (
	"context"
	"database/sql"
	"log"
	"slices"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (h handler) EditVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")
	volume := c.Param("volume")

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Name == "" && requestBody.Description == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var titleID, volumeID uint
	row := tx.Raw(`SELECT titles.id, volumes.id FROM volumes
			INNER JOIN titles ON titles.id = volumes.title_id
			WHERE titles.name = ? AND volumes.name = ?`,
		title, volume).Row()

	if err := row.Scan(&titleID, &volumeID); err != nil {
		log.Println(err)
	}

	if volumeID == 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	var userRoles []string
	tx.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		tx.Rollback()
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	var doesUserTeamTranslatesThisTitle bool
	h.DB.Raw("SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)", titleID, claims.ID).Scan(&doesUserTeamTranslatesThisTitle)
	if !doesUserTeamTranslatesThisTitle {
		tx.Rollback()
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
		return
	}

	editedVolume := models.VolumeOnModeration{
		ExistingID:  sql.NullInt64{Int64: int64(volumeID), Valid: true},
		Name:        requestBody.Name,
		Description: requestBody.Description,
		TitleID:     sql.NullInt64{Int64: int64(titleID), Valid: true},
		CreatorID:   claims.ID,
	}

	if slices.Contains(userRoles, "moder") || slices.Contains(userRoles, "admin") {
		editedVolume.ModeratorID = sql.NullInt64{Int64: int64(claims.ID), Valid: true}
	}

	tx.Raw("SELECT id FROM volumes_on_moderation WHERE existing_id = ?", editedVolume.ExistingID).Scan(&editedVolume.ID)

	if editedVolume.ID == 0 { // Тут пришлось делать так, потому-что метод Save устраивал какую-то вакханалию с временем создания при обновлении записи
		if result := tx.Create(&editedVolume); result.Error != nil {
			log.Println(result.Error)
			tx.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
			return
		}
		tx.Commit()
		c.JSON(200, gin.H{"success": "изменения тома успешно отправлены на модерацию"})
		return
	}

	if result := tx.Exec( // А тут теперь вручную задаются изменения. Текущее время записывается не в updated_at, а в created_at, потому-что это created_at - время отправки на модерацию, а тут как-бы, идёт повторная отправка на модерацию
		`UPDATE volumes_on_moderation SET
		created_at = ?,
		name = ?,
		description = ?,
		creator_id = ?,
		moderator_id = ?`,
		time.Now(), editedVolume.Name, editedVolume.Description, editedVolume.CreatorID, editedVolume.ModeratorID,
	); result.Error != nil {
		log.Println(result.Error)
		tx.Rollback()
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "изменения тома успешно изменены"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutVolumeOnModeration(context.TODO(), &pb.VolumeOnModeration{Name: volume, New: false}); err != nil {
		log.Println(err)
	}
}
