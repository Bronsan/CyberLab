package models

import "time"

type Submission struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index:idx_user_submit;not null" json:"userId"`
	ChallengeID   uint      `gorm:"index:idx_challenge_submit;not null" json:"challengeId"`
	SubmittedFlag string    `gorm:"type:varchar(255)" json:"submittedFlag"`
	IsCorrect     bool      `json:"isCorrect"`
	SubmittedAt   time.Time `json:"submittedAt"`

	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Challenge *Challenge `gorm:"foreignKey:ChallengeID" json:"challenge,omitempty"`
}

func (Submission) TableName() string {
	return "submissions"
}

type UserChallenge struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"userId"`
	ChallengeID uint      `gorm:"index;not null" json:"challengeId"`
	Score       int       `json:"score"`
	SolvedAt    time.Time `json:"solvedAt"`

	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Challenge *Challenge `gorm:"foreignKey:ChallengeID" json:"challenge,omitempty"`
}

func (UserChallenge) TableName() string {
	return "user_challenges"
}
