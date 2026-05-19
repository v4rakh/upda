package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type EventRepository interface {
	Find(id string) (*model.Event, error)
	Window(size int, skip int, orderBy string, order string, updateId *string) ([]*model.Event, error)
	WindowHasNext(size int, skip int, orderBy string, order string, updateId *string) (bool, error)
	Count(state ...string) (int64, error)
	FindAllByState(limit int, state ...string) ([]*model.Event, error)
	Create(name string, state string, payload interface{}) (*model.Event, error)
	UpdateState(id string, state string) (*model.Event, error)
	Delete(id string) (int64, error)
	DeleteByUpdatedAtBeforeAndStates(time time.Time, state ...string) (int64, error)
}

type EventDbRepo struct {
	db *gorm.DB
}

func NewEventDbRepo(db *gorm.DB) *EventDbRepo {
	return &EventDbRepo{
		db: db,
	}
}

func (r *EventDbRepo) Find(id string) (*model.Event, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Event
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *EventDbRepo) Create(name string, state string, payload interface{}) (*model.Event, error) {
	if name == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e = &model.Event{
		Name:  name,
		State: state,
	}

	if payload != nil {
		var bytes []byte
		if bytes, err = json.Marshal(payload); err != nil {
			return nil, err
		}

		e.Payload = bytes
	}

	var res *gorm.DB
	if res = r.db.Create(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *EventDbRepo) UpdateState(id string, state string) (*model.Event, error) {
	if id == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Event

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.State = state

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *EventDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.Event{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}

func (r *EventDbRepo) DeleteByUpdatedAtBeforeAndStates(time time.Time, state ...string) (int64, error) {
	if len(state) == 0 {
		return 0, service_error.ErrValidationNotEmpty
	}

	var res *gorm.DB
	if res = r.db.Where("state IN ?", state).Where("updated_at < ?", time).Delete(&model.Event{}); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *EventDbRepo) Window(size int, skip int, orderBy string, order string, updateId *string) ([]*model.Event, error) {
	if size <= 0 {
		return nil, service_error.ErrValidationSizeGreaterZero
	}
	if orderBy == "" {
		orderBy = "created_at"
	}
	if order == "" {
		order = "asc"
	}

	var e []*model.Event
	if res := r.db.Model(&model.Event{}).
		Scopes(CriterionEventUpdateID(updateId)).
		Order(orderBy + " " + order).
		Offset(skip).
		Limit(size).
		Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *EventDbRepo) WindowHasNext(size int, skip int, orderBy string, order string, updateId *string) (bool, error) {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if order == "" {
		order = "asc"
	}

	var e []*model.Event

	if res := r.db.Model(&model.Event{}).Scopes(CriterionEventUpdateID(updateId)).Order(orderBy + " " + order).Offset(skip + size).Find(&e); res.Error != nil {
		return false, service_error.NewServiceDatabaseError(res.Error)
	}

	return len(e) > 0, nil
}

func (r *EventDbRepo) FindAllByState(limit int, state ...string) ([]*model.Event, error) {
	if len(state) == 0 {
		return nil, service_error.ErrValidationNotEmpty
	}
	if limit <= 0 {
		return nil, service_error.ErrValidationLimitGreaterZero
	}

	var e []*model.Event

	if res := r.db.Model(&model.Event{}).Scopes(CriterionEventState(state)).Order("created_at asc").Limit(limit).Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *EventDbRepo) Count(state ...string) (int64, error) {
	var c int64

	if res := r.db.Model(&model.Event{}).Scopes(CriterionEventState(state)).Count(&c); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return c, nil
}

func CriterionEventState(states []string) func(db *gorm.DB) *gorm.DB {
	if states != nil && len(states) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("state IN (?)", states)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func CriterionEventUpdateID(updateId *string) func(db *gorm.DB) *gorm.DB {
	if updateId == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	names := constant.EventNameNames()
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name IN (?)", names).Where("payload @> ?", fmt.Sprintf(`{"id": "%s"}`, *updateId)).
			Or("name IN (?)", names).Where("payload @> ?", fmt.Sprintf(`{"updateId": "%s"}`, *updateId))
	}
}
