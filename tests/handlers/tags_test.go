package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/tags"
)

func TestAddTag(t *testing.T) {
	scenarios := tags.GetAddTagScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTags(t *testing.T) {
	scenarios := tags.GetGetTagsScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
