package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTitlesCatalogPage(c *gin.Context) {
	c.HTML(200, "titles_catalog.html", nil)
}
