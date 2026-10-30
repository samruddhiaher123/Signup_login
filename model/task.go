package model

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Subject     string    `gorm:"type:text;not null" json:"subject"`
	Topic       string    `gorm:"type:text;not null" json:"topic"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Date        string    `gorm:"type:text;not null" json:"date"`
	Duration    int64     `gorm:"not null" json:"duration"`
	Priority    string    `gorm:"type:text;not null" json:"priority"`
	Status      string    `gorm:"type:text;not null" json:"status"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}
