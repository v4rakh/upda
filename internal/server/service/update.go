package service

import (
	"errors"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/rs/zerolog/log"
)

type UpdateService struct {
	repo                   repository.UpdateRepository
	eventService           *EventService
	stateDefService        *UpdateStateDefinitionService
	stateTransitionService *UpdateStateTransitionService
}

func NewUpdateService(r repository.UpdateRepository, e *EventService, stateDefService *UpdateStateDefinitionService, stateTransitionService *UpdateStateTransitionService) *UpdateService {
	return &UpdateService{
		repo:                   r,
		eventService:           e,
		stateDefService:        stateDefService,
		stateTransitionService: stateTransitionService,
	}
}

func (s *UpdateService) Get(id string) (*model.Update, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.Find(id)
}

func (s *UpdateService) GetAll() ([]*model.Update, error) {
	return s.repo.FindAll()
}

func (s *UpdateService) Upsert(application string, provider string, host string, version string, metadata interface{}) (*model.Update, error) {
	if application == "" || provider == "" || host == "" || version == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Update
	var err error

	e, err = s.repo.FindBy(application, provider, host)

	if err != nil && !errors.Is(err, service_error.ErrResourceNotFound) {
		return nil, err
	} else if err != nil && errors.Is(err, service_error.ErrResourceNotFound) {
		// Get dynamic initial state - required
		initialState, err := s.stateDefService.GetInitial()
		if err != nil {
			return nil, service_error.ErrInitialStateRequired
		}
		if e, err = s.repo.Create(application, provider, host, version, initialState.Name, metadata); err != nil {
			return nil, err
		}
		s.eventService.CreateUpdateCreated(e)
		log.Info().Msgf("Created update '%v'", e)
	} else {
		old := e

		// Check if current state has SkipOnNewVersion flag
		currentStateDef, stateErr := s.stateDefService.GetByName(e.State)
		if stateErr != nil {
			log.Warn().Msgf("Could not find state definition for '%s', proceeding with update", e.State)
		} else if currentStateDef.SkipOnNewVersion {
			log.Info().Msgf("Skipping update '%v' in state '%s' (skipOnNewVersion=true)", e.ID, e.State)
			return nil, nil
		}

		if e, err = s.repo.Update(e.ID.String(), version, metadata); err != nil {
			return nil, err
		}

		s.eventService.CreateUpdateUpdated(old, e)
		log.Info().Msgf("Updated update '%v'", e)

		// If current state is NOT the initial state, reset to initial state
		initialState, initErr := s.stateDefService.GetInitial()
		if initErr != nil {
			log.Warn().Msg("Could not find initial state for reset")
		} else if e.State != initialState.Name {
			log.Info().Msgf("Setting update '%v' state from '%s' to initial state '%s'", e.ID, e.State, initialState.Name)
			if e, err = s.repo.UpdateState(e.ID.String(), initialState.Name); err != nil {
				return nil, err
			}
		}
	}

	return e, err
}

func (s *UpdateService) UpdateState(id string, state string) (*model.Update, error) {
	if id == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Update
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	// Validate target state exists
	_, err = s.stateDefService.GetByName(state)
	if err != nil {
		if errors.Is(err, service_error.ErrResourceNotFound) {
			return nil, service_error.ErrStateNotFound
		}
		return nil, err
	}

	// Validate transition is allowed
	allowed, err := s.stateTransitionService.IsTransitionAllowed(e.State, state)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, service_error.ErrStateTransitionNotAllowed
	}

	oldUpdate := e
	if e, err = s.repo.UpdateState(id, state); err != nil {
		return nil, err
	}

	s.eventService.CreateUpdateUpdated(oldUpdate, e)

	log.Info().Msgf("Modified update '%v'", id)
	return e, nil
}

func (s *UpdateService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	var e *model.Update
	var err error
	if e, err = s.Get(id); err != nil {
		return err
	}

	if _, err = s.repo.Delete(id); err != nil {
		return err
	}

	s.eventService.CreateUpdateDeleted(e)

	log.Info().Msgf("Deleted update '%v'", id)
	return nil
}

func (s *UpdateService) Paginate(page int, pageSize int, orderBy string, order string, searchTerm string, searchIn string, state ...string) ([]*model.Update, error) {
	return s.repo.Paginate(page, pageSize, orderBy, order, searchTerm, searchIn, state...)
}

func (s *UpdateService) Count(searchTerm string, searchIn string, state ...string) (int64, error) {
	return s.repo.Count(searchTerm, searchIn, state...)
}
