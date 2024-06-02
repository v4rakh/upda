package server

import (
	"encoding/json"
	"git.myservermanager.com/varakh/upda/api"
	"gorm.io/gorm"
	"time"
)

type eventRepository interface {
	find(id string) (*Event, error)
	window(size int, skip int, orderBy string, order string) ([]*Event, error)
	windowHasNext(size int, skip int, orderBy string, order string) (bool, error)
	count(state ...api.EventState) (int64, error)
	findAllByState(limit int, state ...api.EventState) ([]*Event, error)
	create(name api.EventName, state api.EventState, payload interface{}) (*Event, error)
	updateState(id string, state api.EventState) (*Event, error)
	delete(id string) (int64, error)
	deleteByUpdatedAtBeforeAndStates(time time.Time, state ...api.EventState) (int64, error)
}

type eventDbRepo struct {
	db *gorm.DB
}

func newEventDbRepo(db *gorm.DB) *eventDbRepo {
	return &eventDbRepo{
		db: db,
	}
}

func (r *eventDbRepo) find(id string) (*Event, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e Event
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *eventDbRepo) create(name api.EventName, state api.EventState, payload interface{}) (*Event, error) {
	if name == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	var e *Event
	unmarshalledPayload := JSONMap{}

	if payload != nil {
		marshalledMetadata, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		err = unmarshalledPayload.UnmarshalJSON(marshalledMetadata)
		if err != nil {
			return nil, err
		}
	}

	e = &Event{
		Name:    name.Value(),
		State:   state.Value(),
		Payload: unmarshalledPayload,
	}

	var res *gorm.DB
	if res = r.db.Create(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *eventDbRepo) updateState(id string, state api.EventState) (*Event, error) {
	if id == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Event

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.State = state.Value()

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *eventDbRepo) delete(id string) (int64, error) {
	if id == "" {
		return 0, errorValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&Event{}, "id = ?", id); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}

func (r *eventDbRepo) deleteByUpdatedAtBeforeAndStates(time time.Time, state ...api.EventState) (int64, error) {
	if len(state) == 0 {
		return 0, errorValidationNotEmpty
	}

	states := make([]string, 0, len(state))
	for _, i := range state {
		states = append(states, i.Value())
	}

	var res *gorm.DB
	if res = r.db.Where("state IN ?", states).Where("updated_at < ?", time).Delete(&Event{}); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *eventDbRepo) window(size int, skip int, orderBy string, order string) ([]*Event, error) {
	if size <= 0 {
		return nil, errorValidationSizeGreaterZero
	}
	if orderBy == "" {
		orderBy = "created_at"
	}
	if order == "" {
		order = "asc"
	}

	var e []*Event
	if res := r.db.Order(orderBy + " " + order).Offset(skip).Limit(size).Find(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *eventDbRepo) windowHasNext(size int, skip int, orderBy string, order string) (bool, error) {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if order == "" {
		order = "asc"
	}

	var e []*Event

	if res := r.db.Order(orderBy + " " + order).Offset(skip + size).Find(&e); res.Error != nil {
		return false, newServiceDatabaseError(res.Error)
	}

	return len(e) > 0, nil
}

func (r *eventDbRepo) findAllByState(limit int, state ...api.EventState) ([]*Event, error) {
	if len(state) == 0 {
		return nil, errorValidationNotEmpty
	}
	if limit <= 0 {
		return nil, errorValidationLimitGreaterZero
	}

	var e []*Event

	states := translateEventState(state...)

	if res := r.db.Model(&Event{}).Scopes(allGetEventCriterion(states)).Order("created_at asc").Limit(limit).Find(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *eventDbRepo) count(state ...api.EventState) (int64, error) {
	var c int64

	states := translateEventState(state...)
	if res := r.db.Model(&Event{}).Scopes(allGetEventCriterion(states)).Count(&c); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return c, nil
}

func translateEventState(state ...api.EventState) []string {
	states := make([]string, 0)
	if len(state) > 0 {
		for _, s := range state {
			states = append(states, s.Value())
		}
	}

	return states
}

func criterionEventState(states []string) func(db *gorm.DB) *gorm.DB {
	if states != nil && len(states) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("state IN (?)", states)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func allGetEventCriterion(states []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(criterionEventState(states))
	}
}
