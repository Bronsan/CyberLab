package models

import "time"

type Challenge struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Title          string    `gorm:"type:varchar(100);not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	Category       string    `gorm:"type:varchar(50);index" json:"category"`
	Difficulty     string    `gorm:"type:varchar(20);index" json:"difficulty"`
	Score          int       `gorm:"default:0" json:"score"`
	DockerImage    string    `gorm:"type:varchar(255)" json:"dockerImage"`
	Flag           string    `gorm:"type:varchar(255)" json:"-"`
	TimeoutMinutes int       `gorm:"default:30" json:"timeoutMinutes"`
	IsActive       bool      `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Tags           []ChallengeTag `gorm:"foreignKey:ChallengeID" json:"tags,omitempty"`
}

func (Challenge) TableName() string {
	return "challenges"
}

type ChallengeTag struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ChallengeID uint   `gorm:"index;not null" json:"challengeId"`
	TagName     string `gorm:"type:varchar(50)" json:"tagName"`
}

func (ChallengeTag) TableName() string {
	return "challenge_tags"
}
