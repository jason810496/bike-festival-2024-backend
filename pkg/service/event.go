package service

import (
	"bikefest/pkg/model"
	"context"

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

func (es *EventService) Update(ctx context.Context, event *model.Event) (rowAffected int64, err error) {
	res := es.db.WithContext(ctx).Model(event).Updates(event)
	rowAffected, err = res.RowsAffected, res.Error
	if err != nil {
		return 0, err
	}
	return
}

func (es *EventService) Delete(ctx context.Context, event *model.Event) (rowAffected int64, err error) {
	res := es.db.WithContext(ctx).Delete(event)
	rowAffected, err = res.RowsAffected, res.Error
	if err != nil {
		return 0, err
	}
	return
}

func (es *EventService) DeleteByUser(ctx context.Context, userID string, eventID string) (rowAffected int64, err error) {
	res := es.db.WithContext(ctx).Where(&model.Event{UserID: userID, EventID: eventID}).Delete(&model.Event{})
	rowAffected, err = res.RowsAffected, res.Error
	if err != nil {
		return 0, err
	}
	return
}

func NewEventService(db *gorm.DB, cache *redis.Client) model.EventService {
	return &EventService{
		db:    db,
		cache: cache,
	}
}
