package service

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"go.uber.org/zap"
)

type FilterPresetService struct {
	repo repository.FilterPresetRepository
}

func NewFilterPresetService(r repository.FilterPresetRepository) *FilterPresetService {
	return &FilterPresetService{
		repo: r,
	}
}

func (s *FilterPresetService) GetByType(t api.FilterPresetType) ([]*model.FilterPreset, error) {
	if t == "" {
		return nil, service_error.ErrValidationNotBlank
	}
	return s.repo.FindByType(t.Value())
}

func (s *FilterPresetService) Create(t api.FilterPresetType, label string, parameters string, color *string) (*model.FilterPreset, error) {
	if t == "" || label == "" || parameters == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.FilterPreset
	var err error

	if e, err = s.repo.Create(t.Value(), label, parameters, color); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Created filter preset '%s' ('%s')", e.Type, e.Label)

	return e, err
}

func (s *FilterPresetService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	var err error
	if _, err = s.repo.FindById(id); err != nil {
		return err
	}

	if _, err = s.repo.Delete(id); err != nil {
		return err
	}

	zap.L().Sugar().Infof("Deleted filter preset '%v'", id)
	return nil
}
