package views

import "github.com/gin-gonic/gin"

func (h handler) ShowProfileChangesCatalog(c *gin.Context) {
	c.HTML(200, "profile_changes_catalog.html", nil)
}