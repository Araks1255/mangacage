package views

import "github.com/gin-gonic/gin"

func (h handler) ShowChapterOnModerationPage(c *gin.Context) {
	c.HTML(200, "chapter_on_moderation_page.html", nil)
}