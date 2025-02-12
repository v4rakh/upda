package server

import (
	"errors"
	"go.uber.org/zap"
)

type constantService struct {
	repo constantRepository
}

func newConstantService(r constantRepository) *constantService {
	return &constantService{
		repo: r,
	}
}

func (s *constantService) get(id string) (*Constant, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	return s.repo.findById(id)
}

func (s *constantService) getValueByKey(key string) (string, error) {
	if key == "" {
		return "", errorValidationNotBlank
	}

	var err error
	var e *Constant

	if e, err = s.repo.findByKey(key); err != nil {
		return "", err
	}

	return e.Value, nil
}

func (s *constantService) getAll() ([]*Constant, error) {
	return s.repo.findAll()
}

func (s *constantService) insert(key string, value string) (*Constant, error) {
	if key == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var e *Constant
	var err error

	e, err = s.repo.findByKey(key)

	if err != nil && !errors.Is(err, errorResourceNotFound) {
		return nil, err
	} else if err != nil && errors.Is(err, errorResourceNotFound) {
		if e, err = s.repo.create(key, value); err != nil {
			return nil, err
		}
		zap.L().Sugar().Infof("Created constant '%s'", e.Key)
	} else {
		return nil, errorResourceConflict
	}

	return e, err
}

func (s *constantService) updateValue(id string, value string) (*Constant, error) {
	if id == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var e *Constant
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.update(id, value); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified constant '%v'", id)
	return e, nil
}

func (s *constantService) delete(id string) error {
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

	zap.L().Sugar().Infof("Deleted constant '%v'", id)
	return nil
}
