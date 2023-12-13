package model

import "gorm.io/gorm"

type PsychoTest struct {
	gorm.Model
	Id    int    `gorm:"type:int;primary_key"`
	Type  string `gorm:"type:varchar(255);unique"`
	Count int    `gorm:"type:int"`
}
