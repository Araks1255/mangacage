package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTitlesOnModerationCatalog(c *gin.Context) {
	c.HTML(200, "titles_on_moderation_catalog.html", nil)
}