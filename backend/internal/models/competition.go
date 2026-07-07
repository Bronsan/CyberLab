package models

import "time"

// Competition represents a time-limited CTF contest.
type Competition struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (Competition) TableName() string {
	return "competitions"
}

// CompetitionChallenge links challenges to competitions.
type CompetitionChallenge struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	CompetitionID uint `gorm:"index;not null" json:"competitionId"`
	ChallengeID   uint `gorm:"index;not null" json:"challengeId"`
}

func (CompetitionChallenge) TableName() string {
	return "competition_challenges"
}

// CompetitionScore tracks per-user scores within a competition.
type CompetitionScore struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CompetitionID uint      `gorm:"index;not null" json:"competitionId"`
	UserID        uint      `gorm:"index;not null" json:"userId"`
	TeamID        uint      `gorm:"default:0" json:"teamId"`
	Score         int       `gorm:"default:0" json:"score"`
	SolvedCount   int       `gorm:"default:0" json:"solvedCount"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (CompetitionScore) TableName() string {
	return "competition_scores"
}

// Team represents a group of users.
type Team struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CaptainID   uint      `gorm:"index" json:"captainId"`
	Score       int       `gorm:"default:0" json:"score"`
	MemberCount int       `gorm:"default:0" json:"memberCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (Team) TableName() string {
	return "teams"
}

// TeamMember links users to teams.
type TeamMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TeamID    uint      `gorm:"index;not null" json:"teamId"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	Role      string    `gorm:"type:varchar(20);default:member" json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
}

func (TeamMember) TableName() string {
	return "team_members"
}
