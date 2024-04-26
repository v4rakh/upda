package server

import (
	"encoding/json"
	"git.myservermanager.com/varakh/upda/api"
	"gorm.io/gorm"
)

type ActionRepository interface {
	paginate(page int, pageSize int, orderBy string, order string) ([]*Action, error)
	count() (int64, error)
	find(id string) (*Action, error)
	findByEnabled(enabled bool) ([]*Action, error)
	findAll() ([]*Action, error)
	create(label string, t api.ActionType, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool) (*Action, error)
	updateLabel(id string, label string) (*Action, error)
	updateMatchEvent(id string, matchEvent *string) (*Action, error)
	updateMatchApplication(id string, matchApplication *string) (*Action, error)
	updateMatchProvider(id string, matchProvider *string) (*Action, error)
	updateMatchHost(id string, matchHost *string) (*Action, error)
	updateTypeAndPayload(id string, t api.ActionType, payload interface{}) (*Action, error)
	updateEnabled(id string, enabled bool) (*Action, error)
	delete(id string) (int64, error)
}

type actionDbRepo struct {
	db *gorm.DB
}

func newActionDbRepo(db *gorm.DB) *actionDbRepo {
	return &actionDbRepo{
		db: db,
	}
}

func (r *actionDbRepo) find(id string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e Action
	var res *gorm.DB
	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *actionDbRepo) findByEnabled(enabled bool) ([]*Action, error) {
	var e []*Action

	res := r.db.Find(&e, "enabled = ?", enabled)

	if res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *actionDbRepo) create(label string, t api.ActionType, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool) (*Action, error) {
	if label == "" || t == "" {
		return nil, errorValidationNotBlank
	}

	e := &Action{
		Label:            label,
		Type:             t.Value(),
		MatchEvent:       matchEvent,
		MatchHost:        matchHost,
		MatchApplication: matchApplication,
		MatchProvider:    matchProvider,
		Enabled:          enabled,
	}

	if payload != nil {
		unmarshalledPayload := JSONMap{}
		marshalledMetadata, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		err = unmarshalledPayload.UnmarshalJSON(marshalledMetadata)
		if err != nil {
			return nil, err
		}
		e.Payload = unmarshalledPayload
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

func (r *actionDbRepo) updateLabel(id string, label string) (*Action, error) {
	if id == "" || label == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.Label = label

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) updateType(id string, t api.ActionType) (*Action, error) {
	if id == "" || t == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.Type = t.Value()

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) updateMatchEvent(id string, matchEvent *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.MatchEvent = matchEvent

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) updateMatchApplication(id string, matchApplication *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.MatchApplication = matchApplication

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) updateMatchProvider(id string, matchProvider *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.MatchProvider = matchProvider

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) updateMatchHost(id string, matchHost *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.MatchHost = matchHost

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) updateTypeAndPayload(id string, t api.ActionType, payload interface{}) (*Action, error) {
	if id == "" || t == "" {
		return nil, errorValidationNotBlank
	}
	if payload == nil {
		return nil, errorValidationNotEmpty
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	unmarshalledPayload := JSONMap{}
	marshalledMetadata, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	err = unmarshalledPayload.UnmarshalJSON(marshalledMetadata)
	if err != nil {
		return nil, err
	}
	e.Payload = unmarshalledPayload
	e.Type = t.Value()

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) updateEnabled(id string, enabled bool) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Action

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.Enabled = enabled

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionDbRepo) delete(id string) (int64, error) {
	if id == "" {
		return 0, errorValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&Action{}, "id = ?", id); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *actionDbRepo) paginate(page int, pageSize int, orderBy string, order string) ([]*Action, error) {
	if page == 0 {
		return nil, errorValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, errorValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*Action
	var res *gorm.DB

	if orderBy != "" && order != "" {
		res = r.db.Order(orderBy + " " + order).Offset(offset).Limit(pageSize).Find(&e)
	} else {
		res = r.db.Offset(offset).Limit(pageSize).Find(&e)
	}

	if res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *actionDbRepo) count() (int64, error) {
	var c int64
	var res *gorm.DB

	if res = r.db.Model(&Action{}).Count(&c); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return c, nil
}

func (r *actionDbRepo) findAll() ([]*Action, error) {
	var e []*Action

	if res := r.db.Model(&Action{}).Order("updated_at desc").Find(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}
