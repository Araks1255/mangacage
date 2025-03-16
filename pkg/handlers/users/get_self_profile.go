package users

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h handler) GetSelfProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var user struct {
		UserName                 string
		AboutYourself            string
		Team                     string
		RegistrationDate         time.Time
		Roles                    pq.StringArray `gorm:"type:text[]"`
		TitlesUserIsSubscribedTo pq.StringArray `gorm:"type:text[]"`
	}

	h.DB.Raw(
		`SELECT u.user_name, u.about_yourself, u.created_at AS registration_date, teams.name AS team,
		(
			SELECT ARRAY(
				SELECT roles.name FROM roles
				INNER JOIN user_roles ON roles.id = user_roles.role_id
				WHERE user_roles.user_id = u.id
			) AS roles
		),
		(
			SELECT ARRAY(
				SELECT titles.name FROM titles
				INNER JOIN user_titles_subscribed_to ON titles.id = user_titles_subscribed_to.title_id
				where user_titles_subscribed_to.user_id = u.id
			) AS titles_user_is_subscribed_to
		) FROM users AS u
		INNER JOIN teams ON teams.id = u.team_id
		WHERE u.id = ?`,
		claims.ID,
	).Scan(&user)

	if user.UserName == "" {
		c.AbortWithStatusJSON(404, gin.H{"error": "не удалось получитб профиль"})
		return
	}

	c.JSON(200, &user)
}
