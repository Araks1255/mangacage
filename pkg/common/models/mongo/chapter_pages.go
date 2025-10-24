package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChapterPages struct {
	ChapterID uint                 `bson:"chapter_id"`
	CreatorID uint                 `bson:"creator_id"`
	PagesIDs  []primitive.ObjectID `bson:"pages_ids"`
}

type ChapterOnModerationPages struct {
	ChapterOnModerationID uint                 `bson:"chapter_on_moderation_id"`
	CreatorID             uint                 `bson:"creator_id"`
	PagesIDs              []primitive.ObjectID `bson:"pages_ids"`
}
