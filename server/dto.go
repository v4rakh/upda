package server

// DTOs

type actionPayloadShoutrrrDto struct {
	Body string   `json:"body" binding:"required" validate:"required"`
	Urls []string `json:"urls" binding:"required" validate:"required,min=1"`
}

type eventPayloadInformationDto struct {
	Host        string
	Application string
	Provider    string
	Version     string
}
