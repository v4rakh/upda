package server

import (
	"git.myservermanager.com/varakh/upda/api"
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

	var eventName api.EventName

	if old.State == new.State {
		eventName = api.EventNameUpdateUpdated
	} else {
		switch old.State {
		case api.UpdateStatePending.Value():
			eventName = api.EventNameUpdateUpdatedPending
			break
		case api.UpdateStateApproved.Value():
			eventName = api.EventNameUpdateUpdatedApproved
			break
		case api.UpdateStateIgnored.Value():
			eventName = api.EventNameUpdateUpdatedIgnored
			break
		}
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
	})

	return nil
}

func (s *eventService) createWebhookCreated(e *Webhook) *Event {
	if e == nil {
		return nil
	}

	s.createWithWarnOnly(api.EventNameWebhookCreated, &api.EventPayloadWebhookCreatedDto{
		ID:         e.ID,
		Label:      e.Label,
		Type:       e.Type,
		IgnoreHost: e.IgnoreHost,
	})

	return nil
}

func (s *eventService) createWebhookUpdated(old *Webhook, new *Webhook) *Event {
	if old == nil || new == nil {
		return nil
	}

	var eventName api.EventName

	if old.Label == new.Label {
		eventName = api.EventNameWebhookUpdatedIgnoreHost
	} else {
		eventName = api.EventNameWebhookUpdatedLabel
	}

	s.createWithWarnOnly(eventName, &api.EventPayloadWebhookUpdatedDto{
		ID:              new.ID,
		LabelPrior:      old.Label,
		Label:           new.Label,
		IgnoreHostPrior: old.IgnoreHost,
		IgnoreHost:      new.IgnoreHost,
		Type:            new.Type,
	})

	return nil
}

func (s *eventService) createWebhookDeleted(e *Webhook) *Event {
	if e == nil {
		return nil
	}

	s.createWithWarnOnly(api.EventNameWebhookDeleted, &api.EventPayloadWebhookDeletedDto{
		Label:      e.Label,
		Type:       e.Type,
		IgnoreHost: e.IgnoreHost,
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
