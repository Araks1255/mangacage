package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/search"
)

func TestSearch(t *testing.T) {
	scenarios := search.GetSearchScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
