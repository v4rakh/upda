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
	EventNameUpdateCreated            EventName = "update_created"
	EventNameUpdateUpdated            EventName = "update_updated"
	EventNameUpdateUpdatedPending     EventName = "update_updated_state_pending"
	EventNameUpdateUpdatedApproved    EventName = "update_updated_state_approved"
	EventNameUpdateUpdatedIgnored     EventName = "update_updated_state_ignored"
	EventNameUpdateDeleted            EventName = "update_deleted"
	EventNameWebhookCreated           EventName = "webhook_created"
	EventNameWebhookUpdatedLabel      EventName = "webhook_updated_label"
	EventNameWebhookUpdatedIgnoreHost EventName = "webhook_updated_ignore_host"
	EventNameWebhookDeleted           EventName = "webhook_deleted"
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
	EventStateCreated EventState = "created"
)

func (e *EventState) Scan(value interface{}) error {
	*e = EventState(value.([]byte))
	return nil
}

func (e EventState) Value() string {
	return string(e)
}
