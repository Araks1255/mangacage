package models

type Genre struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"unique"`
}

type GenreDTO struct { // На будущее, если вдруг появится описание жанра или ещё что-то
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
