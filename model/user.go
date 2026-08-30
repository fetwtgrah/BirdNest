package model

type User struct {
	ID       int64     `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Passages []Passage `gorm:"foreignKey:UserID" json:"passages"`
}
