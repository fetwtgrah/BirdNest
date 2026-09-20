package model

import "gorm.io/gorm"

type Passage struct {
	gorm.Model
	UserID  uint
	Content string `json:"content"`
	Tags    []Tag  `gorm:"many2many:passage_tags;" json:"tags"`
}
