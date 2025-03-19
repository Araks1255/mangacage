package teams

import (
	"context"
	"io"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

type TeamCover struct {
	TeamID uint   `bson:"team_id"`
	Cover  []byte `bson:"cover"`
}

func (h handler) CreateTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID).Scan(&userRoles)

	if slices.Contains(userRoles, "team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже являетесь владельцем другой команды"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе нет имени команды"})
		return
	}

	name := form.Value["name"][0]
	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	newTeam := models.Team{
		Name:        name,
		Description: description,
	}

	transaction := h.DB.Begin()

	if result := transaction.Create(&newTeam); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось создать команду"})
		return
	}

	if result := transaction.Exec("UPDATE users SET team_id = ? WHERE id = ?", newTeam.ID, claims.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось присоеденить вас к команде"})
		return
	}

	if result := transaction.Exec(`INSERT INTO user_roles (user_id, role_id)
		VALUES (?, (SELECT id FROM roles WHERE name = 'team_leader')),
		(?, (SELECT id FROM roles WHERE name = 'translater'))`,
		claims.ID, claims.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось назначить вас лидером команды"})
		return
	}

	transaction.Commit()

	c.JSON(201, gin.H{"success": "Команда успешно создана, и вы являетесь её лидером"})

	cover, err := c.FormFile("cover")
	if err != nil {
		log.Println(err)
		return
	}

	file, err := cover.Open()
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	teamCover := TeamCover{
		TeamID: newTeam.ID,
		Cover:  data,
	}

	if _, err := h.Collection.InsertOne(context.Background(), teamCover); err != nil {
		log.Println(err)
	}
}
