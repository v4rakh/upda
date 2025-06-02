package dto

// ContextSession represents context which gets loaded in each route the middleware is attached
type ContextSession struct {
	User string
}

// ActionPayloadShoutrrrDto payload for shoutrrr
type ActionPayloadShoutrrrDto struct {
	Body string   `json:"body" binding:"required" validate:"required"`
	Urls []string `json:"urls" binding:"required" validate:"required,min=1"`
}

// EventPayloadInformationDto general event payload
type EventPayloadInformationDto struct {
	Host        string
	Application string
	Provider    string
	Version     string
	State       string
}
