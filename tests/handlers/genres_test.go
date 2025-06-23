package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/genres"
)

func TestAddGenre(t *testing.T) {
	scenarios := genres.GetAddGenreScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetGenres(t *testing.T) {
	scenarios := genres.GetGetGenresScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
