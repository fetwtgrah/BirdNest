package model

type Passage struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
	//CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time `json:"updated_at"`
	User User `gorm:"foreignKey:UserID" json:"-"`
}
