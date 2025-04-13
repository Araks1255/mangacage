package seeder

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func seedUsersAndTeams(ctx context.Context, db *gorm.DB, usersProfilePictures, teamsCovers *mongo.Collection) error {
	var userID uint
	if result := db.Raw(
		`INSERT INTO users (user_name, password, about_yourself, tg_user_id, team_id)
		VALUES ('user_test', '', '', 0, NULL)
		ON CONFLICT DO NOTHING
		RETURNING id`,
	).Scan(&userID); result.Error != nil {
		return result.Error
	}

	if userID == 0 { // Тут транзакция, так что ситуации, когда есть юзер и нет команды случиться не должно
		log.Println("Юзер уже создан")
		return nil
	}

	if result := db.Exec(
		`INSERT INTO user_roles (user_id, role_id)
		SELECT u.id, r.id
		FROM (SELECT id FROM users WHERE user_name = 'user_test') AS u
		CROSS JOIN UNNEST(ARRAY['user', 'admin', 'team_leader']) AS role_name
		JOIN roles AS r ON r.name = role_name
		ON CONFLICT DO NOTHING`,
	); result.Error != nil {
		return result.Error
	}

	var userProfilePicture struct {
		UserID         uint   `bson:"user_id"`
		ProfilePicture []byte `bson:"profile_picture"`
	}

	var err error
	userProfilePicture.UserID = userID
	userProfilePicture.ProfilePicture, err = os.ReadFile("test_data/user_profile_picture.png")
	if err != nil {
		return err
	}

	if _, err = usersProfilePictures.InsertOne(ctx, userProfilePicture); err != nil {
		return err
	}

	var teamID uint
	if result := db.Raw(
		`INSERT INTO teams (name, description, creator_id, moderator_id)
		SELECT 'team_test', '', u.id, u.id FROM
		(SELECT id FROM users WHERE user_name = 'user_test') AS u
		ON CONFLICT DO NOTHING
		RETURNING id `,
	).Scan(&teamID); result.Error != nil {
		return result.Error
	}

	if teamID == 0 {
		log.Println("Команда уже существует")
		return nil
	}

	var teamCover struct {
		TeamID uint   `bson:"team_id"`
		Cover  []byte `bson:"cover"`
	}

	teamCover.TeamID = teamID
	teamCover.Cover, err = os.ReadFile("test_data/team_cover.png")
	if err != nil {
		return err
	}

	if _, err = teamsCovers.InsertOne(ctx, teamCover); err != nil {
		return err
	}

	if result := db.Exec(
		`UPDATE users SET team_id = (SELECT id FROM teams WHERE name = 'team_test')`,
	); result.Error != nil {
		return result.Error
	}

	return nil
}
