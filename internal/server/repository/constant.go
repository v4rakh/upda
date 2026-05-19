//nolint:dupl
package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type ConstantRepository interface {
	FindAll() ([]*model.Constant, error)
	FindById(id string) (*model.Constant, error)
	FindByKey(key string) (*model.Constant, error)
	Create(key string, value string) (*model.Constant, error)
	Update(id string, value string) (*model.Constant, error)
	Delete(id string) (int64, error)
}

type ConstantDbRepo struct {
	db *gorm.DB
}

func NewConstantDbRepo(db *gorm.DB) *ConstantDbRepo {
	return &ConstantDbRepo{
		db: db,
	}
}

func (r *ConstantDbRepo) FindAll() ([]*model.Constant, error) {
	var e []*model.Constant
	var res *gorm.DB

	if res = r.db.Order("key asc").Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *ConstantDbRepo) FindById(id string) (*model.Constant, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Constant
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *ConstantDbRepo) FindByKey(key string) (*model.Constant, error) {
	if key == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Constant
	var res *gorm.DB

	if res = r.db.Find(&e, "key = ?", key); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *ConstantDbRepo) Create(key string, value string) (*model.Constant, error) {
	if key == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e = &model.Constant{
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

func (r *ConstantDbRepo) Update(id string, value string) (*model.Constant, error) {
	if id == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Constant

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

func (r *ConstantDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.Constant{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}
