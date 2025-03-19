package volumes

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")

	var requestBody struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE lower(name) = lower(?)", title).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Тайтл не найден"})
		return
	}

	var existingVolumeID uint
	h.DB.Raw("SELECT id FROM volumes WHERE lower(name) = lower(?) AND title_id = ?", requestBody.Name, titleID).Scan(&existingVolumeID)
	if existingVolumeID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "Такой том уже существует"})
		return
	}

	volume := models.Volume{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		TitleID:     titleID,
		CreatorID:   claims.ID,
	}

	if result := h.DB.Create(&volume); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error})
		return
	}

	c.JSON(201, gin.H{"success": "Том успешно отправлен на модерацию"})
}
