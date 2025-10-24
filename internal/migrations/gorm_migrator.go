package migrations

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"
)

func gormMigrate(db *gorm.DB) error {
	if err := createStubTables(db); err != nil {
		return err
	}

	if err := executeSQLScripts(db, "./internal/migrations/sql/extensions"); err != nil {
		return err
	}
	if err := executeSQLScripts(db, "./internal/migrations/sql/types"); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Role{}, &models.Genre{}, &models.Tag{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Team{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Author{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Title{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Chapter{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.UserOnModeration{},
		&models.TeamOnModeration{},
		&models.GenreOnModeration{},
		&models.TagOnModeration{},
		&models.AuthorOnModeration{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.TitleOnModeration{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.ChapterOnModeration{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Page{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.TeamJoinRequest{},
		&models.TitleTranslateRequest{},
		&models.UserViewedChapter{},
		&models.TitleRate{},
	); err != nil {
		return err
	}

	if err := executeSQLScripts(db, "./internal/migrations/sql", "extensions", "types"); err != nil {
		return err
	}

	return nil
}

func createStubTables(db *gorm.DB) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS teams (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS authors_on_moderation (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS chapters_on_moderation (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS genres_on_moderation (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS tags_on_moderation (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS teams_on_moderation (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS titles_on_moderation (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS users_on_moderation (id BIGSERIAL PRIMARY KEY)`).Error; err != nil {
		return err
	}
	return nil
}

func executeSQLScripts(db *gorm.DB, rootDir string, ignoreDirsNames ...string) error {
	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		for i := 0; i < len(ignoreDirsNames); i++ {
			if d.IsDir() && d.Name() == ignoreDirsNames[i] {
				return fs.SkipDir
			}
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			scriptBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if err := db.Exec(string(scriptBytes)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
