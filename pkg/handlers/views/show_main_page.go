package views

import "github.com/gin-gonic/gin"

func (h handler) ShowMainPage(c *gin.Context) {
	c.File("./html/main_page.html")
}
