package service

import (
	"encoding/json"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type ActionService struct {
	repo         repository.ActionRepository
	eventService *EventService
}

func NewActionService(r repository.ActionRepository, e *EventService) *ActionService {
	return &ActionService{
		repo:         r,
		eventService: e,
	}
}

func (s *ActionService) Get(id string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e, err := s.repo.Find(id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *ActionService) Create(label string, t constant.ActionType, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool) (*model.Action, error) {
	if label == "" || t == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	if isValid, validationErr := s.IsValidPayload(t, payload); !isValid {
		return nil, service_error.NewServiceError(service_error.ErrCodeIllegalArgument, validationErr)
	}

	var err error
	var e *model.Action
	if e, err = s.repo.Create(label, t.String(), matchEvent, matchHost, matchApplication, matchProvider, payload, enabled); err != nil {
		return nil, err
	} else {
		log.Info().Msg("Created action")
		return e, nil
	}
}

func (s *ActionService) IsValidPayload(t constant.ActionType, payload interface{}) (bool, error) {
	if t == "" {
		return false, service_error.ErrValidationNotBlank
	}
	if payload == nil {
		return false, service_error.ErrValidationNotEmpty
	}

	var err error
	if constant.ActionTypeShoutrrr == t {
		var pb []byte
		if pb, err = json.Marshal(payload); err != nil {
			return false, err
		}

		var p dto.ActionPayloadShoutrrrDto
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

func (s *ActionService) UpdateLabel(id string, label string) (*model.Action, error) {
	if id == "" || label == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Action
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateLabel(id, label); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified action '%v'", id)
	return e, nil
}

func (s *ActionService) UpdateMatchEvent(id string, matchEvent *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Action
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateMatchEvent(id, matchEvent); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified action '%v'", id)
	return e, nil
}

func (s *ActionService) UpdateMatchApplication(id string, matchApplication *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Action
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateMatchApplication(id, matchApplication); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified action '%v'", id)
	return e, nil
}

func (s *ActionService) UpdateMatchProvider(id string, matchProvider *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Action
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateMatchProvider(id, matchProvider); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified action '%v'", id)
	return e, nil
}

func (s *ActionService) UpdateMatchHost(id string, matchHost *string) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Action
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateMatchHost(id, matchHost); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified action '%v'", id)
	return e, nil
}

func (s *ActionService) UpdateTypeAndPayload(id string, t constant.ActionType, payload interface{}) (*model.Action, error) {
	if id == "" || t == "" {
		return nil, service_error.ErrValidationNotBlank
	}
	if payload == nil {
		return nil, service_error.ErrValidationNotEmpty
	}

	var e *model.Action
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if isValid, validationErr := s.IsValidPayload(t, payload); !isValid {
		return nil, service_error.NewServiceError(service_error.ErrCodeIllegalArgument, validationErr)
	}

	if e, err = s.repo.UpdateTypeAndPayload(id, t.String(), payload); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified action '%v'", id)
	return e, nil
}

func (s *ActionService) UpdateEnabled(id string, enabled bool) (*model.Action, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Action
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateEnabled(id, enabled); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified action '%v'", id)
	return e, nil
}

func (s *ActionService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	e, err := s.Get(id)
	if err != nil {
		return err
	}

	if _, err = s.repo.Delete(e.ID.String()); err != nil {
		return err
	}

	log.Info().Msgf("Deleted action '%v'", id)

	return nil
}

func (s *ActionService) Paginate(page int, pageSize int, orderBy string, order string) ([]*model.Action, error) {
	return s.repo.Paginate(page, pageSize, orderBy, order)
}

func (s *ActionService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *ActionService) GetAll() ([]*model.Action, error) {
	return s.repo.FindAll()
}

func (s *ActionService) GetByEnabled(enabled bool) ([]*model.Action, error) {
	return s.repo.FindByEnabled(enabled)
}
