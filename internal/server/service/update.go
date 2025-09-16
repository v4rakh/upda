package service

import (
	"errors"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/rs/zerolog/log"
	"time"
)

type UpdateService struct {
	repo         repository.UpdateRepository
	eventService *EventService
}

func NewUpdateService(r repository.UpdateRepository, e *EventService) *UpdateService {
	return &UpdateService{
		repo:         r,
		eventService: e,
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
		if e, err = s.repo.Create(application, provider, host, version, constant.UpdateStatePending.String(), metadata); err != nil {
			return nil, err
		}
		s.eventService.CreateUpdateCreated(e)
		log.Info().Msgf("Created update '%v'", e)
	} else {
		old := e
		skip := e.State == constant.UpdateStateIgnored.String()

		if skip {
			log.Info().Msgf("Skipping ignored update '%v'", e.ID)
			return nil, nil
		}

		if e, err = s.repo.Update(e.ID.String(), version, metadata); err != nil {
			return nil, err
		}

		s.eventService.CreateUpdateUpdated(old, e)
		log.Info().Msgf("Updated update '%v'", e)

		if constant.UpdateStateApproved.String() == e.State {
			log.Info().Msgf("Setting update '%v' state to '%v'", e.ID, constant.UpdateStatePending)
			if e, err = s.repo.UpdateState(e.ID.String(), constant.UpdateStatePending.String()); err != nil {
				return nil, err
			}
		}
	}

	return e, err
}

func (s *UpdateService) UpdateState(id string, state constant.UpdateState) (*model.Update, error) {
	if id == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Update
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	oldUpdate := e
	if e, err = s.repo.UpdateState(id, state.String()); err != nil {
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

func (s *UpdateService) CleanStale(time time.Time, state ...constant.UpdateState) (int64, error) {
	if len(state) == 0 {
		return 0, service_error.ErrValidationNotEmpty
	}

	return s.repo.DeleteByUpdatedAtBeforeAndStates(time, constant.FromVariadicToStr(state...)...)
}

func (s *UpdateService) Paginate(page int, pageSize int, orderBy string, order string, searchTerm string, searchIn string, state ...constant.UpdateState) ([]*model.Update, error) {
	return s.repo.Paginate(page, pageSize, orderBy, order, searchTerm, searchIn, constant.FromVariadicToStr(state...)...)
}

func (s *UpdateService) Count(searchTerm string, searchIn string, state ...constant.UpdateState) (int64, error) {
	return s.repo.Count(searchTerm, searchIn, constant.FromVariadicToStr(state...)...)
}
