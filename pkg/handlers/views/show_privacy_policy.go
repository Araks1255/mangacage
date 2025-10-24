package views

import "github.com/gin-gonic/gin"

func (h handler) ShowPrivacyPolicy(c *gin.Context) {
	c.HTML(200, "privacy_policy.html", nil)
}
