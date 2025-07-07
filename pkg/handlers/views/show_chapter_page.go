package views

import (
	"github.com/gin-gonic/gin"
)

func (h handler) ShowChapterPage(c *gin.Context) {
	c.HTML(200, "chapter_page.html", gin.H{})
}
