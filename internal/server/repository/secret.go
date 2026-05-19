//nolint:dupl
package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type SecretRepository interface {
	FindAll() ([]*model.Secret, error)
	FindById(id string) (*model.Secret, error)
	FindByKey(key string) (*model.Secret, error)
	Create(key string, value string) (*model.Secret, error)
	Update(id string, value string) (*model.Secret, error)
	Delete(id string) (int64, error)
}

type SecretDbRepo struct {
	db *gorm.DB
}

func NewSecretDbRepo(db *gorm.DB) *SecretDbRepo {
	return &SecretDbRepo{
		db: db,
	}
}

func (r *SecretDbRepo) FindAll() ([]*model.Secret, error) {
	var e []*model.Secret
	var res *gorm.DB

	if res = r.db.Order("key asc").Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *SecretDbRepo) FindById(id string) (*model.Secret, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Secret
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *SecretDbRepo) FindByKey(key string) (*model.Secret, error) {
	if key == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Secret
	var res *gorm.DB

	if res = r.db.Find(&e, "key = ?", key); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *SecretDbRepo) Create(key string, value string) (*model.Secret, error) {
	if key == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e = &model.Secret{
		Key:   key,
		Value: value,
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

func (r *SecretDbRepo) Update(id string, value string) (*model.Secret, error) {
	if id == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Secret

	if e, err = r.FindById(id); err != nil {
		return nil, err
	}

	e.Value = value

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *SecretDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.Secret{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}
