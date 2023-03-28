package model

type MovieModel struct {
	ID     int      `gorm:"primaryKey"`
	Title  string   `gorm:"not null"`
	Year   int      `gorm:"not null"`
	Rating float32  `gorm:"not null"`
	Genres []string `gorm:"-"`
}
