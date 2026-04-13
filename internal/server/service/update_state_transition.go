package service

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/rs/zerolog/log"
)

type UpdateStateTransitionService struct {
	repo         repository.UpdateStateTransitionRepository
	stateDefRepo repository.UpdateStateDefinitionRepository
}

func NewUpdateStateTransitionService(r repository.UpdateStateTransitionRepository, stateDefRepo repository.UpdateStateDefinitionRepository) *UpdateStateTransitionService {
	return &UpdateStateTransitionService{
		repo:         r,
		stateDefRepo: stateDefRepo,
	}
}

func (s *UpdateStateTransitionService) Get(id string) (*model.UpdateStateTransition, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.Find(id)
}

func (s *UpdateStateTransitionService) GetAll() ([]*model.UpdateStateTransition, error) {
	return s.repo.FindAll()
}

func (s *UpdateStateTransitionService) GetByFromStateId(fromStateId string) ([]*model.UpdateStateTransition, error) {
	if fromStateId == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.FindByFromStateId(fromStateId)
}

func (s *UpdateStateTransitionService) GetByFromStateName(fromStateName string) ([]*model.UpdateStateTransition, error) {
	if fromStateName == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.FindByFromStateName(fromStateName)
}

// IsTransitionAllowed checks if a transition from one state to another is allowed.
// Returns true if:
// - The transition explicitly exists in the transitions table
// - There are no transitions defined from the source state (allow-all fallback)
func (s *UpdateStateTransitionService) IsTransitionAllowed(fromStateName string, toStateName string) (bool, error) {
	if fromStateName == "" || toStateName == "" {
		return false, service_error.ErrValidationNotBlank
	}

	// Same state transition is always allowed
	if fromStateName == toStateName {
		return true, nil
	}

	// Check if target state exists
	_, err := s.stateDefRepo.FindByName(toStateName)
	if err != nil {
		return false, err
	}

	// Check if transition exists
	allowed, err := s.repo.IsTransitionAllowed(fromStateName, toStateName)
	if err != nil {
		return false, err
	}
	if allowed {
		return true, nil
	}

	// Check if there are any transitions from the source state
	// If not, allow all transitions (fallback behavior)
	transitions, err := s.repo.FindByFromStateName(fromStateName)
	if err != nil {
		return false, err
	}
	if len(transitions) == 0 {
		return true, nil
	}

	return false, nil
}

func (s *UpdateStateTransitionService) Create(fromStateId string, toStateId string) (*model.UpdateStateTransition, error) {
	if fromStateId == "" || toStateId == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	// Verify from state exists
	_, err := s.stateDefRepo.Find(fromStateId)
	if err != nil {
		return nil, err
	}

	// Verify to state exists
	_, err = s.stateDefRepo.Find(toStateId)
	if err != nil {
		return nil, err
	}

	// Cannot create self-transition
	if fromStateId == toStateId {
		return nil, service_error.ErrValidationNotBlank
	}

	// Check if transition already exists
	exists, err := s.repo.Exists(fromStateId, toStateId)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, service_error.ErrResourceConflict
	}

	e, err := s.repo.Create(fromStateId, toStateId)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Created update state transition from '%s' to '%s'", e.FromState.Name, e.ToState.Name)
	return e, nil
}

func (s *UpdateStateTransitionService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	existing, err := s.repo.Find(id)
	if err != nil {
		return err
	}

	if _, err := s.repo.Delete(id); err != nil {
		return err
	}

	log.Info().Msgf("Deleted update state transition from '%s' to '%s'", existing.FromState.Name, existing.ToState.Name)
	return nil
}
