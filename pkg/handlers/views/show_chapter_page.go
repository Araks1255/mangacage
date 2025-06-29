package views

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) ShowChapterPage(c *gin.Context) {
	c.HTML(200, "chapter_page.html", gin.H{})

	claims, ok := c.Get("claims")
	if !ok {
		return
	}

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return
	}

	view := models.UserViewedChapter{
		UserID:    claims.(*auth.Claims).ID,
		ChapterID: uint(chapterID),
	}

	h.DB.Create(&view)
}
