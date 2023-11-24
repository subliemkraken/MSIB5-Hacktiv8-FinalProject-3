package model

import (
	"time"
)

type Category struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Type      string `gorm:"type:varchar(200)" json:"type" valid:"required~Type cannot be empty"`
	Task      []Task
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Category) TableName() string {
	return "category"
}
