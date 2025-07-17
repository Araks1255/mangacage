package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/users"
	"github.com/Araks1255/mangacage/tests/handlers/users/favorites"
	"github.com/Araks1255/mangacage/tests/handlers/users/moderation"
)

// Users
func TestEditProfile(t *testing.T) {
	scenarios := users.GetEditProfileScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyProfile(t *testing.T) {
	scenarios := users.GetGetMyProfileScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyProfilePicture(t *testing.T) {
	scenarios := users.GetGetMyProfilePictureScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetUser(t *testing.T) {
	scenarios := users.GetGetUserScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetUsers(t *testing.T) {
	scenarios := users.GetGetUsersScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

// Favorites

func TestAddTitleToFavorites(t *testing.T) {
	scenarios := favorites.GetAddTitleToFavoritesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestDeleteTitleFromFavorites(t *testing.T) {
	scenarios := favorites.GetDeleteTitleFromFavoritesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestAddGenreToFavorites(t *testing.T) {
	scenarios := favorites.GetAddGenreToFavoritesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestDeleteGenreFromFavorites(t *testing.T) {
	scenarios := favorites.GetDeleteGenreFromFavoritesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestAddChapterToFavorites(t *testing.T) {
	scenarios := favorites.GetAddChapterToFavoritesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestDeleteChapterFromFavorites(t *testing.T) {
	scenarios := favorites.GetDeleteChapterFromFavoritesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

// Moderation

func TestCancelAppealForModeration(t *testing.T) {
	scenarios := moderation.GetCancelAppealForModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestCancelAppealForProfileChanges(t *testing.T) {
	scenarios := moderation.GetCancelAppealForProfileChangesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyProfilePictureOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyProfilePictureOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyChapterOnModerationPage(t *testing.T) {
	scenarios := moderation.GetGetMyChapterOnModerationPageScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyTitleOnModerationCover(t *testing.T) {
	scenarios := moderation.GetGetMyTitleOnModerationCoverScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyTitlesOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyTitlesOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyChaptersOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyChaptersOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyGenresOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyGenresOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyTagsOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyTagsOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyAuthorsOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyAuthorsOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyTeamOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyTeamOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyProfileChangesOnModeration(t *testing.T) {
	scenarios := moderation.GetGetMyProfileChangesOnModerationScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
