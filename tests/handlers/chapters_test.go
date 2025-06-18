package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/chapters"
)

func TestCreateChapter(t *testing.T) {
	scenarios := chapters.GetCreateChapterScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

// func TestDeleteChapter(t *testing.T) {
// 	scenarios := chapters.GetDeleteChapterScenarios(env)

// 	for name, scenario := range scenarios {
// 		t.Run(name, scenario)
// 	}
// }

func TestEditChapter(t *testing.T) {
	scenarios := chapters.GetEditChapterScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetChapterPage(t *testing.T) {
	scenarios := chapters.GetGetChapterPageScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetChapter(t *testing.T) {
	scenarios := chapters.GetGetChapterScenarios(env)

	for name, scescenario := range scenarios {
		t.Run(name, scescenario)
	}
}

func TestGetVolumeChapters(t *testing.T) {
	scenarios := chapters.GetGetVolumeChaptersScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
