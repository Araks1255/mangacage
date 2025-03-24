package moderation

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetSelfProfileChangesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var editedProfile struct {
		CreatedAt     time.Time
		UserName      string
		AboutYourself string
	}

	h.DB.Raw(
		`SELECT u.created_at, u.user_name, u.about_yourself
		FROM users_on_moderation AS u
		WHERE u.existing_id = ?`, claims.ID,
	).Scan(&editedProfile)

	response := make(map[string]string, 3)
	if editedProfile.UserName != "" {
		response["newUserName"] = editedProfile.UserName
	}
	if editedProfile.AboutYourself != "" {
		response["newAboutYourself"] = editedProfile.AboutYourself
	}
	response["createdAt"] = editedProfile.CreatedAt.Format(time.DateTime)

	c.JSON(200, &response)
}
