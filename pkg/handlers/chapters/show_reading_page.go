package chapters

import (
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h handler) ShowReadingPage(c *gin.Context) {
	c.String(200, "Типа глава\n\nЧитай")
	//c.HTML(200, "page", gin.H{})

	title := strings.ToLower(c.Param("title"))
	volume := strings.ToLower(c.Param("volume"))
	chapter := strings.ToLower(c.Param("chapter"))

	var titleID uint
	h.DB.Raw(`SELECT titles.id FROM titles
		INNER JOIN volumes ON titles.id = volumes.title_id
		INNER JOIN chapters ON volumes.id = chapters.volume_id
		WHERE titles.name = ?
		AND volumes.name = ?
		AND chapters.name = ?`,
		title,
		volume,
		chapter,
	).Scan(&titleID)

	if titleID == 0 {
		return
	}

	cookie, err := c.Cookie("mangacage_token")
	if err != nil {
		return
	}

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	claims, err := utils.ParseToken(cookie, secretKey)
	if err != nil {
		return
	}

	if result := h.DB.Exec("INSERT INTO user_viewed_titles (user_id, title_id) VALUES (?, ?)", claims.ID, titleID); result.Error != nil {
		return
	}
}
