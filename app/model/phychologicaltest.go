package model

type PhychoTest struct {
	Id    int    `gorm:"type:int;primary_key"`
	Type  string `gorm:"type:varchar(255)"`
	Count int    `gorm:"type:int"`
}
