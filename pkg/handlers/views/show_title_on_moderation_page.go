package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTitleOnModerationPage(c *gin.Context) {
	c.HTML(200, "title_on_moderation_page.html", nil)
}