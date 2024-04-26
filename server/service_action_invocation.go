package server

import (
	"errors"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/util"
	"github.com/containrrr/shoutrrr"
	"go.uber.org/zap"
	"strings"
	"time"
)

type actionInvocationService struct {
	repo          ActionInvocationRepository
	actionService *actionService
	eventService  *eventService
	secretService *secretService
}

func newActionInvocationService(r ActionInvocationRepository, a *actionService, e *eventService, s *secretService) *actionInvocationService {
	return &actionInvocationService{
		repo:          r,
		actionService: a,
		eventService:  e,
		secretService: s,
	}
}

func (s *actionInvocationService) enqueue(batchSize int) error {
	if batchSize <= 0 {
		return newServiceError(General, errors.New("cannot enqueue actions from events with invalid configured batch size"))
	}

	var events []*Event
	var err error

	if events, err = s.eventService.getByState(batchSize, api.EventStateCreated); err != nil {
		return err
	}

	var actions []*Action
	if actions, err = s.actionService.getByEnabled(true); err != nil {
		return err
	}

	for _, event := range events {
		if err = s.enqueueFromEvent(event, actions); err != nil {
			zap.L().Sugar().Errorf("Could not enqueue action for event '%s' (%s). Reason: %s", event.Name, event.ID, err.Error())
		}
	}

	return nil
}

func (s *actionInvocationService) enqueueFromEvent(event *Event, actions []*Action) error {
	if event == nil || actions == nil {
		return newServiceError(IllegalArgument, errorValidationNotEmpty)
	}

	var err error

	// match requires event payload
	var eventPayload *eventPayloadInformationDto
	if eventPayload, err = s.eventService.extractPayloadInfo(event); err != nil {
		return err
	}

	var filteredActions []*Action
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
		if _, err = s.create(event, action, api.ActionInvocationStateCreated); err != nil {
			zap.L().Sugar().Errorf("Could not enqueue action '%s' (%v). Reason: %s", action.Label, action.ID, err.Error())
			continue
		}
	}

	// mark event as enqueued
	if _, err = s.eventService.updateState(event.ID.String(), api.EventStateEnqueued); err != nil {
		return err
	}

	return nil
}

func (s *actionInvocationService) invoke(batchSize int, maxRetries int) error {
	if batchSize <= 0 {
		return newServiceError(General, errors.New("cannot invoke actions with invalid configured batch size"))
	}
	if maxRetries <= 0 {
		return newServiceError(General, errors.New("cannot invoke actions with invalid configured max retries"))
	}

	var err error
	var actionInvocations []*ActionInvocation

	if actionInvocations, err = s.getByState(batchSize, maxRetries, api.ActionInvocationStateCreated, api.ActionInvocationStateRetrying); err != nil {
		return err
	}

	if len(actionInvocations) == 0 {
		zap.L().Sugar().Debugf("No action invocations found to process")
		return nil
	}

	for _, actionInvocation := range actionInvocations {
		if _, err = s.updateState(actionInvocation.ID.String(), api.ActionInvocationStateRunning); err != nil {
			zap.L().Sugar().Errorf("Could not mark action invocation '%v' as running. Reason: %s", actionInvocation.ID, err.Error())
			continue
		}

		zap.L().Sugar().Debugf("Invoking action '%v' for event '%v'", actionInvocation.ActionID, actionInvocation.EventID)

		var event *Event
		if event, err = s.eventService.get(actionInvocation.EventID); err != nil {
			zap.L().Sugar().Errorf("Could not find event '%v' for action '%v' and action invocation '%v'. Reason: %s", actionInvocation.EventID, actionInvocation.ActionID, actionInvocation.ID, err.Error())
			// with cascade, cannot happen
			continue
		}

		var eventPayload *eventPayloadInformationDto
		if eventPayload, err = s.eventService.extractPayloadInfo(event); err != nil {
			zap.L().Sugar().Errorf("Could not extract event's '%v' information for action '%v' and action invocation '%v'. Reason: %s", actionInvocation.EventID, actionInvocation.ActionID, actionInvocation.ID, err.Error())
			// with layout of attached payload, cannot happen
			continue
		}

		var action *Action
		if action, err = s.actionService.get(actionInvocation.ActionID); err != nil {
			zap.L().Sugar().Errorf("Could not find action '%v' for action invocation '%v'. Reason: %s", actionInvocation.ActionID, actionInvocation.ID, err.Error())
			// with cascade, cannot happen
			continue
		}

		if err = s.execute(action, eventPayload); err != nil {
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

			if _, err = s.updateState(actionInvocation.ID.String(), newState); err != nil {
				zap.L().Sugar().Errorf("Could not mark action invocation '%v' as '%v'. Reason: %s", actionInvocation.ID, newState, err.Error())
			}

			if _, err = s.updateRetryCount(actionInvocation.ID.String(), newRetryCount); err != nil {
				zap.L().Sugar().Errorf("Could not update action invocation '%v' retry count to '%d'. Reason: %s", actionInvocation.ID, newRetryCount, err.Error())
			}

			msg := cause.Error()
			if _, err = s.updateMessage(actionInvocation.ID.String(), &msg); err != nil {
				zap.L().Sugar().Errorf("Could not update action invocation '%v' message. Reason: %s", actionInvocation.ID, err.Error())
			}

			continue
		}

		zap.L().Sugar().Debugf("Processed action invocation '%v' for event '%s' (%v) and action '%s' (%v)", actionInvocation.ID, event.Name, event.ID, action.Label, action.ID)
		if _, err = s.updateState(actionInvocation.ID.String(), api.ActionInvocationStateSuccess); err != nil {
			zap.L().Sugar().Errorf("Could not mark action invocation '%v' as success. Reason: %s", actionInvocation.ID, err.Error())
		}
		if _, err = s.updateMessage(actionInvocation.ID.String(), nil); err != nil {
			zap.L().Sugar().Errorf("Could not update action invocation '%v' message. Reason: %s", actionInvocation.ID, err.Error())
		}
	}

	return nil
}

func (s *actionInvocationService) execute(action *Action, eventPayloadInfo *eventPayloadInformationDto) error {
	if action == nil || eventPayloadInfo == nil {
		return errorValidationNotEmpty
	}

	var err error
	var bytes []byte

	if bytes, err = action.Payload.MarshalJSON(); err != nil {
		return newServiceError(General, err)
	}

	switch action.Type {
	case api.ActionTypeShoutrrr.Value():
		var payload actionPayloadShoutrrrDto
		if payload, err = util.UnmarshalGenericJSON[actionPayloadShoutrrrDto](bytes); err != nil {
			return newServiceError(General, err)
		}

		body := s.replaceVars(payload.Body, eventPayloadInfo)
		body = s.replaceSecrets(body)

		for _, url := range payload.Urls {
			url = s.replaceSecrets(url)
			url = s.replaceVars(url, eventPayloadInfo)
			if err = shoutrrr.Send(url, body); err != nil {
				return err
			}
		}
		break
	default:
		return newServiceError(General, errors.New("no matching action type found for invocation"))
	}

	return nil
}

func (s *actionInvocationService) replaceSecrets(str string) string {
	if str == "" {
		return str
	}

	var matches [][]string

	matches = util.ExtractBetween(str, "<SECRET>", "</SECRET>")
	var err error

	for _, match := range matches {
		var val string
		if val, err = s.secretService.getValueByKey(match[1]); err != nil {
			zap.L().Sugar().Warnf("Could not inject secret '%s'. Reason: %s", match[1], err.Error())
			continue
		}
		str = strings.ReplaceAll(str, match[0], val)
	}

	return str
}

func (s *actionInvocationService) replaceVars(str string, eventPayloadInfo *eventPayloadInformationDto) string {
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

func (s *actionInvocationService) paginate(page int, pageSize int, orderBy string, order string) ([]*ActionInvocation, error) {
	return s.repo.paginate(page, pageSize, orderBy, order)
}

func (s *actionInvocationService) get(id string) (*ActionInvocation, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	e, err := s.repo.find(id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *actionInvocationService) getByState(limit int, maxRetries int, state ...api.ActionInvocationState) ([]*ActionInvocation, error) {
	if len(state) == 0 {
		return nil, errorValidationNotEmpty
	}
	if limit <= 0 {
		return nil, errorValidationLimitGreaterZero
	}
	if maxRetries <= 0 {
		return nil, errorValidationMaxRetriesGreaterZero
	}

	return s.repo.findAllByState(limit, maxRetries, state...)
}

func (s *actionInvocationService) count() (int64, error) {
	return s.repo.count()
}

func (s *actionInvocationService) delete(id string) error {
	if id == "" {
		return errorValidationNotBlank
	}

	e, err := s.get(id)
	if err != nil {
		return err
	}

	if _, err = s.repo.delete(e.ID.String()); err != nil {
		return err
	}

	zap.L().Sugar().Infof("Deleted action '%v'", id)

	return nil
}

func (s *actionInvocationService) updateState(id string, state api.ActionInvocationState) (*ActionInvocation, error) {
	if id == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	var e *ActionInvocation
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateState(id, state); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action invocation '%v'", id)
	return e, nil
}

func (s *actionInvocationService) updateMessage(id string, message *string) (*ActionInvocation, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *ActionInvocation
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateMessage(id, message); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action invocation '%v'", id)
	return e, nil
}

func (s *actionInvocationService) updateRetryCount(id string, retryCount int) (*ActionInvocation, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *ActionInvocation
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateRetryCount(id, retryCount); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified action invocation '%v'", id)
	return e, nil
}

func (s *actionInvocationService) create(event *Event, action *Action, state api.ActionInvocationState) (*ActionInvocation, error) {
	if state == "" {
		return nil, errorValidationNotBlank
	}
	if action == nil || event == nil {
		return nil, errorValidationNotEmpty
	}

	var err error
	var e *ActionInvocation
	if e, err = s.repo.create(event.ID.String(), action.ID.String(), state); err != nil {
		return nil, err
	} else {
		zap.L().Sugar().Info("Created action invocation")
		return e, nil
	}
}

func (s *actionInvocationService) cleanStale(time time.Time, maxRetries int, state ...api.ActionInvocationState) (int64, error) {
	if len(state) == 0 {
		return 0, errorValidationNotEmpty
	}

	return s.repo.deleteByUpdatedAtBeforeAndStates(time, maxRetries, state...)
}
