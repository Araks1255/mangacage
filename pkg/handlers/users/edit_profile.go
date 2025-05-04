package users

import (
	"context"
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) EditProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["userName"]) == 0 && len(form.Value["aboutYourself"]) == 0 && len(form.File["profilePicture"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	if len(form.File["profilePicture"]) != 0 && form.File["profilePicture"][0].Size > 10<<20 {
		c.AbortWithStatusJSON(400, gin.H{"error": "превышен максимальный размер аватарки (10мб)"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if len(form.Value["userName"]) != 0 {
		var existing struct {
			UserOnModerationID uint
			UserID             uint
		}

		tx.Raw(
			`SELECT
				(SELECT id FROM users_on_moderation WHERE lower(user_name) = lower(?)) AS user_on_moderation_id,
				(SELECT id FROM users WHERE lower(user_name) = lower(?)) AS user_id`,
			form.Value["userName"][0], form.Value["userName"][0],
		).Scan(&existing)

		if existing.UserOnModerationID != 0 {
			c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже ожидает верификации аккаунта"})
			return
		}
		if existing.UserID != 0 {
			c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже существует"})
			return
		}
	}

	editedProfile := models.UserOnModeration{
		ExistingID: sql.NullInt64{Int64: int64(claims.ID), Valid: true},
	}

	if len(form.Value["userName"]) != 0 {
		editedProfile.UserName = sql.NullString{String: form.Value["userName"][0]}
	}
	if len(form.Value["aboutYourself"]) != 0 {
		editedProfile.AboutYourself = form.Value["aboutYourself"][0]
	}

	if result := tx.Raw(
		`INSERT INTO users_on_moderation (created_at, user_name, about_yourself, existing_id)
		VALUES(NOW(), ?, ?, ?)
		ON CONFLICT (existing_id) DO UPDATE
		SET
			updated_at = EXCLUDED.created_at,
			user_name = EXCLUDED.user_name,
			about_yourself = EXCLUDED.about_yourself
		RETURNING id`,
		editedProfile.UserName, editedProfile.AboutYourself, editedProfile.ExistingID,
	).Scan(&editedProfile.ID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error})
		return
	}

	if len(form.File["profilePicture"]) != 0 {
		profilePicture, err := utils.ReadMultipartFile(form.File["profilePicture"][0], 10<<20)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"user_on_moderation_id": editedProfile.ID}
		update := bson.M{"$set": bson.M{"profile_picture": profilePicture}}
		opts := options.Update().SetUpsert(true)

		if _, err = h.UsersOnModerationProfilePictures.UpdateOne(context.Background(), filter, update, opts); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения профиля успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutUserOnModeration(context.Background(), &pb.User{ID: uint64(editedProfile.ExistingID.Int64), New: false}); err != nil {
		log.Println(err)
	}
}
