package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	TagName  string    `json:"tagName"`
	Passages []Passage `gorm:"many2many:passage_tags;" json:"passages"`
}
