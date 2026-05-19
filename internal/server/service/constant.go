//nolint:dupl
package service

import (
	"errors"

	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/rs/zerolog/log"
)

type ConstantService struct {
	repo repository.ConstantRepository
}

func NewConstantService(r repository.ConstantRepository) *ConstantService {
	return &ConstantService{
		repo: r,
	}
}

func (s *ConstantService) Get(id string) (*model.Constant, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.FindById(id)
}

func (s *ConstantService) GetValueByKey(key string) (string, error) {
	if key == "" {
		return "", service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Constant

	if e, err = s.repo.FindByKey(key); err != nil {
		return "", err
	}

	return e.Value, nil
}

func (s *ConstantService) GetAll() ([]*model.Constant, error) {
	return s.repo.FindAll()
}

func (s *ConstantService) Insert(key string, value string) (*model.Constant, error) {
	if key == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Constant
	var err error

	_, err = s.repo.FindByKey(key)

	if err != nil && !errors.Is(err, service_error.ErrResourceNotFound) {
		return nil, err
	} else if err != nil && errors.Is(err, service_error.ErrResourceNotFound) {
		if e, err = s.repo.Create(key, value); err != nil {
			return nil, err
		}
		log.Info().Msgf("Created constant '%s'", e.Key)
	} else {
		return nil, service_error.ErrResourceConflict
	}

	return e, err
}

func (s *ConstantService) UpdateValue(id string, value string) (*model.Constant, error) {
	if id == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Constant
	var err error

	if _, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.Update(id, value); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified constant '%v'", id)
	return e, nil
}

func (s *ConstantService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	var err error
	if _, err = s.Get(id); err != nil {
		return err
	}

	if _, err = s.repo.Delete(id); err != nil {
		return err
	}

	log.Info().Msgf("Deleted constant '%v'", id)
	return nil
}
