package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTeamPage(c *gin.Context) {
	c.HTML(200, "team_page.html", gin.H{})
}