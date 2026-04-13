package dto

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
	StateLabel  string
}

// UpdateStateReorderItem represents an item to reorder with its new sort order
type UpdateStateReorderItem struct {
	ID        string
	SortOrder int
}
