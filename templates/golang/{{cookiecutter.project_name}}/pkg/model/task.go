package model

import (
	"time"

	"gorm.io/datatypes"
)

// Task 后台任务
type Task struct {
	BaseModel
	ID        int64          `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(128);not null"`
	Args      datatypes.JSON `json:"args" gorm:"type:json"`
	Result    datatypes.JSON `json:"result" gorm:"type:json"`
	StartedAt time.Time      `json:"startedAt" gorm:"type:datetime;default:null"`
	Duration  time.Duration  `json:"duration" gorm:"type:bigint;default:null"`
}

// PeriodicTask 周期任务
type PeriodicTask struct {
	BaseModel
	ID      int64          `json:"id" gorm:"primaryKey"`
	Cron    string         `json:"cron" gorm:"type:varchar(32);not null"`
	Name    string         `json:"name" gorm:"type:varchar(128);not null"`
	Args    datatypes.JSON `json:"args" gorm:"type:json"`
	Enabled bool           `json:"enabled" gorm:"not null;default:true"`
}
