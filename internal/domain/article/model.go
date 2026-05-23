package article

import (
	"time"

	"github.com/nullablenone/go-news-api/internal/domain/user"

	"gorm.io/gorm"
)

type Article struct {
	ID       uint   `gorm:"primaryKey"`
	Title    string `gorm:"size:255;not null"`
	Slug     string `gorm:"size:255;uniqueIndex;not null"`
	Summary  string `gorm:"type:text;not null"`
	Category string `gorm:"size:100;not null"`
	ReadTime string `gorm:"size:50;not null"`

	AuthorID uint      `gorm:"not null"`
	Author   user.User `gorm:"foreignKey:AuthorID"`

	Content []map[string]interface{} `gorm:"serializer:json;type:jsonb;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
