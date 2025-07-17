package views

import "github.com/gin-gonic/gin"

func (h handler) ShowProfilePage(c *gin.Context) {
	c.HTML(200, "profile_page.html", nil)
}