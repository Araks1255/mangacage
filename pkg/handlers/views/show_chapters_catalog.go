package views

import "github.com/gin-gonic/gin"

func (h handler) ShowChaptersCatalogPage(c *gin.Context) {
	c.HTML(200, "chapters_catalog.html", nil)
}