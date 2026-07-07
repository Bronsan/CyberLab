package models

import "time"

type AIHint struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"userId"`
	ChallengeID uint      `gorm:"index;not null" json:"challengeId"`
	Question    string    `gorm:"type:text" json:"question"`
	Answer      string    `gorm:"type:text" json:"answer"`
	CreatedAt   time.Time `json:"createdAt"`

	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Challenge *Challenge `gorm:"foreignKey:ChallengeID" json:"challenge,omitempty"`
}

func (AIHint) TableName() string {
	return "ai_hints"
}
