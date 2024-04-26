package api

// UpdateState state of an update
type UpdateState string

const (
	UpdateStatePending  UpdateState = "pending"
	UpdateStateApproved UpdateState = "approved"
	UpdateStateIgnored  UpdateState = "ignored"
)

func (e *UpdateState) Scan(value interface{}) error {
	*e = UpdateState(value.([]byte))
	return nil
}

func (e UpdateState) Value() string {
	return string(e)
}

// WebhookType type of webhook
type WebhookType string

const (
	WebhookTypeGeneric WebhookType = "generic"
	WebhookTypeDiun    WebhookType = "diun"
)

func (e *WebhookType) Scan(value interface{}) error {
	*e = WebhookType(value.([]byte))
	return nil
}

func (e WebhookType) Value() string {
	return string(e)
}

// EventName name of event
type EventName string

const (
	EventNameUpdateCreated         EventName = "update_created"
	EventNameUpdateUpdated         EventName = "update_updated"
	EventNameUpdateUpdatedPending  EventName = "update_updated_state_pending"
	EventNameUpdateUpdatedApproved EventName = "update_updated_state_approved"
	EventNameUpdateUpdatedIgnored  EventName = "update_updated_state_ignored"
	EventNameUpdateDeleted         EventName = "update_deleted"
)

func (e *EventName) Scan(value interface{}) error {
	*e = EventName(value.([]byte))
	return nil
}

func (e EventName) Value() string {
	return string(e)
}

// EventState name of event
type EventState string

const (
	EventStateCreated  EventState = "created"
	EventStateEnqueued EventState = "enqueued"
)

func (e *EventState) Scan(value interface{}) error {
	*e = EventState(value.([]byte))
	return nil
}

func (e EventState) Value() string {
	return string(e)
}

// ActionType state of an update
type ActionType string

const (
	ActionTypeShoutrrr ActionType = "shoutrrr"
)

func (e *ActionType) Scan(value interface{}) error {
	*e = ActionType(value.([]byte))
	return nil
}

func (e ActionType) Value() string {
	return string(e)
}

// ActionInvocationState state of an action invocation
type ActionInvocationState string

const (
	ActionInvocationStateCreated  ActionInvocationState = "created"
	ActionInvocationStateRunning  ActionInvocationState = "running"
	ActionInvocationStateRetrying ActionInvocationState = "retrying"
	ActionInvocationStateSuccess  ActionInvocationState = "success"
	ActionInvocationStateError    ActionInvocationState = "error"
)

func (e *ActionInvocationState) Scan(value interface{}) error {
	*e = ActionInvocationState(value.([]byte))
	return nil
}

func (e ActionInvocationState) Value() string {
	return string(e)
}
