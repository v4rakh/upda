package server

import (
	"errors"
	"git.myservermanager.com/varakh/upda/api"
	"go.uber.org/zap"
	"time"
)

type updateService struct {
	repo              updateRepository
	eventService      *eventService
	prometheusService *prometheusService
}

func newUpdateService(r updateRepository, e *eventService, p *prometheusService) *updateService {
	return &updateService{
		repo:              r,
		eventService:      e,
		prometheusService: p,
	}
}

func (s *updateService) get(id string) (*Update, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	return s.repo.find(id)
}

func (s *updateService) getAll() ([]*Update, error) {
	return s.repo.findAll()
}

func (s *updateService) upsert(application string, provider string, host string, version string, metadata interface{}) (*Update, error) {
	if application == "" || provider == "" || host == "" || version == "" {
		return nil, errorValidationNotBlank
	}

	var e *Update
	var err error

	e, err = s.repo.findBy(application, provider, host)

	if err != nil && !errors.Is(err, errorResourceNotFound) {
		return nil, err
	} else if err != nil && errors.Is(err, errorResourceNotFound) {
		if e, err = s.repo.create(application, provider, host, version, metadata); err != nil {
			return nil, err
		}
		s.eventService.createUpdateCreated(e)
		zap.L().Sugar().Infof("Created update '%v'", e)
	} else {
		old := e
		skip := e.State == api.UpdateStateIgnored.Value()

		if skip {
			zap.L().Sugar().Infof("Skipping ignored update '%v'", e.ID)
			return nil, nil
		}

		if e, err = s.repo.update(e.ID.String(), version, metadata); err != nil {
			return nil, err
		}

		s.eventService.createUpdateUpdated(old, e)
		zap.L().Sugar().Infof("Updated update '%v'", e)

		if api.UpdateStateApproved.Value() == e.State {
			zap.L().Sugar().Infof("Setting update '%v' state to '%v'", e.ID, api.UpdateStatePending)
			if e, err = s.repo.updateState(e.ID.String(), api.UpdateStatePending); err != nil {
				return nil, err
			}
		}
	}

	return e, err
}

func (s *updateService) updateState(id string, state api.UpdateState) (*Update, error) {
	if id == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	var e *Update
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	oldUpdate := e
	if e, err = s.repo.updateState(id, state); err != nil {
		return nil, err
	}

	s.eventService.createUpdateUpdated(oldUpdate, e)

	zap.L().Sugar().Infof("Modified update '%v'", id)
	return e, nil
}

func (s *updateService) delete(id string) error {
	if id == "" {
		return errorValidationNotBlank
	}

	var e *Update
	var err error
	if e, err = s.get(id); err != nil {
		return err
	}

	if _, err = s.repo.delete(id); err != nil {
		return err
	}

	s.eventService.createUpdateDeleted(e)

	if err = s.prometheusService.setGauge(metricUpdates, []string{e.Application, e.Provider, e.Host}, -1); err != nil {
		zap.L().Sugar().Errorf("Could not refresh updates prometheus metric for deleted update '%v'. Reason: %v", e.ID, err)
	}

	zap.L().Sugar().Infof("Deleted update '%v'", id)
	return nil
}

func (s *updateService) cleanStale(time time.Time, state ...api.UpdateState) (int64, error) {
	if len(state) == 0 {
		return 0, errorValidationNotEmpty
	}

	return s.repo.deleteByUpdatedAtBeforeAndStates(time, state...)
}

func (s *updateService) paginate(page int, pageSize int, orderBy string, order string, searchTerm string, searchIn string, state ...api.UpdateState) ([]*Update, error) {
	return s.repo.paginate(page, pageSize, orderBy, order, searchTerm, searchIn, state...)
}

func (s *updateService) count(searchTerm string, searchIn string, state ...api.UpdateState) (int64, error) {
	return s.repo.count(searchTerm, searchIn, state...)
}
