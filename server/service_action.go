package server

import (
	"encoding/json"
	"git.myservermanager.com/varakh/upda/api"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type actionService struct {
	repo         ActionRepository
	eventService *eventService
}

func newActionService(r ActionRepository, e *eventService) *actionService {
	return &actionService{
		repo:         r,
		eventService: e,
	}
}

func (s *actionService) get(id string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	e, err := s.repo.find(id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *actionService) create(label string, t api.ActionType, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool) (*Action, error) {
	if label == "" || t == "" {
		return nil, errorValidationNotBlank
	}

	if isValid, validationErr := s.isValidPayload(t, payload); !isValid {
		return nil, newServiceError(IllegalArgument, validationErr)
	}

	var err error
	var e *Action
	if e, err = s.repo.create(label, t, matchEvent, matchHost, matchApplication, matchProvider, payload, enabled); err != nil {
		return nil, err
	} else {
		zap.L().Sugar().Info("Created action")
		return e, nil
	}
}

func (s *actionService) isValidPayload(t api.ActionType, payload interface{}) (bool, error) {
	if t == "" {
		return false, errorValidationNotBlank
	}
	if payload == nil {
		return false, errorValidationNotEmpty
	}

	var err error
	if api.ActionTypeShoutrrr == t {
		var pb []byte
		if pb, err = json.Marshal(payload); err != nil {
			return false, err
		}

		var p actionPayloadShoutrrrDto
		if err = json.Unmarshal(pb, &p); err != nil {
			return false, err
		}

		valid := validator.New()
		if err = valid.Struct(p); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *actionService) updateLabel(id string, label string) (*Action, error) {
	if id == "" || label == "" {
		return nil, errorValidationNotBlank
	}

	var e *Action
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateLabel(id, label); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action '%v'", id)
	return e, nil
}

func (s *actionService) updateMatchEvent(id string, matchEvent *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *Action
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateMatchEvent(id, matchEvent); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action '%v'", id)
	return e, nil
}

func (s *actionService) updateMatchApplication(id string, matchApplication *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *Action
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateMatchApplication(id, matchApplication); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action '%v'", id)
	return e, nil
}

func (s *actionService) updateMatchProvider(id string, matchProvider *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *Action
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateMatchProvider(id, matchProvider); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action '%v'", id)
	return e, nil
}

func (s *actionService) updateMatchHost(id string, matchHost *string) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *Action
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateMatchHost(id, matchHost); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action '%v'", id)
	return e, nil
}

func (s *actionService) updateTypeAndPayload(id string, t api.ActionType, payload interface{}) (*Action, error) {
	if id == "" || t == "" {
		return nil, errorValidationNotBlank
	}
	if payload == nil {
		return nil, errorValidationNotEmpty
	}

	var e *Action
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if isValid, validationErr := s.isValidPayload(t, payload); !isValid {
		return nil, newServiceError(IllegalArgument, validationErr)
	}

	if e, err = s.repo.updateTypeAndPayload(id, t, payload); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action '%v'", id)
	return e, nil
}

func (s *actionService) updateEnabled(id string, enabled bool) (*Action, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *Action
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateEnabled(id, enabled); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action '%v'", id)
	return e, nil
}

func (s *actionService) delete(id string) error {
	if id == "" {
		return errorValidationNotBlank
	}

	e, err := s.get(id)
	if err != nil {
		return err
	}

	if _, err = s.repo.delete(e.ID.String()); err != nil {
		return err
	}

	zap.L().Sugar().Infof("Deleted action '%v'", id)

	return nil
}

func (s *actionService) paginate(page int, pageSize int, orderBy string, order string) ([]*Action, error) {
	return s.repo.paginate(page, pageSize, orderBy, order)
}

func (s *actionService) count() (int64, error) {
	return s.repo.count()
}

func (s *actionService) getAll() ([]*Action, error) {
	return s.repo.findAll()
}

func (s *actionService) getByEnabled(enabled bool) ([]*Action, error) {
	return s.repo.findByEnabled(enabled)
}
