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

func TestGetTitles(t *testing.T) {
	scenarios := titles.GetGetTitlesScenarios(env)

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

func TestQuitTranslatingTitle(t *testing.T) {
	scenarios := titles.GetQuitTranslatingTitleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestSubscribeToTitle(t *testing.T) {
	scenarios := titles.GetSubscribeToTitleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestTranslateTitle(t *testing.T) {
	scenarios := titles.GetTranslateTitleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestCancelTitleTranslateRequest(t *testing.T) {
	scenarios := titles.GetCancelTitleTranslateRequestScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTitleTranslateRequests(t *testing.T) {
	scenarios := titles.GetGetTitleTranslateRequests(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
