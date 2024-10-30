package models

import (
	"time"
)

type Blog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"size:255;not null;index" json:"title" validate:"required,min=1,max=255"`
	Subtitle      string    `gorm:"size:255" json:"subtitle" validate:"max=255"`
	Content       string    `gorm:"type:text" json:"content" validate:"required,min=20"`
	PublishedDate time.Time `json:"published_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AuthorID      uint      `gorm:"notnull;index" json:"author_id"`
	Author        *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
	Email     string    `gorm:"size:100;unique;not null" json:"email" validate:"required,email,max=100"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Blogs     []Blog    `gorm:"foreignKey:AuthorID" json:"blogs,omitempty"`
}
