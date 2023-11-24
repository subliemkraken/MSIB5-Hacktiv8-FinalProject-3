package model

import (
	"time"
)

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:varchar(100)" json:"description"`
	Status      bool   `json:"status"`
	UserID      uint   `json:"user_id"`
	CategoryID  uint   `json:"category_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        User     `gorm:"foreignKey:UserID"`
	Category    Category `gorm:"foreignKey:CategoryID"`
}

type NewTaskRequest struct {
	Title       string `gorm:"type:varchar(50)" json:"title" valid:"required~Title cannot be empty"`
	Description string `gorm:"type:varchar(100)" json:"description" valid:"required~Description cannot be empty"`
	CategoryID  uint   `json:"category_id"`
}

type PutRequest struct {
	Title       string `gorm:"type:varchar(50)" json:"title" valid:"required~Title cannot be empty"`
	Description string `gorm:"type:varchar(100)" json:"description" valid:"required~Description cannot be empty"`
}

type PatchStatusRequest struct {
	Status bool `json:"status" valid:"required~Status cannot be empty"`
}

type PatchCatIdRequest struct {
	CategoryID uint `json:"category_id"`
}

func (Task) TableName() string {
	return "task"
}
