package api

// UpdateState state of a model.Update
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

// WebhookType type of model.Webhook
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

// EventName name of a model.Event
type EventName string

const (
	EventNameUpdateCreated        EventName = "update_created"
	EventNameUpdateUpdated        EventName = "update_updated"
	EventNameUpdateUpdatedState   EventName = "update_updated_state"
	EventNameUpdateUpdatedVersion EventName = "update_updated_version"
	EventNameUpdateDeleted        EventName = "update_deleted"
)

func (e *EventName) Scan(value interface{}) error {
	*e = EventName(value.([]byte))
	return nil
}

func (e EventName) Value() string {
	return string(e)
}

// EventState state of a model.Event
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

// ActionType state of a model.Action
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

// ActionInvocationState state of a model.ActionInvocation
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

// FilterPresetType the type for the model.FilterPreset
type FilterPresetType string

const (
	FilterPresetTypeUpdate FilterPresetType = "update"
)

func (e *FilterPresetType) Scan(value interface{}) error {
	*e = FilterPresetType(value.([]byte))
	return nil
}

func (e FilterPresetType) Value() string {
	return string(e)
}

// FromVariadicToStr converts variadic notation to string array if type is of string
func FromVariadicToStr[T ~string](s ...T) []string {
	arr := make([]string, 0, len(s))
	for _, i := range s {
		arr = append(arr, string(i))
	}
	return arr
}
