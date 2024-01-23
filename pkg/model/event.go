package model

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	// the event id is defne at the frontend, if frontend don't have event id, the event id would be calculated by the hash of event detail and event time
	ID             *string `gorm:"type:varchar(36);primary_key" json:"id"`
	EventTimeStart *string `gorm:"type:varchar(255)" json:"event_time_start"`
	EventTimeEnd   *string `gorm:"type:varchar(255)" json:"event_time_end"`
	// the `EventDetail` field store the event detail in json format, this would be parsed when send to line message API
	EventDetail *string `gorm:"type:varchar(1024)" json:"event_detail"`
}

func (e *Event) BeforeCreate(*gorm.DB) error {
	if e.ID == nil {
		uuidStr := uuid.New().String()
		e.ID = &uuidStr
	}
	return nil
}

func CaculateEventID(event *Event) (string, error) {
	eventMap := make(map[string]interface{})
	eventMap["event_time_start"] = event.EventTimeStart
	eventMap["event_time_end"] = event.EventTimeEnd
	eventMap["event_detail"] = event.EventDetail

	// stringfy the event map and calculate the hash
	eventJson, err := json.Marshal(eventMap)
	if err != nil {
		return "", err
	}
	return uuid.NewSHA1(uuid.Nil, eventJson).String(), nil
}

type CreateEventRequest struct {
	ID             *string `json:"id"`
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
	FindAll(ctx context.Context, page, limit int64) ([]*Event, error)
	FindByID(ctx context.Context, id string) (*Event, error)
	Store(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) (rowAffected int64, err error)
	Delete(ctx context.Context, event *Event) (rowAffected int64, err error)
}

type AsynqNotificationService interface {
	EnqueueEvent(userId, eventId, eventStartTime string)
	// TODO: delete event notification by event id
}
