package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTitleTranslateRequestsCatalog(c *gin.Context) {
	c.HTML(200, "titles_translate_requests_catalog.html", nil)
}
