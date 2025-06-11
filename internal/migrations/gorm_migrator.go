package migrations

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"

	"gorm.io/gorm"
)

func GormMigrate(db *gorm.DB) error {
	db.Exec(
		`CREATE TABLE users (
    		id BIGSERIAL PRIMARY KEY,
   			user_name TEXT,
    		team_id BIGINT
		)`,
	)

	db.Exec(
		`CREATE TABLE teams (
    		id BIGSERIAL PRIMARY KEY,
    		name TEXT,
    		creator_id BIGINT,
    		moderator_id BIGINT
		)`,
	)

	sqlTypesDir := "./internal/migrations/sql/types"
	pathsToScriptsWithTypes := make([]string, 0, 3)

	filepath.WalkDir(sqlTypesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if !d.IsDir() {
			if !strings.HasSuffix(d.Name(), ".sql") {
				panic("в каталоге с sql миграциями не sql файл")
			}

			pathsToScriptsWithTypes = append(pathsToScriptsWithTypes, path)
		}

		return nil
	})

	for i := 0; i < len(pathsToScriptsWithTypes); i++ {
		scriptWithTypeBytes, err := os.ReadFile(pathsToScriptsWithTypes[i])
		if err != nil {
			panic(err)
		}

		db.Exec(string(scriptWithTypeBytes))
	}

	err := db.AutoMigrate(
		&models.Role{},
		&models.Genre{},
		&models.Tag{},
	)
	if err != nil {
		return err
	}

	if err = db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&models.Team{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(
		&models.Author{},
		&models.Title{},
		&models.Volume{},
		&models.Chapter{},
	); err != nil {
		return err
	}

	if err = db.AutoMigrate(
		&models.TitleOnModeration{},
		&models.VolumeOnModeration{},
		&models.ChapterOnModeration{},
		&models.UserOnModeration{},
		&models.TeamOnModeration{},
	); err != nil {
		return err
	}

	if err = db.AutoMigrate(
		&models.TeamJoinRequest{},
		&models.UserViewedChapter{},
	); err != nil {
		return err
	}

	sqlDir := "./internal/migrations/sql"
	pathsToScripts := make([]string, 0, 20)

	filepath.WalkDir(sqlDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if d.IsDir() && d.Name() == "types" {
			return fs.SkipDir
		}

		if !d.IsDir() {
			if !strings.HasSuffix(d.Name(), ".sql") {
				panic("в каталоге с sql миграциями не sql файл")
			}

			pathsToScripts = append(pathsToScripts, path)
		}

		return nil
	})

	for i := 0; i < len(pathsToScripts); i++ {
		scriptBytes, err := os.ReadFile(pathsToScripts[i])
		if err != nil {
			panic(err)
		}

		if err := db.Exec(string(scriptBytes)).Error; err != nil {
			panic(err)
		}
	}

	return nil
}
