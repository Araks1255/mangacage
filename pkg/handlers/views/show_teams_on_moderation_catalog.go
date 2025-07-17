package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTeamsOnModerationCatalog(c *gin.Context) {
	c.HTML(200, "teams_on_moderation_catalog.html", nil)
}