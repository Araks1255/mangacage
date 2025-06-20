package mongo

type ChapterPages struct {
	ChapterID uint     `bson:"chapter_id"`
	CreatorID uint     `bson:"creator_id"`
	Pages     [][]byte `bson:"pages"`
}

type ChapterOnModerationPages struct {
	ChapterOnModerationID uint     `bson:"chapter_on_moderation_id"`
	CreatorID             uint     `bson:"creator_id"`
	Pages                 [][]byte `bson:"pages"`
}
