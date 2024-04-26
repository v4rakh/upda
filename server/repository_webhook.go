package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"gorm.io/gorm"
)

type WebhookRepository interface {
	paginate(page int, pageSize int, orderBy string, order string) ([]*Webhook, error)
	count() (int64, error)
	find(id string) (*Webhook, error)
	create(label string, t api.WebhookType, token string, ignoreHost bool) (*Webhook, error)
	updateLabel(id string, label string) (*Webhook, error)
	updateIgnoreHost(id string, ignoreHost bool) (*Webhook, error)
	delete(id string) (int64, error)
}

type webhookDbRepo struct {
	db *gorm.DB
}

func newWebhookDbRepo(db *gorm.DB) *webhookDbRepo {
	return &webhookDbRepo{
		db: db,
	}
}

func (r *webhookDbRepo) find(id string) (*Webhook, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e Webhook
	var res *gorm.DB
	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *webhookDbRepo) create(label string, t api.WebhookType, token string, ignoreHost bool) (*Webhook, error) {
	if label == "" || t == "" || token == "" {
		return nil, errorValidationNotBlank
	}

	e := &Webhook{
		Label:      label,
		Type:       t.Value(),
		Token:      token,
		IgnoreHost: ignoreHost,
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

func (r *webhookDbRepo) updateLabel(id string, label string) (*Webhook, error) {
	if id == "" || label == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Webhook

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

func (r *webhookDbRepo) updateIgnoreHost(id string, ignoreHost bool) (*Webhook, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Webhook

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.IgnoreHost = ignoreHost

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *webhookDbRepo) delete(id string) (int64, error) {
	if id == "" {
		return 0, errorValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&Webhook{}, "id = ?", id); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *webhookDbRepo) paginate(page int, pageSize int, orderBy string, order string) ([]*Webhook, error) {
	if page == 0 {
		return nil, errorValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, errorValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*Webhook
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

func (r *webhookDbRepo) count() (int64, error) {
	var c int64
	var res *gorm.DB

	if res = r.db.Model(&Webhook{}).Count(&c); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return c, nil
}
