package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type FilterPresetRepository interface {
	FindByType(t string) ([]*model.FilterPreset, error)
	FindById(id string) (*model.FilterPreset, error)
	Create(t string, label string, parameters string, color *string) (*model.FilterPreset, error)
	Delete(id string) (int64, error)
}

type FilterPresetDbRepo struct {
	db *gorm.DB
}

func NewFilterPresetDbRepo(db *gorm.DB) *FilterPresetDbRepo {
	return &FilterPresetDbRepo{
		db: db,
	}
}

func (r *FilterPresetDbRepo) FindByType(t string) ([]*model.FilterPreset, error) {
	if t == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e []*model.FilterPreset
	var res *gorm.DB

	if res = r.db.Order("label asc").Find(&e, "type = ?", t); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *FilterPresetDbRepo) FindById(id string) (*model.FilterPreset, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.FilterPreset
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *FilterPresetDbRepo) Create(t string, label string, parameters string, color *string) (*model.FilterPreset, error) {
	if t == "" || label == "" || parameters == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e = &model.FilterPreset{
		Type:       t,
		Label:      label,
		Parameters: parameters,
		Color:      color,
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

func (r *FilterPresetDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.FilterPreset{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}
