package seeder

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, mongoDB *mongo.Database, mode string) error {
	var err error

	ctx := context.TODO()

	titlesCovers := mongoDB.Collection("titles_covers")
	chaptersPages := mongoDB.Collection("chapters_pages")
	usersProfilePictures := mongoDB.Collection("users_profile_pictures")
	teamsCovers := mongoDB.Collection("teams_covers")

	tx := db.Begin()
	defer tx.Rollback()
	defer func() {
		if err != nil {
			cleanUpCollections(titlesCovers, chaptersPages, usersProfilePictures, teamsCovers)
		}
	}()

	switch mode {
	case "test":
		if err = seedRoles(tx); err != nil {
			return err
		}
		if err = seedGenres(tx); err != nil {
			return err
		}
		if err = seedUsersAndTeams(ctx, tx, usersProfilePictures, teamsCovers); err != nil {
			return err
		}
		if err = seedAuthors(tx); err != nil {
			return err
		}
		if err = seedTitles(ctx, tx, titlesCovers); err != nil {
			return err
		}
		if err = seedVolumes(tx); err != nil {
			return err
		}
		if err = seedChapters(ctx, tx, chaptersPages); err != nil {
			return err
		}

	case "prod":
		if err = seedRoles(tx); err != nil {
			return err
		}

	default:
		return errors.New("Неккоректный тип сида")
	}

	tx.Commit()

	return nil
}

func cleanUpCollections(titlesCovers, chaptersPages, usersProfilePictures, teamsCovers *mongo.Collection) {
	titlesCovers.DeleteMany(nil, bson.M{})
	chaptersPages.DeleteMany(nil, bson.M{})
	usersProfilePictures.DeleteMany(nil, bson.M{})
	teamsCovers.DeleteMany(nil, bson.M{})
}
