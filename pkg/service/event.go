package service

import (
	"bikefest/pkg/model"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EventServiceImpl struct {
	db    *gorm.DB
	cache *redis.Client
}

func (es *EventServiceImpl) FindAll(ctx context.Context, page, limit int64) (events []*model.Event, err error) {
	err = es.db.WithContext(ctx).Limit(int(limit)).Offset(int((page - 1) * limit)).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return
}

func (es *EventServiceImpl) FindByID(ctx context.Context, id string) (event *model.Event, err error) {
	err = es.db.WithContext(ctx).Where(&model.Event{ID: &id}).First(&event).Error
	if err != nil {
		return nil, err
	}
	return
}

func (es *EventServiceImpl) Store(ctx context.Context, event *model.Event) error {
	err := es.db.WithContext(ctx).Create(event).Error
	if err != nil {
		return err
	}
	return nil
}

func (es *EventServiceImpl) Update(ctx context.Context, event *model.Event) (rowAffected int64, err error) {
	res := es.db.WithContext(ctx).Model(event).Updates(event)
	rowAffected, err = res.RowsAffected, res.Error
	if err != nil {
		return 0, err
	}
	return
}

func (es *EventServiceImpl) Delete(ctx context.Context, event *model.Event) (rowAffected int64, err error) {
	res := es.db.WithContext(ctx).Delete(event)
	rowAffected, err = res.RowsAffected, res.Error
	if err != nil {
		return 0, err
	}
	return
}

func NewEventService(db *gorm.DB, cache *redis.Client) model.EventService {
	return &EventServiceImpl{
		db:    db,
		cache: cache,
	}
}
