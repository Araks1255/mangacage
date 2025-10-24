package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTeamsCatalog(c *gin.Context) {
	c.HTML(200, "teams_catalog.html", nil)
}
