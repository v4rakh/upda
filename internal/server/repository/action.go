package repository

import (
	"encoding/json"

	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type ActionRepository interface {
	Paginate(page int, pageSize int, orderBy string, order string) ([]*model.Action, error)
	Count() (int64, error)
	Find(id string) (*model.Action, error)
	FindByEnabled(enabled bool) ([]*model.Action, error)
	FindAll() ([]*model.Action, error)
	Create(label string, t string, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool) (*model.Action, error)
	UpdateLabel(id string, label string) (*model.Action, error)
	UpdateMatchEvent(id string, matchEvent *string) (*model.Action, error)
	UpdateMatchApplication(id string, matchApplication *string) (*model.Action, error)
	UpdateMatchProvider(id string, matchProvider *string) (*model.Action, error)
	UpdateMatchHost(id string, matchHost *string) (*model.Action, error)
	UpdateTypeAndPayload(id string, t string, payload interface{}) (*model.Action, error)
	UpdateEnabled(id string, enabled bool) (*model.Action, error)
	Delete(id string) (int64, error)
}

type ActionDbRepo struct {
	db *gorm.DB
}

func NewActionDbRepo(db *gorm.DB) *ActionDbRepo {
	return &ActionDbRepo{
		db: db,
	}
}

func (r *ActionDbRepo) Find(id string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Action
	var res *gorm.DB
	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *ActionDbRepo) FindByEnabled(enabled bool) ([]*model.Action, error) {
	var e []*model.Action

	res := r.db.Find(&e, "enabled = ?", enabled)

	if res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *ActionDbRepo) Create(label string, t string, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool) (*model.Action, error) {
	if label == "" || t == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e := &model.Action{
		Label:            label,
		Type:             t,
		MatchEvent:       matchEvent,
		MatchHost:        matchHost,
		MatchApplication: matchApplication,
		MatchProvider:    matchProvider,
		Enabled:          enabled,
	}

	if payload != nil {
		bytes, err := json.Marshal(payload)
		if err != nil {
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

func (r *ActionDbRepo) UpdateLabel(id string, label string) (*model.Action, error) {
	if id == "" || label == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.Label = label

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) UpdateType(id string, t string) (*model.Action, error) {
	if id == "" || t == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.Type = t

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) UpdateMatchEvent(id string, matchEvent *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.MatchEvent = matchEvent

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) UpdateMatchApplication(id string, matchApplication *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.MatchApplication = matchApplication

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) UpdateMatchProvider(id string, matchProvider *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.MatchProvider = matchProvider

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) UpdateMatchHost(id string, matchHost *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.MatchHost = matchHost

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) UpdateTypeAndPayload(id string, t string, payload interface{}) (*model.Action, error) {
	if id == "" || t == "" {
		return nil, service_error.ErrValidationNotBlank
	}
	if payload == nil {
		return nil, service_error.ErrValidationNotEmpty
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.Type = t

	var b []byte
	if b, err = json.Marshal(payload); err != nil {
		return nil, err
	}
	e.Payload = b

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) UpdateEnabled(id string, enabled bool) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Action

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.Enabled = enabled

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.Action{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *ActionDbRepo) Paginate(page int, pageSize int, orderBy string, order string) ([]*model.Action, error) {
	if page == 0 {
		return nil, service_error.ErrValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, service_error.ErrValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*model.Action
	var res *gorm.DB

	if orderBy != "" && order != "" {
		res = r.db.Order(orderBy + " " + order).Offset(offset).Limit(pageSize).Find(&e)
	} else {
		res = r.db.Order("label asc").Offset(offset).Limit(pageSize).Find(&e)
	}

	if res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *ActionDbRepo) Count() (int64, error) {
	var c int64
	var res *gorm.DB

	if res = r.db.Model(&model.Action{}).Count(&c); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return c, nil
}

func (r *ActionDbRepo) FindAll() ([]*model.Action, error) {
	var e []*model.Action

	if res := r.db.Model(&model.Action{}).Order("updated_at desc").Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}
