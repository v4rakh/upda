package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type UpdateStateTransitionRepository interface {
	FindAll() ([]*model.UpdateStateTransition, error)
	Find(id string) (*model.UpdateStateTransition, error)
	FindByFromStateId(fromStateId string) ([]*model.UpdateStateTransition, error)
	FindByFromStateName(fromStateName string) ([]*model.UpdateStateTransition, error)
	IsTransitionAllowed(fromStateName string, toStateName string) (bool, error)
	Exists(fromStateId string, toStateId string) (bool, error)
	Create(fromStateId string, toStateId string) (*model.UpdateStateTransition, error)
	Delete(id string) (int64, error)
}

type UpdateStateTransitionDbRepo struct {
	db *gorm.DB
}

func NewUpdateStateTransitionDbRepo(db *gorm.DB) *UpdateStateTransitionDbRepo {
	return &UpdateStateTransitionDbRepo{
		db: db,
	}
}

func (r *UpdateStateTransitionDbRepo) FindAll() ([]*model.UpdateStateTransition, error) {
	var e []*model.UpdateStateTransition
	var res *gorm.DB

	if res = r.db.Preload("FromState").Preload("ToState").Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *UpdateStateTransitionDbRepo) Find(id string) (*model.UpdateStateTransition, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.UpdateStateTransition
	var res *gorm.DB

	if res = r.db.Preload("FromState").Preload("ToState").Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *UpdateStateTransitionDbRepo) FindByFromStateId(fromStateId string) ([]*model.UpdateStateTransition, error) {
	if fromStateId == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e []*model.UpdateStateTransition
	var res *gorm.DB

	if res = r.db.Preload("FromState").Preload("ToState").Where("from_state_id = ?", fromStateId).Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *UpdateStateTransitionDbRepo) FindByFromStateName(fromStateName string) ([]*model.UpdateStateTransition, error) {
	if fromStateName == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e []*model.UpdateStateTransition
	var res *gorm.DB

	if res = r.db.Preload("FromState").Preload("ToState").
		Joins("JOIN update_state_definitions ON update_state_definitions.id = update_state_transitions.from_state_id").
		Where("update_state_definitions.name = ?", fromStateName).
		Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *UpdateStateTransitionDbRepo) IsTransitionAllowed(fromStateName string, toStateName string) (bool, error) {
	if fromStateName == "" || toStateName == "" {
		return false, service_error.ErrValidationNotBlank
	}

	var count int64
	if res := r.db.Model(&model.UpdateStateTransition{}).
		Joins("JOIN update_state_definitions AS from_state ON from_state.id = update_state_transitions.from_state_id").
		Joins("JOIN update_state_definitions AS to_state ON to_state.id = update_state_transitions.to_state_id").
		Where("from_state.name = ? AND to_state.name = ?", fromStateName, toStateName).
		Count(&count); res.Error != nil {
		return false, service_error.NewServiceDatabaseError(res.Error)
	}

	return count > 0, nil
}

func (r *UpdateStateTransitionDbRepo) Exists(fromStateId string, toStateId string) (bool, error) {
	if fromStateId == "" || toStateId == "" {
		return false, service_error.ErrValidationNotBlank
	}

	var count int64
	if res := r.db.Model(&model.UpdateStateTransition{}).Where("from_state_id = ? AND to_state_id = ?", fromStateId, toStateId).Count(&count); res.Error != nil {
		return false, service_error.NewServiceDatabaseError(res.Error)
	}

	return count > 0, nil
}

func (r *UpdateStateTransitionDbRepo) Create(fromStateId string, toStateId string) (*model.UpdateStateTransition, error) {
	if fromStateId == "" || toStateId == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e := &model.UpdateStateTransition{
		FromStateID: fromStateId,
		ToStateID:   toStateId,
	}

	var res *gorm.DB
	if res = r.db.Create(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrDatabaseRowsExpected
	}

	// Reload with preloaded associations
	return r.Find(e.ID.String())
}

func (r *UpdateStateTransitionDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.UpdateStateTransition{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}
