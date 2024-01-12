package service

import (
	"bikefest/pkg/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EventService struct {
	db    *gorm.DB
	cache *redis.Client
}

func (es *EventService) FindAll(ctx context.Context, page, limit uint64) (events []model.Event, err error) {
	err = es.db.WithContext(ctx).Limit(int(limit)).Offset(int((page - 1) * limit)).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return
}

func (es *EventService) FindByID(ctx context.Context, id string) (event *model.Event, err error) {
	err = es.db.WithContext(ctx).Where(&model.Event{ID: id}).First(&event).Error
	if err != nil {
		return nil, err
	}
	return
}

func (es *EventService) FindByUserID(ctx context.Context, userID string) (events []model.Event, err error) {
	err = es.db.WithContext(ctx).Where(&model.Event{UserID: userID}).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return
}

func (es *EventService) Store(ctx context.Context, event *model.Event) error {
	err := es.db.WithContext(ctx).Create(event).Error
	if err != nil {
		return err
	}
	// TODO: add event to distributed scheduler
	return nil
}

func (es *EventService) Update(ctx context.Context, event *model.Event) error {
	err := es.db.WithContext(ctx).Updates(event).Error
	if err != nil {
		return err
	}
	// TODO: update event time to distributed scheduler
	return nil
}

func (es *EventService) Delete(ctx context.Context, event *model.Event) error {
	err := es.db.WithContext(ctx).Delete(event).Error
	if err != nil {
		return err
	}
	//TODO: delete event from distributed scheduler
	return nil
}

func NewEventService(db *gorm.DB, cache *redis.Client) model.EventService {
	return &EventService{
		db:    db,
		cache: cache,
	}
}
