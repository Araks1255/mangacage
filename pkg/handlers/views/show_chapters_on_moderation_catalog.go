package views

import "github.com/gin-gonic/gin"

func (h handler) ShowChaptersOnModerationCatalog(c *gin.Context) {
	c.HTML(200, "chapters_on_moderation_catalog.html", nil)
}