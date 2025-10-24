package chapters_pages_compressor

func (c *ChaptersPagesCompressor) EnqueueChapterOnModerationID(chapterOnModerationID uint) {
	c.chaptersOnModerationIDs <- chapterOnModerationID
}
