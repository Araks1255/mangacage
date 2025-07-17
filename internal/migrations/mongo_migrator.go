package migrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionsAndEntitiesNames = map[string]string{
	mongodb.ChaptersPagesCollection:        "chapter",
	mongodb.TitlesCoversCollection:         "title",
	mongodb.TeamsCoversCollection:          "team",
	mongodb.UsersProfilePicturesCollection: "user",
}

func mongoMigrate(ctx context.Context, db *mongo.Database) error {
	for collectionName, entityName := range collectionsAndEntitiesNames {
		err := migrateCollection(ctx, db, collectionName, entityName)
		if err != nil {
			return err
		}
	}
	return nil
}

func migrateCollection(ctx context.Context, db *mongo.Database, collectionName, entityName string) error {
	var cmdErr mongo.CommandError

	if err := db.CreateCollection(ctx, collectionName); errors.As(err, &cmdErr) && cmdErr.Code != 48 { // Код повторного создания коллекции
		return err
	}

	collection := db.Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{fmt.Sprintf("%s_id", entityName): 1},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.M{fmt.Sprintf("%s_on_moderation_id", entityName): 1},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	}

	if entityName == "team" || entityName == "user" {
		indexes = append(indexes, mongo.IndexModel{
			Keys:    bson.M{"creator_id": 1},
			Options: options.Index().SetUnique(true),
		})
	}

	if _, err := collection.Indexes().CreateMany(ctx, indexes); errors.As(err, &cmdErr) && cmdErr.Code != 85 && cmdErr.Code != 86 { // Это коды повторного создания
		return err
	}

	return nil
}
