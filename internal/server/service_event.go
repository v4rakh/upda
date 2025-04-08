package server

import (
	"errors"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/commons"
	"go.uber.org/zap"
	"time"
)

type eventService struct {
	repo eventRepository
}

func newEventService(r eventRepository) *eventService {
	return &eventService{
		repo: r,
	}
}

func (s *eventService) createUpdateCreated(e *Update) *Event {
	if e == nil {
		return nil
	}

	s.createWithWarnOnly(api.EventNameUpdateCreated, &api.EventPayloadUpdateCreatedDto{
		ID:          e.ID,
		Application: e.Application,
		Provider:    e.Provider,
		Host:        e.Host,
		Version:     e.Version,
		State:       e.State,
	})

	return nil
}

func (s *eventService) createUpdateUpdated(old *Update, new *Update) *Event {
	if old == nil || new == nil {
		return nil
	}

	eventName := api.EventNameUpdateUpdated

	if old.State != new.State {
		eventName = api.EventNameUpdateUpdatedState
	}

	if old.Version != new.Version {
		eventName = api.EventNameUpdateUpdatedVersion
	}

	s.createWithWarnOnly(eventName, &api.EventPayloadUpdateUpdatedDto{
		ID:           new.ID,
		Application:  new.Application,
		Provider:     new.Provider,
		Host:         new.Host,
		VersionPrior: old.Version,
		Version:      new.Version,
		StatePrior:   old.State,
		State:        new.State,
	})

	return nil
}

func (s *eventService) createUpdateDeleted(e *Update) *Event {
	if e == nil {
		return nil
	}

	s.createWithWarnOnly(api.EventNameUpdateDeleted, &api.EventPayloadUpdateDeletedDto{
		Application: e.Application,
		Provider:    e.Provider,
		Host:        e.Host,
		Version:     e.Version,
		State:       e.State,
	})

	return nil
}

func (s *eventService) createWithWarnOnly(name api.EventName, payload interface{}) *Event {
	var e *Event
	var err error

	if e, err = s.create(name, payload); err != nil {
		zap.L().Sugar().Warnf("Could not create event '%s': %v", name, err)
		return nil
	}

	return e
}

func (s *eventService) create(name api.EventName, payload interface{}) (*Event, error) {
	if name == "" {
		return nil, errorValidationNotBlank
	}

	var e *Event
	var err error

	if e, err = s.repo.create(name, api.EventStateCreated, payload); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Created event '%v'", e.Name)

	return e, err
}

func (s *eventService) get(id string) (*Event, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	return s.repo.find(id)
}

func (s *eventService) delete(id string) error {
	if id == "" {
		return errorValidationNotBlank
	}

	if _, err := s.get(id); err != nil {
		return err
	}

	if _, err := s.repo.delete(id); err != nil {
		return err
	}

	zap.L().Sugar().Infof("Deleted event '%v'", id)
	return nil
}

func (s *eventService) cleanStale(time time.Time, state ...api.EventState) (int64, error) {
	if len(state) == 0 {
		return 0, errorValidationNotEmpty
	}

	return s.repo.deleteByUpdatedAtBeforeAndStates(time, state...)
}

func (s *eventService) window(size int, skip int, orderBy string, order string) ([]*Event, error) {
	return s.repo.window(size, skip, orderBy, order)
}

func (s *eventService) windowHasNext(size int, skip int, orderBy string, order string) (bool, error) {
	return s.repo.windowHasNext(size, skip, orderBy, order)
}

func (s *eventService) count(state ...api.EventState) (int64, error) {
	return s.repo.count(state...)
}

func (s *eventService) getByState(limit int, state ...api.EventState) ([]*Event, error) {
	if len(state) == 0 {
		return nil, errorValidationNotEmpty
	}
	if limit <= 0 {
		return nil, errorValidationLimitGreaterZero
	}

	return s.repo.findAllByState(limit, state...)
}

func (s *eventService) updateState(id string, state api.EventState) (*Event, error) {
	if id == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	var e *Event
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateState(id, state); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified event '%v'", id)
	return e, nil
}

func (s *eventService) extractPayloadInfo(event *Event) (*eventPayloadInformationDto, error) {
	if event == nil {
		return nil, errorValidationNotEmpty
	}

	var err error
	var bytes []byte

	if bytes, err = event.Payload.MarshalJSON(); err != nil {
		return nil, newServiceError(general, err)
	}

	switch event.Name {
	case api.EventNameUpdateCreated.Value():
		var p api.EventPayloadUpdateCreatedDto
		if p, err = commons.UnmarshalGenericJSON[api.EventPayloadUpdateCreatedDto](bytes); err != nil {
			return nil, newServiceError(general, err)
		}
		return &eventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case api.EventNameUpdateDeleted.Value():
		var p api.EventPayloadUpdateDeletedDto
		if p, err = commons.UnmarshalGenericJSON[api.EventPayloadUpdateDeletedDto](bytes); err != nil {
			return nil, newServiceError(general, err)
		}
		return &eventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case api.EventNameUpdateUpdatedState.Value():
		var p api.EventPayloadUpdateUpdatedDto
		if p, err = commons.UnmarshalGenericJSON[api.EventPayloadUpdateUpdatedDto](bytes); err != nil {
			return nil, newServiceError(general, err)
		}
		return &eventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case api.EventNameUpdateUpdatedVersion.Value():
		var p api.EventPayloadUpdateUpdatedDto
		if p, err = commons.UnmarshalGenericJSON[api.EventPayloadUpdateUpdatedDto](bytes); err != nil {
			return nil, newServiceError(general, err)
		}
		return &eventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case api.EventNameUpdateUpdated.Value():
		var p api.EventPayloadUpdateUpdatedDto
		if p, err = commons.UnmarshalGenericJSON[api.EventPayloadUpdateUpdatedDto](bytes); err != nil {
			return nil, newServiceError(general, err)
		}
		return &eventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	}

	return nil, newServiceError(general, errors.New("no matching event found"))
}
