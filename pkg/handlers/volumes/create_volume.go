package volumes

import (
	"context"
	"io"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
	}

	cover, err := c.FormFile("cover")
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
		return
	}

	title := strings.ToLower(c.Param("title"))

	name := strings.ToLower(form.Value["name"][0])
	description := strings.ToLower(form.Value["description"][0])

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE name = ?", title).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Тайтл не найден"})
		return
	}

	var existingVolumeID uint
	h.DB.Raw("SELECT id FROM volumes WHERE name = ? AND title_id = ?", name, titleID).Scan(&existingVolumeID)
	if existingVolumeID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "Такой том уже существует"})
		return
	}

	volume := models.Volume{
		Name:         name,
		Description:  description,
		TitleID:      titleID,
		CreatorID:    claims.ID,
		OnModeration: true,
	}

	var volumeCover struct {
		VolumeID uint   `bson:"volume_id"`
		Cover    []byte `bson:"cover"`
	}

	tx := h.DB.Begin()

	if result := tx.Create(&volume); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error})
		return
	}

	file, err := cover.Open()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	volumeCover.VolumeID = volume.ID
	volumeCover.Cover = data

	if _, err := h.Collection.InsertOne(context.Background(), volumeCover); err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "Том успешно отправлен на модерацию"})
}
