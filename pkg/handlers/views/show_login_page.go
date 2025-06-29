package views

import "github.com/gin-gonic/gin"

func (h handler) ShowLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}