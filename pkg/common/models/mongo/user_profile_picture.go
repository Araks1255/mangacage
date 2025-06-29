package mongo

type UserProfilePicture struct {
	UserID         uint   `bson:"user_id"`
	CreatorID      uint   `bson:"creator_id"`
	ProfilePicture []byte `bson:"profile_picture"`
	Visible        bool
}

type UserOnModerationProfilePicture struct {
	UserOnModerationID uint   `bson:"user_on_moderation_id"`
	CreatorID          uint   `bson:"creator_id"`
	ProfilePicture     []byte `bson:"profile_picture"`
}
