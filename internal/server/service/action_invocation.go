package service

import (
	"errors"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/json"
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	strutil "git.myservermanager.com/varakh/upda/internal/str"
	"github.com/containrrr/shoutrrr"
	"go.uber.org/zap"
	"strings"
	"time"
)

type ActionInvocationService struct {
	repo            repository.ActionInvocationRepository
	actionService   *ActionService
	eventService    *EventService
	secretService   *SecretService
	constantService *ConstantService
}

func NewActionInvocationService(r repository.ActionInvocationRepository, a *ActionService, e *EventService, s *SecretService, c *ConstantService) *ActionInvocationService {
	return &ActionInvocationService{
		repo:            r,
		actionService:   a,
		eventService:    e,
		secretService:   s,
		constantService: c,
	}
}

func (s *ActionInvocationService) Enqueue(batchSize int) error {
	if batchSize <= 0 {
		return service_error.NewServiceError(service_error.ErrCodeGeneral, errors.New("cannot enqueue actions from events with invalid configured batch size"))
	}

	var events []*model.Event
	var err error

	if events, err = s.eventService.GetByState(batchSize, api.EventStateCreated); err != nil {
		return err
	}

	var actions []*model.Action
	if actions, err = s.actionService.GetByEnabled(true); err != nil {
		return err
	}

	for _, event := range events {
		if err = s.EnqueueFromEvent(event, actions); err != nil {
			zap.L().Sugar().Errorf("Could not enqueue action for event '%s' (%s). Reason: %s", event.Name, event.ID, err.Error())
		}
	}

	return nil
}

func (s *ActionInvocationService) EnqueueFromEvent(event *model.Event, actions []*model.Action) error {
	if event == nil || actions == nil {
		return service_error.NewServiceError(service_error.ErrCodeIllegalArgument, service_error.ErrValidationNotEmpty)
	}

	var err error

	// match requires event payload
	var eventPayload *dto.EventPayloadInformationDto
	if eventPayload, err = s.eventService.ExtractPayloadInfo(event); err != nil {
		return err
	}

	filteredActions := make([]*model.Action, 0)

	for _, action := range actions {
		matchesEvent := action.MatchEvent == nil || *action.MatchEvent == event.Name
		matchesHost := action.MatchHost == nil || *action.MatchHost == eventPayload.Host
		matchesApplication := action.MatchApplication == nil || *action.MatchApplication == eventPayload.Application
		matchesProvider := action.MatchProvider == nil || *action.MatchProvider == eventPayload.Provider

		if matchesEvent && matchesHost && matchesApplication && matchesProvider {
			filteredActions = append(filteredActions, action)
		}
	}

	if len(filteredActions) == 0 {
		zap.L().Sugar().Debugf("No actions found which match event '%s', nothing to enqueue", event.Name)
	}

	for _, action := range filteredActions {
		if _, err = s.Create(event, action, api.ActionInvocationStateCreated); err != nil {
			zap.L().Sugar().Errorf("Could not enqueue action '%s' (%v). Reason: %s", action.Label, action.ID, err.Error())
			continue
		}
	}

	// mark event as enqueued
	if _, err = s.eventService.UpdateState(event.ID.String(), api.EventStateEnqueued); err != nil {
		return err
	}

	return nil
}

func (s *ActionInvocationService) Invoke(batchSize int, maxRetries int) error {
	if batchSize <= 0 {
		return service_error.NewServiceError(service_error.ErrCodeGeneral, errors.New("cannot invoke actions with invalid configured batch size"))
	}
	if maxRetries <= 0 {
		return service_error.NewServiceError(service_error.ErrCodeGeneral, errors.New("cannot invoke actions with invalid configured max retries"))
	}

	var err error
	var actionInvocations []*model.ActionInvocation

	if actionInvocations, err = s.GetByState(batchSize, maxRetries, api.ActionInvocationStateCreated, api.ActionInvocationStateRetrying); err != nil {
		return err
	}

	if len(actionInvocations) == 0 {
		zap.L().Sugar().Debugf("No action invocations found to process")
		return nil
	}

	for _, actionInvocation := range actionInvocations {
		if _, err = s.UpdateState(actionInvocation.ID.String(), api.ActionInvocationStateRunning); err != nil {
			zap.L().Sugar().Errorf("Could not mark action invocation '%v' as running. Reason: %s", actionInvocation.ID, err.Error())
			continue
		}

		zap.L().Sugar().Debugf("Invoking action '%v' for event '%v'", actionInvocation.ActionID, actionInvocation.EventID)

		var event *model.Event
		if event, err = s.eventService.Get(actionInvocation.EventID); err != nil {
			zap.L().Sugar().Errorf("Could not find event '%v' for action '%v' and action invocation '%v'. Reason: %s", actionInvocation.EventID, actionInvocation.ActionID, actionInvocation.ID, err.Error())
			// with cascade, cannot happen
			continue
		}

		var eventPayload *dto.EventPayloadInformationDto
		if eventPayload, err = s.eventService.ExtractPayloadInfo(event); err != nil {
			zap.L().Sugar().Errorf("Could not extract event's '%v' information for action '%v' and action invocation '%v'. Reason: %s", actionInvocation.EventID, actionInvocation.ActionID, actionInvocation.ID, err.Error())
			// with layout of attached payload, cannot happen
			continue
		}

		var action *model.Action
		if action, err = s.actionService.Get(actionInvocation.ActionID); err != nil {
			zap.L().Sugar().Errorf("Could not find action '%v' for action invocation '%v'. Reason: %s", actionInvocation.ActionID, actionInvocation.ID, err.Error())
			// with cascade, cannot happen
			continue
		}

		if err = s.Execute(action, eventPayload); err != nil {
			var cause error
			cause = err

			zap.L().Sugar().Errorf("Could not invoke action '%s' (%v) for action invocation '%v'. Reason: %s", action.Label, action.ID, actionInvocation.ID, err.Error())

			var newState api.ActionInvocationState
			newRetryCount := actionInvocation.RetryCount + 1
			newState = api.ActionInvocationStateRetrying

			if newRetryCount >= maxRetries {
				zap.L().Sugar().Infof("Action invocation '%v' exceeded max retry count of '%d'. Not trying again.", actionInvocation.ID, newRetryCount)
				newState = api.ActionInvocationStateError
			}

			if _, err = s.UpdateState(actionInvocation.ID.String(), newState); err != nil {
				zap.L().Sugar().Errorf("Could not mark action invocation '%v' as '%v'. Reason: %s", actionInvocation.ID, newState, err.Error())
			}

			if _, err = s.UpdateRetryCount(actionInvocation.ID.String(), newRetryCount); err != nil {
				zap.L().Sugar().Errorf("Could not update action invocation '%v' retry count to '%d'. Reason: %s", actionInvocation.ID, newRetryCount, err.Error())
			}

			msg := cause.Error()
			if _, err = s.UpdateMessage(actionInvocation.ID.String(), &msg); err != nil {
				zap.L().Sugar().Errorf("Could not update action invocation '%v' message. Reason: %s", actionInvocation.ID, err.Error())
			}

			continue
		}

		zap.L().Sugar().Debugf("Processed action invocation '%v' for event '%s' (%v) and action '%s' (%v)", actionInvocation.ID, event.Name, event.ID, action.Label, action.ID)
		if _, err = s.UpdateState(actionInvocation.ID.String(), api.ActionInvocationStateSuccess); err != nil {
			zap.L().Sugar().Errorf("Could not mark action invocation '%v' as success. Reason: %s", actionInvocation.ID, err.Error())
		}
		if _, err = s.UpdateMessage(actionInvocation.ID.String(), nil); err != nil {
			zap.L().Sugar().Errorf("Could not update action invocation '%v' message. Reason: %s", actionInvocation.ID, err.Error())
		}
	}

	return nil
}

func (s *ActionInvocationService) Execute(action *model.Action, eventPayloadInfo *dto.EventPayloadInformationDto) error {
	if action == nil || eventPayloadInfo == nil {
		return service_error.ErrValidationNotEmpty
	}

	var err error
	var bytes []byte

	if bytes, err = action.Payload.MarshalJSON(); err != nil {
		return service_error.NewServiceError(service_error.ErrCodeGeneral, err)
	}

	switch action.Type {
	case api.ActionTypeShoutrrr.Value():
		var payload dto.ActionPayloadShoutrrrDto
		if payload, err = json.UnmarshalGenericJSON[dto.ActionPayloadShoutrrrDto](bytes); err != nil {
			return service_error.NewServiceError(service_error.ErrCodeGeneral, err)
		}

		body := s.replaceConstants(payload.Body)
		body = s.replaceVars(body, eventPayloadInfo)
		body = s.replaceSecrets(body)

		for _, url := range payload.Urls {
			url = s.replaceConstants(url)
			url = s.replaceVars(url, eventPayloadInfo)
			url = s.replaceSecrets(url)
			if err = shoutrrr.Send(url, body); err != nil {
				return err
			}
		}
		break
	default:
		return service_error.NewServiceError(service_error.ErrCodeGeneral, errors.New("no matching action type found for invocation"))
	}

	return nil
}

func (s *ActionInvocationService) replaceSecrets(str string) string {
	if str == "" {
		return str
	}

	var matches [][]string

	matches = strutil.ExtractBetween(str, "<SECRET>", "</SECRET>")
	var err error

	for _, match := range matches {
		var val string
		if val, err = s.secretService.GetValueByKey(match[1]); err != nil {
			zap.L().Sugar().Warnf("Could not inject secret '%s'. Reason: %s", match[1], err.Error())
			continue
		}
		str = strings.ReplaceAll(str, match[0], val)
	}

	return str
}

func (s *ActionInvocationService) replaceConstants(str string) string {
	if str == "" {
		return str
	}

	var matches [][]string

	matches = strutil.ExtractBetween(str, "<CONST>", "</CONST>")
	var err error

	for _, match := range matches {
		var val string
		if val, err = s.constantService.GetValueByKey(match[1]); err != nil {
			zap.L().Sugar().Warnf("Could not inject constant '%s'. Reason: %s", match[1], err.Error())
			continue
		}
		str = strings.ReplaceAll(str, match[0], val)
	}

	return str
}

func (s *ActionInvocationService) replaceVars(str string, eventPayloadInfo *dto.EventPayloadInformationDto) string {
	if str == "" || eventPayloadInfo == nil {
		return str
	}

	str = strings.ReplaceAll(str, "<VAR>APPLICATION</VAR>", eventPayloadInfo.Application)
	str = strings.ReplaceAll(str, "<VAR>PROVIDER</VAR>", eventPayloadInfo.Provider)
	str = strings.ReplaceAll(str, "<VAR>HOST</VAR>", eventPayloadInfo.Host)
	str = strings.ReplaceAll(str, "<VAR>VERSION</VAR>", eventPayloadInfo.Version)
	str = strings.ReplaceAll(str, "<VAR>STATE</VAR>", eventPayloadInfo.State)

	return str
}

func (s *ActionInvocationService) Paginate(page int, pageSize int, orderBy string, order string) ([]*model.ActionInvocation, error) {
	return s.repo.Paginate(page, pageSize, orderBy, order)
}

func (s *ActionInvocationService) Get(id string) (*model.ActionInvocation, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e, err := s.repo.Find(id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *ActionInvocationService) GetByState(limit int, maxRetries int, state ...api.ActionInvocationState) ([]*model.ActionInvocation, error) {
	if len(state) == 0 {
		return nil, service_error.ErrValidationNotEmpty
	}
	if limit <= 0 {
		return nil, service_error.ErrValidationLimitGreaterZero
	}
	if maxRetries <= 0 {
		return nil, service_error.ErrValidationMaxRetriesGreaterZero
	}

	return s.repo.FindAllByState(limit, maxRetries, api.FromVariadicToStr(state...)...)
}

func (s *ActionInvocationService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *ActionInvocationService) Delete(id string) error {
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

	zap.L().Sugar().Infof("Deleted action '%v'", id)

	return nil
}

func (s *ActionInvocationService) UpdateState(id string, state api.ActionInvocationState) (*model.ActionInvocation, error) {
	if id == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.ActionInvocation
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateState(id, state.Value()); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action invocation '%v'", id)
	return e, nil
}

func (s *ActionInvocationService) UpdateMessage(id string, message *string) (*model.ActionInvocation, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.ActionInvocation
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateMessage(id, message); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action invocation '%v'", id)
	return e, nil
}

func (s *ActionInvocationService) UpdateRetryCount(id string, retryCount int) (*model.ActionInvocation, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.ActionInvocation
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateRetryCount(id, retryCount); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action invocation '%v'", id)
	return e, nil
}

func (s *ActionInvocationService) Create(event *model.Event, action *model.Action, state api.ActionInvocationState) (*model.ActionInvocation, error) {
	if state == "" {
		return nil, service_error.ErrValidationNotBlank
	}
	if action == nil || event == nil {
		return nil, service_error.ErrValidationNotEmpty
	}

	var err error
	var e *model.ActionInvocation
	if e, err = s.repo.Create(event.ID.String(), action.ID.String(), state.Value()); err != nil {
		return nil, err
	} else {
		zap.L().Sugar().Info("Created action invocation")
		return e, nil
	}
}

func (s *ActionInvocationService) CleanStale(time time.Time, maxRetries int, state ...api.ActionInvocationState) (int64, error) {
	if len(state) == 0 {
		return 0, service_error.ErrValidationNotEmpty
	}

	return s.repo.DeleteByUpdatedAtBeforeAndStates(time, maxRetries, api.FromVariadicToStr(state...)...)
}
