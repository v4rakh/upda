//nolint:dupl
package service

import (
	"errors"

	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/rs/zerolog/log"
)

type SecretService struct {
	repo repository.SecretRepository
}

func NewSecretService(r repository.SecretRepository) *SecretService {
	return &SecretService{
		repo: r,
	}
}

func (s *SecretService) Get(id string) (*model.Secret, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.FindById(id)
}

func (s *SecretService) GetValueByKey(key string) (string, error) {
	if key == "" {
		return "", service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Secret

	if e, err = s.repo.FindByKey(key); err != nil {
		return "", err
	}

	return e.Value, nil
}

func (s *SecretService) GetAll() ([]*model.Secret, error) {
	return s.repo.FindAll()
}

func (s *SecretService) Insert(key string, value string) (*model.Secret, error) {
	if key == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Secret
	var err error

	e, err = s.repo.FindByKey(key)

	if err != nil && !errors.Is(err, service_error.ErrResourceNotFound) {
		return nil, err
	} else if err != nil && errors.Is(err, service_error.ErrResourceNotFound) {
		if e, err = s.repo.Create(key, value); err != nil {
			return nil, err
		}
		log.Info().Msgf("Created secret '%s'", e.Key)
	} else {
		return nil, service_error.ErrResourceConflict
	}

	return e, err
}

func (s *SecretService) UpdateValue(id string, value string) (*model.Secret, error) {
	if id == "" || value == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Secret
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.Update(id, value); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified secret '%v'", id)
	return e, nil
}

func (s *SecretService) Delete(id string) error {
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

	log.Info().Msgf("Deleted secret '%v'", id)
	return nil
}
