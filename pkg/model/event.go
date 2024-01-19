package model

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID     string `gorm:"type:varchar(36);primary_key" json:"id"`
	UserID string `gorm:"type:varchar(36);index;not null" json:"user_id"`
	// the event id is defne at the frontend
	EventID        string  `gorm:"type:varchar(36);index;not null" json:"event_id" binding:"required"`
	EventTimeStart *string `gorm:"type:varchar(255)" json:"event_time_start"`
	EventTimeEnd   *string `gorm:"type:varchar(255)" json:"event_time_end"`
	// the `EventDetail` field store the event detail in json format, this would be parsed when send to line message API
	EventDetail *string `gorm:"type:varchar(1024)" json:"event_detail"`
}

func (e *Event) BeforeCreate(*gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

type CreateEventRequest struct {
	EventID        string  `json:"event_id" binding:"required"`
	EventTimeStart *string `json:"event_time_start"`
	EventTimeEnd   *string `json:"event_time_end"`
	EventDetail    *string `json:"event_detail"`
}

type EventResponse struct {
	Msg  string `json:"msg"`
	Data *Event `json:"data"`
}

type EventListResponse struct {
	Msg  string   `json:"msg"`
	Data []*Event `json:"data"`
}

type EventService interface {
	FindAll(ctx context.Context, page, limit uint64) ([]*Event, error)
	FindByID(ctx context.Context, id string) (*Event, error)
	FindByUserID(ctx context.Context, userID string) ([]*Event, error)
	Store(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) (rowAffected int64, err error)
	Delete(ctx context.Context, event *Event) (rowAffected int64, err error)
	DeleteByUser(ctx context.Context, userID string, eventID string) (rowAffected int64, err error)
}

type AsynqNotificationService interface {
	EnqueueEvent(user_id, event_id, event_start_time string)
}
