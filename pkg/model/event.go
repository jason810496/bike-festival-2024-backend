package model

import (
	"context"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID     string `gorm:"type:varchar(36);primary_key"`
	UserID string `gorm:"type:varchar(36);index;not null"`
	// the event id is defne at the frontend
	EventID        string  `gorm:"type:varchar(36);index;not null"`
	EventTimeStart *string `gorm:"type:varchar(255)"`
	EventTimeEnd   *string `gorm:"type:varchar(255)"`
	// the `EventDetail` field store the event detail in json format, this would be parsed when send to line message API
	EventDetail *string `gorm:"type:varchar(1024)"`
}

type EventService interface {
	FindAll(ctx context.Context, page, limit uint64) ([]Event, error)
	FindByID(ctx context.Context, id string) (*Event, error)
	FindByUserID(ctx context.Context, userID string) ([]Event, error)
	Store(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, event *Event) error
}
