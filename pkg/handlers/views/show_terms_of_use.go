package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTermsOfUse(c *gin.Context) {
	c.HTML(200, "terms_of_use.html", nil)
}
