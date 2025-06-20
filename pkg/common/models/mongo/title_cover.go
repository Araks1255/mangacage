package mongo

type TitleCover struct {
	TitleID   uint   `bson:"title_id"`
	CreatorID uint   `bson:"creator_id"`
	Cover     []byte `bson:"cover"`
}

type TitleOnModerationCover struct {
	TitleOnModerationID uint   `bson:"title_on_moderation_id"`
	CreatorID           uint   `bson:"creator_id"`
	Cover               []byte `bson:"cover"`
}
