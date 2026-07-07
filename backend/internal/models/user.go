package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Avatar       string    `gorm:"type:varchar(255)" json:"avatar"`
	Bio          string    `gorm:"type:text" json:"bio"`
	Role         string    `gorm:"type:varchar(20);default:user" json:"role"`
	Score        int       `gorm:"default:0" json:"score"`
	SolvedCount  int       `gorm:"default:0" json:"solvedCount"`
	Status       int       `gorm:"default:1" json:"status"` // 1=active, 0=banned
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
