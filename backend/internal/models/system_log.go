package models

import "time"

type SystemLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Action    string    `gorm:"type:varchar(100);not null" json:"action"`
	IP        string    `gorm:"type:varchar(50)" json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
}

func (SystemLog) TableName() string {
	return "system_logs"
}

// Action constants
const (
	ActionLogin          = "login"
	ActionRegister       = "register"
	ActionStartContainer = "start_container"
	ActionStopContainer  = "stop_container"
	ActionSubmitFlag     = "submit_flag"
	ActionAIRequest      = "ai_request"
	ActionAdminOp        = "admin_operation"
)
