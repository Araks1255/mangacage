package handlers

import (
	"fmt"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/handlers/chapters"
	"github.com/lib/pq"
)

func TestNewPagesCreation(t *testing.T) {
	pages := make([][]byte, 0, 70)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 35; i++ {
		page, err := os.ReadFile("./test_data/.jpg")
		if err != nil {
			t.Fatal(err)
		}
		pages = append(pages, page)
	}
	for i := 0; i < 35; i++ {
		page, err := os.ReadFile("./test_data/chapter_page9.jpg")
		if err != nil {
			t.Fatal(err)
		}
		pages = append(pages, page)
	}

	rows, err := env.DB.Raw("EXPLAIN ANALYZE INSERT INTO chapters_pages (chapter_id, page) select ?, UNNEST(?::BYTEA[])", chapterID, pq.Array(pages)).Rows()
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			t.Fatal(err)
		}
		fmt.Println(line)
	}
} hh

func TestCreateChapter(t *testing.T) {
	scenarios := chapters.GetCreateChapterScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

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

func TestGetChapters(t *testing.T) {
	scenarios := chapters.GetGetChaptersScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
