package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTeamJoinRequestsCatalog(c *gin.Context) {
	c.HTML(200, "team_join_requests_catalog.html", nil)
}
