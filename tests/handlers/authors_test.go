package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/authors"
)

func TestAddAuthor(t *testing.T) {
	scenarios := authors.GetAddAuthorScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetAuthors(t *testing.T) {
	scenarios := authors.GetGetAuthorsScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetAuthor(t *testing.T) {
	scenarios := authors.GetGetAuthorScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
