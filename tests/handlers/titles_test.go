package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/titles"
)

func TestCreateTitle(t *testing.T) {
	scenarios := titles.GetCreateTitleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestDeleteTitle(t *testing.T) {
	// Логика удаления скорее всего будет изменена, поэтому тестов пока нет
}

func TestEditTitle(t *testing.T) {
	scenarios := titles.GetEditTitleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMostPopularTitles(t *testing.T) {
	scenarios := titles.GetGetMostPopularTitlesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetNewTitles(t *testing.T) {
	scenarios := titles.GetGetNewTitlesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetRecentlyUpdatedTitles(t *testing.T) {
	scenarios := titles.GetGetRecentlyUpdatedTitlesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTitleCover(t *testing.T) {
	scenarios := titles.GetGetTitleCoverScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTitle(t *testing.T) {
	scenarios := titles.GetGetTitleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestQuitTranslatingTitle(t *testing.T) { // Эта логика скорее всего будет изменена, поэтому тестов на неё нет
}

func TestSubscribeToTitle(t *testing.T) {
	scenarios := titles.GetSubscribeToTitleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestTranslateTitle(t *testing.T) { // Аналогично quit translating title
}
