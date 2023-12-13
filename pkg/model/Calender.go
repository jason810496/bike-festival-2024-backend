package model

import "gorm.io/gorm"

type Calender struct {
	gorm.Model
	Id       int    `gorm:"type:int;primary_key"`
	GoogleId string `gorm:"type:varchar(255)"`
	EventId  int    `gorm:"type:int"`
}
