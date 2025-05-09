package volumes

import (
	"database/sql"
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	var requestBody struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для создания тома"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		DoesTitleExist                 bool
		DoesUserTeamTranslateTitle     bool
		DoesVolumeWithTheSameNameExist bool
	}

	if err := tx.Raw(
		`SELECT
			EXISTS(SELECT 1 FROM titles WHERE id = ?) AS does_title_exist,
			(SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)) AS does_user_team_translate_title,
			EXISTS(SELECT 1 FROM volumes WHERE lower(name) = lower(?)) AS does_volume_with_the_same_name_exist`,
		titleID, titleID, claims.ID, requestBody.Name,
	).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !check.DoesTitleExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}
	if !check.DoesUserTeamTranslateTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит этот тайтл"})
		return
	}
	if check.DoesVolumeWithTheSameNameExist {
		c.AbortWithStatusJSON(409, gin.H{"error": "том с таким названием уже существует"})
		return
	}

	volume := models.VolumeOnModeration{
		Name:        sql.NullString{String: requestBody.Name, Valid: true},
		Description: requestBody.Description,
		TitleID:     uint(titleID),
		CreatorID:   claims.ID,
	}

	err = tx.Create(&volume).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniVolumeTitle) {
			c.AbortWithStatusJSON(409, gin.H{"error": "том с таким названием уже ожидает модерации в этом тайтле"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "том успешно отправлен на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutVolumeOnModeration(c.Request.Context(), &pb.VolumeOnModeration{ID: uint64(volume.ID), New: true}); err != nil {
		log.Println(err)
	}
}
