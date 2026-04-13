package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type UpdateStateDefinitionRepository interface {
	FindAll() ([]*model.UpdateStateDefinition, error)
	Find(id string) (*model.UpdateStateDefinition, error)
	FindByName(name string) (*model.UpdateStateDefinition, error)
	FindInitial() (*model.UpdateStateDefinition, error)
	ExistsByName(name string) (bool, error)
	MaxSortOrder() (int, error)
	Create(name string, label string, color string, icon string, description *string, isInitial bool, skipOnNewVersion bool, sortOrder int) (*model.UpdateStateDefinition, error)
	Update(id string, name string, label string, color string, icon string, description *string, isInitial bool, skipOnNewVersion bool, sortOrder int) (*model.UpdateStateDefinition, error)
	Reorder(items []dto.UpdateStateReorderItem) error
	Delete(id string) (int64, error)
	ClearInitial() error
}

type UpdateStateDefinitionDbRepo struct {
	db *gorm.DB
}

func NewUpdateStateDefinitionDbRepo(db *gorm.DB) *UpdateStateDefinitionDbRepo {
	return &UpdateStateDefinitionDbRepo{
		db: db,
	}
}

func (r *UpdateStateDefinitionDbRepo) FindAll() ([]*model.UpdateStateDefinition, error) {
	var e []*model.UpdateStateDefinition
	var res *gorm.DB

	if res = r.db.Order("sort_order asc, name asc").Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *UpdateStateDefinitionDbRepo) Find(id string) (*model.UpdateStateDefinition, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.UpdateStateDefinition
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *UpdateStateDefinitionDbRepo) FindByName(name string) (*model.UpdateStateDefinition, error) {
	if name == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.UpdateStateDefinition
	var res *gorm.DB

	if res = r.db.Find(&e, "name = ?", name); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *UpdateStateDefinitionDbRepo) FindInitial() (*model.UpdateStateDefinition, error) {
	var e model.UpdateStateDefinition
	var res *gorm.DB

	if res = r.db.Find(&e, "is_initial = ?", true); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *UpdateStateDefinitionDbRepo) ExistsByName(name string) (bool, error) {
	if name == "" {
		return false, service_error.ErrValidationNotBlank
	}

	var count int64
	if res := r.db.Model(&model.UpdateStateDefinition{}).Where("name = ?", name).Count(&count); res.Error != nil {
		return false, service_error.NewServiceDatabaseError(res.Error)
	}

	return count > 0, nil
}

func (r *UpdateStateDefinitionDbRepo) Create(name string, label string, color string, icon string, description *string, isInitial bool, skipOnNewVersion bool, sortOrder int) (*model.UpdateStateDefinition, error) {
	if name == "" || label == "" || color == "" || icon == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e := &model.UpdateStateDefinition{
		Name:             name,
		Label:            label,
		Color:            color,
		Icon:             icon,
		Description:      description,
		IsInitial:        isInitial,
		SkipOnNewVersion: skipOnNewVersion,
		SortOrder:        sortOrder,
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

func (r *UpdateStateDefinitionDbRepo) Update(id string, name string, label string, color string, icon string, description *string, isInitial bool, skipOnNewVersion bool, sortOrder int) (*model.UpdateStateDefinition, error) {
	if id == "" || name == "" || label == "" || color == "" || icon == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.UpdateStateDefinition

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.Name = name
	e.Label = label
	e.Color = color
	e.Icon = icon
	e.Description = description
	e.IsInitial = isInitial
	e.SkipOnNewVersion = skipOnNewVersion
	e.SortOrder = sortOrder

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *UpdateStateDefinitionDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.UpdateStateDefinition{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}

func (r *UpdateStateDefinitionDbRepo) ClearInitial() error {
	if res := r.db.Model(&model.UpdateStateDefinition{}).Where("is_initial = ?", true).Update("is_initial", false); res.Error != nil {
		return service_error.NewServiceDatabaseError(res.Error)
	}
	return nil
}

func (r *UpdateStateDefinitionDbRepo) MaxSortOrder() (int, error) {
	var maxSortOrder *int
	if res := r.db.Model(&model.UpdateStateDefinition{}).Select("MAX(sort_order)").Scan(&maxSortOrder); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	if maxSortOrder == nil {
		return -1, nil
	}
	return *maxSortOrder, nil
}

func (r *UpdateStateDefinitionDbRepo) Reorder(items []dto.UpdateStateReorderItem) error {
	if len(items) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if res := tx.Model(&model.UpdateStateDefinition{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder); res.Error != nil {
				return service_error.NewServiceDatabaseError(res.Error)
			}
		}
		return nil
	})
}
