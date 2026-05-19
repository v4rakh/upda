package service

import (
	"errors"

	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/rs/zerolog/log"
)

type UpdateStateDefinitionService struct {
	repo       repository.UpdateStateDefinitionRepository
	updateRepo repository.UpdateRepository
}

func NewUpdateStateDefinitionService(r repository.UpdateStateDefinitionRepository, updateRepo repository.UpdateRepository) *UpdateStateDefinitionService {
	return &UpdateStateDefinitionService{
		repo:       r,
		updateRepo: updateRepo,
	}
}

func (s *UpdateStateDefinitionService) Get(id string) (*model.UpdateStateDefinition, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.Find(id)
}

func (s *UpdateStateDefinitionService) GetByName(name string) (*model.UpdateStateDefinition, error) {
	if name == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.FindByName(name)
}

func (s *UpdateStateDefinitionService) GetInitial() (*model.UpdateStateDefinition, error) {
	return s.repo.FindInitial()
}

func (s *UpdateStateDefinitionService) GetAll() ([]*model.UpdateStateDefinition, error) {
	return s.repo.FindAll()
}

func (s *UpdateStateDefinitionService) ExistsByName(name string) (bool, error) {
	if name == "" {
		return false, service_error.ErrValidationNotBlank
	}

	return s.repo.ExistsByName(name)
}

func (s *UpdateStateDefinitionService) Create(name string, label string, color string, icon string, description *string, isInitial bool, skipOnNewVersion bool) (*model.UpdateStateDefinition, error) {
	if name == "" || label == "" || color == "" || icon == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	// Check if name already exists
	exists, err := s.repo.ExistsByName(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, service_error.ErrResourceConflict
	}

	// If this is being set as initial, clear other initial flags
	if isInitial {
		if err := s.repo.ClearInitial(); err != nil {
			return nil, err
		}
	}

	// Auto-assign sortOrder based on max existing + 1
	maxSortOrder, err := s.repo.MaxSortOrder()
	if err != nil {
		return nil, err
	}
	sortOrder := maxSortOrder + 1

	e, err := s.repo.Create(name, label, color, icon, description, isInitial, skipOnNewVersion, sortOrder)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Created update state definition '%s'", e.Name)
	return e, nil
}

func (s *UpdateStateDefinitionService) Update(id string, name string, label string, color string, icon string, description *string, isInitial bool, skipOnNewVersion bool, sortOrder int) (*model.UpdateStateDefinition, error) {
	if id == "" || name == "" || label == "" || color == "" || icon == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	// Check if state exists
	existing, err := s.repo.Find(id)
	if err != nil {
		return nil, err
	}

	// Check if name is being changed to one that already exists
	if existing.Name != name {
		exists, err := s.repo.ExistsByName(name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, service_error.ErrResourceConflict
		}
	}

	// If this is being set as initial, clear other initial flags
	if isInitial && !existing.IsInitial {
		if err := s.repo.ClearInitial(); err != nil {
			return nil, err
		}
	}

	// If this was the initial state and is being unset, ensure there's still an initial state
	if existing.IsInitial && !isInitial {
		return nil, service_error.ErrInitialStateRequired
	}

	e, err := s.repo.Update(id, name, label, color, icon, description, isInitial, skipOnNewVersion, sortOrder)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Updated update state definition '%s'", e.Name)
	return e, nil
}

func (s *UpdateStateDefinitionService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	// Check if state exists
	existing, err := s.repo.Find(id)
	if err != nil {
		return err
	}

	// Prevent deletion of initial state
	if existing.IsInitial {
		return service_error.ErrInitialStateRequired
	}

	// Check if any updates are using this state
	updates, err := s.updateRepo.Paginate(1, 1, "id", "asc", "", "", existing.Name)
	if err != nil && !errors.Is(err, service_error.ErrResourceNotFound) {
		return err
	}
	if len(updates) > 0 {
		return service_error.ErrStateInUse
	}

	if _, err := s.repo.Delete(id); err != nil {
		return err
	}

	log.Info().Msgf("Deleted update state definition '%s'", existing.Name)
	return nil
}

func (s *UpdateStateDefinitionService) Reorder(items []dto.UpdateStateReorderItem) error {
	if len(items) == 0 {
		return nil
	}

	if err := s.repo.Reorder(items); err != nil {
		return err
	}

	log.Info().Msgf("Reordered %d update state definitions", len(items))
	return nil
}
