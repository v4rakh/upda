package server

import (
	"errors"
	"go.uber.org/zap"
)

type secretService struct {
	repo secretRepository
}

func newSecretService(r secretRepository) *secretService {
	return &secretService{
		repo: r,
	}
}

func (s *secretService) get(id string) (*Secret, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	return s.repo.findById(id)
}

func (s *secretService) getValueByKey(key string) (string, error) {
	if key == "" {
		return "", errorValidationNotBlank
	}

	var err error
	var e *Secret

	if e, err = s.repo.findByKey(key); err != nil {
		return "", err
	}

	return e.Value, nil
}

func (s *secretService) getAll() ([]*Secret, error) {
	return s.repo.findAll()
}

func (s *secretService) insert(key string, value string) (*Secret, error) {
	if key == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var e *Secret
	var err error

	e, err = s.repo.findByKey(key)

	if err != nil && !errors.Is(err, errorResourceNotFound) {
		return nil, err
	} else if err != nil && errors.Is(err, errorResourceNotFound) {
		if e, err = s.repo.create(key, value); err != nil {
			return nil, err
		}
		zap.L().Sugar().Infof("Created secret '%s'", e.Key)
	} else {
		return nil, errorResourceConflict
	}

	return e, err
}

func (s *secretService) updateValue(id string, value string) (*Secret, error) {
	if id == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var e *Secret
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.update(id, value); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified secret '%v'", id)
	return e, nil
}

func (s *secretService) delete(id string) error {
	if id == "" {
		return errorValidationNotBlank
	}

	var err error
	if _, err = s.get(id); err != nil {
		return err
	}

	if _, err = s.repo.delete(id); err != nil {
		return err
	}

	zap.L().Sugar().Infof("Deleted secret '%v'", id)
	return nil
}
