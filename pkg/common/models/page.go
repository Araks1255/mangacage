package models

type Page struct {
	ID     uint `gorm:"primaryKey;autoIncrement:true"`
	Number uint `gorm:"not null;check:page_number_positive_chk,number > 0"`

	ChapterID *uint
	Chapter   Chapter `gorm:"foreignKey:ChapterID;references:id;constraint:OnDelete:CASCADE"`

	ChapterOnModerationID *uint

	Path   string `gorm:"not null"`
	Format string `gorm:"not null"`

	Hidden bool `gorm:"not null;default:false"`
}
