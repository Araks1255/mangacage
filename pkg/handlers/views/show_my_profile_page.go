package views

import (

	"github.com/gin-gonic/gin"
)

func (h handler) ShowMyProfilePage(c *gin.Context) {
	c.HTML(200, "my_profile_page.html", nil)
}
