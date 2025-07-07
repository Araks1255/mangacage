package teams

import (
	"context"
	"mime/multipart"

	"github.com/Araks1255/mangacage/pkg/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpsertTeamOnModerationCover(ctx context.Context, collection *mongo.Collection, coverFileHeader *multipart.FileHeader, teamOnModerationID, userID uint) error {
	cover, err := utils.ReadMultipartFile(coverFileHeader, 2<<20)
	if err != nil {
		return err
	}

	filter := bson.M{"team_on_moderation_id": teamOnModerationID}
	update := bson.M{"$set": bson.M{"cover": cover, "creator_id": userID}}
	opts := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return err
	}

	return nil
}
