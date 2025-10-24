package views

import "github.com/gin-gonic/gin"

func (h handler) ShowRules(c *gin.Context) {
	c.HTML(200, "rules.html", nil)
}
