package article

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"size:255;not null"`
	Slug      string         `gorm:"size:255;uniqueIndex;not null"`
	Content   string         `gorm:"type:text;not null"`
	CreatedAt time.Time      
	UpdatedAt time.Time      
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
