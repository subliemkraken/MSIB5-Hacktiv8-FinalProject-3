package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Full_name string `gorm:"type:varchar(200)" valid:"required~Full name cannot be empty"`
	Email     string `gorm:"type:varchar(100);unique" valid:"required~Email cannot be empty,email~Invalid format for email"`
	Password  string `gorm:"type:varchar(200)" valid:"required~Password cannot be empty,minstringlength(6)~Password must be more than 6 words"`
	Role      string `gorm:"type:varchar(10)" valid:"required~Role cannot be empty"`
	Tasks     []Task
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "user"
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
