package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTeamOnModerationCatalog(c *gin.Context) {
	c.HTML(200, "team_on_moderation_catalog.html", nil)
}
