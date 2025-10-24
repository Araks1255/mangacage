package views

import "github.com/gin-gonic/gin"

func (h handler) ShowUsersCatalog(c *gin.Context) {
	c.HTML(200, "users_catalog.html", nil)
}
