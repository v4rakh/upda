package service

import (
	"errors"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/json"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"go.uber.org/zap"
	"time"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(r repository.EventRepository) *EventService {
	return &EventService{
		repo: r,
	}
}

func (s *EventService) CreateUpdateCreated(e *model.Update) *model.Event {
	if e == nil {
		return nil
	}

	s.CreateWithWarnOnly(constant.EventNameUpdateCreated, &api.EventPayloadUpdateCreatedDto{
		ID:          e.ID,
		Application: e.Application,
		Provider:    e.Provider,
		Host:        e.Host,
		Version:     e.Version,
		State:       e.State,
	})

	return nil
}

func (s *EventService) CreateUpdateUpdated(old *model.Update, new *model.Update) *model.Event {
	if old == nil || new == nil {
		return nil
	}

	eventName := constant.EventNameUpdateUpdated

	if old.State != new.State {
		eventName = constant.EventNameUpdateUpdatedState
	}

	if old.Version != new.Version {
		eventName = constant.EventNameUpdateUpdatedVersion
	}

	s.CreateWithWarnOnly(eventName, &api.EventPayloadUpdateUpdatedDto{
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

func (s *EventService) CreateUpdateDeleted(e *model.Update) *model.Event {
	if e == nil {
		return nil
	}

	s.CreateWithWarnOnly(constant.EventNameUpdateDeleted, &api.EventPayloadUpdateDeletedDto{
		Application: e.Application,
		Provider:    e.Provider,
		Host:        e.Host,
		Version:     e.Version,
		State:       e.State,
	})

	return nil
}

func (s *EventService) CreateWithWarnOnly(name constant.EventName, payload interface{}) *model.Event {
	var e *model.Event
	var err error

	if e, err = s.Create(name, payload); err != nil {
		zap.L().Sugar().Warnf("Could not create event '%s': %v", name, err)
		return nil
	}

	return e
}

func (s *EventService) Create(name constant.EventName, payload interface{}) (*model.Event, error) {
	if name == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Event
	var err error

	if e, err = s.repo.Create(name.String(), constant.EventStateCreated.String(), payload); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Created event '%v'", e.Name)

	return e, err
}

func (s *EventService) Get(id string) (*model.Event, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.Find(id)
}

func (s *EventService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	if _, err := s.Get(id); err != nil {
		return err
	}

	if _, err := s.repo.Delete(id); err != nil {
		return err
	}

	zap.L().Sugar().Infof("Deleted event '%v'", id)
	return nil
}

func (s *EventService) CleanStale(time time.Time, state ...constant.EventState) (int64, error) {
	if len(state) == 0 {
		return 0, service_error.ErrValidationNotEmpty
	}

	return s.repo.DeleteByUpdatedAtBeforeAndStates(time, constant.FromVariadicToStr(state...)...)
}

func (s *EventService) Window(size int, skip int, orderBy string, order string) ([]*model.Event, error) {
	return s.repo.Window(size, skip, orderBy, order)
}

func (s *EventService) WindowHasNext(size int, skip int, orderBy string, order string) (bool, error) {
	return s.repo.WindowHasNext(size, skip, orderBy, order)
}

func (s *EventService) Count(state ...constant.EventState) (int64, error) {
	return s.repo.Count(constant.FromVariadicToStr(state...)...)
}

func (s *EventService) GetByState(limit int, state ...constant.EventState) ([]*model.Event, error) {
	if len(state) == 0 {
		return nil, service_error.ErrValidationNotEmpty
	}
	if limit <= 0 {
		return nil, service_error.ErrValidationLimitGreaterZero
	}

	return s.repo.FindAllByState(limit, constant.FromVariadicToStr(state...)...)
}

func (s *EventService) UpdateState(id string, state constant.EventState) (*model.Event, error) {
	if id == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Event
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateState(id, state.String()); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified event '%v'", id)
	return e, nil
}

func (s *EventService) ExtractPayloadInfo(event *model.Event) (*dto.EventPayloadInformationDto, error) {
	if event == nil {
		return nil, service_error.ErrValidationNotEmpty
	}

	var err error
	var bytes []byte

	if bytes, err = event.Payload.MarshalJSON(); err != nil {
		return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, err)
	}

	switch event.Name {
	case constant.EventNameUpdateCreated.String():
		var p api.EventPayloadUpdateCreatedDto
		if p, err = json.UnmarshalGenericJSON[api.EventPayloadUpdateCreatedDto](bytes); err != nil {
			return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, err)
		}
		return &dto.EventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case constant.EventNameUpdateDeleted.String():
		var p api.EventPayloadUpdateDeletedDto
		if p, err = json.UnmarshalGenericJSON[api.EventPayloadUpdateDeletedDto](bytes); err != nil {
			return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, err)
		}
		return &dto.EventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case constant.EventNameUpdateUpdatedState.String():
		var p api.EventPayloadUpdateUpdatedDto
		if p, err = json.UnmarshalGenericJSON[api.EventPayloadUpdateUpdatedDto](bytes); err != nil {
			return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, err)
		}
		return &dto.EventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case constant.EventNameUpdateUpdatedVersion.String():
		var p api.EventPayloadUpdateUpdatedDto
		if p, err = json.UnmarshalGenericJSON[api.EventPayloadUpdateUpdatedDto](bytes); err != nil {
			return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, err)
		}
		return &dto.EventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	case constant.EventNameUpdateUpdated.String():
		var p api.EventPayloadUpdateUpdatedDto
		if p, err = json.UnmarshalGenericJSON[api.EventPayloadUpdateUpdatedDto](bytes); err != nil {
			return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, err)
		}
		return &dto.EventPayloadInformationDto{Host: p.Host, Application: p.Application, Provider: p.Provider, Version: p.Version, State: p.State}, nil
	}

	return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, errors.New("no matching event found"))
}
