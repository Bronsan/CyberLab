package models

import "time"

type ChallengeInstance struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index:idx_user;not null" json:"userId"`
	ChallengeID   uint      `gorm:"index;not null" json:"challengeId"`
	ContainerID   string    `gorm:"type:varchar(255)" json:"containerId"`
	ContainerName string    `gorm:"type:varchar(255)" json:"containerName"`
	HostPort      int       `json:"hostPort"`
	Status        string    `gorm:"type:varchar(20);index:idx_status;default:CREATED" json:"status"`
	StartTime     time.Time `json:"startTime"`
	ExpireTime    time.Time `gorm:"index:idx_expire" json:"expireTime"`

	// Relations
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Challenge *Challenge `gorm:"foreignKey:ChallengeID" json:"challenge,omitempty"`
}

func (ChallengeInstance) TableName() string {
	return "challenge_instances"
}

// Status constants
const (
	InstanceStatusCreated   = "CREATED"
	InstanceStatusRunning   = "RUNNING"
	InstanceStatusStopped   = "STOPPED"
	InstanceStatusDestroyed = "DESTROYED"
	InstanceStatusError     = "ERROR"
)
