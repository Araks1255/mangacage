package constraints

const (
	UniUsersOnModerationUsername  = "uni_users_on_moderation_user_name"
	UserFavoriteTitlesPkey        = "user_favorite_titles_pkey"
	FkUserFavoriteTitlesTitle     = "fk_user_favorite_titles_title"
	UserFavoriteChaptersPkey      = "user_favorite_chapters_pkey"
	FkUserFavoriteChaptersChapter = "fk_user_favorite_chapters_chapter"
	UserFavoriteGenresPkey        = "user_favorite_genres_pkey"
	FkUserFavoriteGenresGenre     = "fk_user_favorite_genres_genre"
	UniTitlesOnModerationName     = "uni_titles_on_moderation_name"
	FkTitlesOnModerationAuthor    = "fk_titles_on_moderation_author"
	FkTitleGenresGenre            = "fk_title_genres_genre"
	FkVolumesTitle                = "fk_volumes_title"
	UniVolumeTitle                = "uniq_volume_title"
	FkChaptersVolume              = "fk_chapters_volume"
	UniChapterVolume              = "uniq_chapter_volume"
	FkChaptersOnModerationVolume  = "fk_chapters_on_moderation_volume"
	UniTeamsOnModerationName      = "uni_teams_on_moderation_name"
	UniTeamsOnModerationCreatorID = "uni_teams_on_moderation_creator_id"
	UserRolesPkey                 = "user_roles_pkey"
	FkUserRolesRole               = "fk_user_roles_role"
)
