package dto

// ActionPayloadShoutrrrDto payload for shoutrrr
type ActionPayloadShoutrrrDto struct {
	Body string   `binding:"required" json:"body" validate:"required"`
	Urls []string `binding:"required" json:"urls" validate:"required,min=1"`
}

// EventPayloadInformationDto general event payload
type EventPayloadInformationDto struct {
	Host           string
	Application    string
	Provider       string
	Version        string
	State          string
	StateLabel     string
	CommentAuthor  string
	CommentContent string
}

// UpdateStateReorderItem represents an item to reorder with its new sort order
type UpdateStateReorderItem struct {
	ID        string
	SortOrder int
}
