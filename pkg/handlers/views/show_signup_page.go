package views

import "github.com/gin-gonic/gin"

func (h handler) ShowSignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}
