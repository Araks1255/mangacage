package teams

// import (
// 	//"fmt"
// 	"log"
// 	//"strings"

// 	"github.com/Araks1255/mangacage/pkg/common/models"
// 	"github.com/gin-gonic/gin"
// 	//"gorm.io/gorm"
// )

// func (h handler) CreateTeam(c *gin.Context) {
// 	claims := c.MustGet("claims").(*models.Claims)

// 	if claims.Role == "team_owner" {
// 		c.AbortWithStatusJSON(403, gin.H{"error":"Вы уже владеете командой перевода"})
// 	}

// 	var newTeam models.Team

// 	if err := c.ShouldBindJSON(&newTeam); err != nil {
// 		log.Println(err)
// 		c.AbortWithStatusJSON(400, gin.H{"error":err.Error()})
// 		return
// 	}

// 	// Не забыть присваивание роли тим овнер
// }
