package views

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h handler) ShowReadingPage(c *gin.Context) {
	c.HTML(200, "reading_page.html", gin.H{})

	cookie, err := c.Cookie("mangacage_token")
	if err != nil {
		return
	}

	title := c.Param("title")
	volume := c.Param("volume")
	chapter := c.Param("chapter")

	var chapterID uint
	h.DB.Raw(`SELECT chapters.id FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE lower(titles.name) = lower(?)
		AND lower(volumes.name) = lower(?)
		AND lower(chapters.name) = lower(?)
		AND NOT chapters.on_moderation`,
		title,
		volume,
		chapter).Scan(&chapterID)

	if chapterID == 0 {
		return
	}

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	claims, err := utils.ParseToken(cookie, secretKey)
	if err != nil {
		return
	}

	if result := h.DB.Exec("INSERT INTO user_viewed_chapters (user_id, chapter_id) VALUES (?, ?)", claims.ID, chapterID); result.Error != nil {
		log.Println(result.Error)
		return
	}
}
