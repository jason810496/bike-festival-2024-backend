package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID   string `json:"id" gorm:"type:varchar(36);primary_key"`
	Name string `json:"name" gorm:"type:varchar(255);index"`
}

func (u *User) BeforeCreate(*gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
