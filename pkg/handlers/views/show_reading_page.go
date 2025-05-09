package views

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
)

func (h handler) ShowReadingPage(c *gin.Context) {
	c.HTML(200, "reading_page.html", nil)

	cookie, err := c.Cookie("mangacage_token")
	if err != nil {
		return
	}

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	claims, err := utils.ParseToken(cookie, secretKey)
	if err != nil {
		return
	}

	view := models.UserViewedChapter{
		UserID:    claims.ID,
		ChapterID: uint(chapterID),
	}

	if err := h.DB.Create(&view).Error; err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return
		}

		log.Println(err)
	}
}
