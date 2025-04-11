package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTitlePage(c *gin.Context) {
	c.HTML(200, "title_page.html", gin.H{})
}
