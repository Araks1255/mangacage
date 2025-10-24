package views

import "github.com/gin-gonic/gin"

func (h handler) ShowCopyright(c *gin.Context) {
	c.HTML(200, "copyright.html", nil)
}
