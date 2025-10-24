package chapters_pages_compressor

import (
	"context"

	"gorm.io/gorm"
)

type ChaptersPagesCompressor struct {
	chaptersOnModerationIDs chan uint
	db                      *gorm.DB
	ctx                     context.Context
}

func NewChaptersPagesCompressor(ctx context.Context, db *gorm.DB, buffer int) *ChaptersPagesCompressor {
	return &ChaptersPagesCompressor{
		chaptersOnModerationIDs: make(chan uint, buffer),
		db:                      db,
		ctx:                     ctx,
	}
}
