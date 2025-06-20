package mongo

type TeamCover struct {
	TeamID    uint   `bson:"team_id"`
	CreatorID uint   `bson:"creator_id"`
	Cover     []byte `bson:"cover"`
}

type TeamOnModerationCover struct {
	TeamOnModerationID uint   `bson:"team_on_moderation_id"`
	CreatorID          uint   `bson:"creator_id"`
	Cover              []byte `bson:"cover"`
}
