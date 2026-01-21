package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type WebhookRepository interface {
	Paginate(page int, pageSize int, orderBy string, order string) ([]*model.Webhook, error)
	Count() (int64, error)
	Find(id string) (*model.Webhook, error)
	Create(label string, t string, token string, ignoreHost bool, ignoreHostReplacement string) (*model.Webhook, error)
	UpdateLabel(id string, label string) (*model.Webhook, error)
	UpdateIgnoreHost(id string, ignoreHost bool) (*model.Webhook, error)
	UpdateIgnoreHostReplacement(id string, ignoreHostReplacement string) (*model.Webhook, error)
	Delete(id string) (int64, error)
}

type WebhookDbRepo struct {
	db *gorm.DB
}

func NewWebhookDbRepo(db *gorm.DB) *WebhookDbRepo {
	return &WebhookDbRepo{
		db: db,
	}
}

func (r *WebhookDbRepo) Find(id string) (*model.Webhook, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Webhook
	var res *gorm.DB
	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *WebhookDbRepo) Create(label string, t string, token string, ignoreHost bool, ignoreHostReplacement string) (*model.Webhook, error) {
	if label == "" || t == "" || token == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e := &model.Webhook{
		Label:                 label,
		Type:                  t,
		Token:                 token,
		IgnoreHost:            ignoreHost,
		IgnoreHostReplacement: ignoreHostReplacement,
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

func (r *WebhookDbRepo) UpdateLabel(id string, label string) (*model.Webhook, error) {
	if id == "" || label == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Webhook

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

func (r *WebhookDbRepo) UpdateIgnoreHost(id string, ignoreHost bool) (*model.Webhook, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Webhook

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.IgnoreHost = ignoreHost

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *WebhookDbRepo) UpdateIgnoreHostReplacement(id string, ignoreHostReplacement string) (*model.Webhook, error) {
	if id == "" || ignoreHostReplacement == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Webhook

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.IgnoreHostReplacement = ignoreHostReplacement

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *WebhookDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.Webhook{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *WebhookDbRepo) Paginate(page int, pageSize int, orderBy string, order string) ([]*model.Webhook, error) {
	if page == 0 {
		return nil, service_error.ErrValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, service_error.ErrValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*model.Webhook
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

func (r *WebhookDbRepo) Count() (int64, error) {
	var c int64
	var res *gorm.DB

	if res = r.db.Model(&model.Webhook{}).Count(&c); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return c, nil
}
