package api

import (
	"github.com/google/uuid"
	"time"
)

// Requests

type ModifyUpdateStateRequest struct {
	State string `json:"state" binding:"required,oneof=pending approved ignored"`
}

type ModifyWebhookLabelRequest struct {
	Label string `json:"label" binding:"required,min=1,max=255"`
}

type ModifyWebhookIgnoreHostRequest struct {
	IgnoreHost bool `json:"ignoreHost"`
}

type CreateWebhookRequest struct {
	Label      string `json:"label" binding:"required,min=1,max=255"`
	Type       string `json:"type" binding:"required,oneof=generic diun"`
	IgnoreHost bool   `json:"ignoreHost"`
}

type CreateSecretRequest struct {
	Key   string `json:"key" binding:"required,min=1"`
	Value string `json:"value" binding:"required,min=1"`
}

type CreateConstantRequest struct {
	Key   string `json:"key" binding:"required,min=1"`
	Value string `json:"value" binding:"required,min=1"`
}

type CreateActionRequest struct {
	Label            string      `json:"label" binding:"required,min=1,max=255"`
	Type             string      `json:"type" binding:"required,oneof=shoutrrr"`
	MatchEvent       *string     `json:"matchEvent"`
	MatchHost        *string     `json:"matchHost"`
	MatchApplication *string     `json:"matchApplication"`
	MatchProvider    *string     `json:"matchProvider"`
	Payload          interface{} `json:"payload"`
	Enabled          bool        `json:"enabled"`
}

type ModifySecretValueRequest struct {
	Value string `json:"value" binding:"required,min=1"`
}

type ModifyConstantValueRequest struct {
	Value string `json:"value" binding:"required,min=1"`
}

type ModifyActionLabelRequest struct {
	Label string `json:"label" binding:"required,min=1,max=255"`
}

type ModifyActionMatchEventRequest struct {
	MatchEvent *string `json:"matchEvent"`
}

type ModifyActionMatchHostRequest struct {
	MatchHost *string `json:"matchHost"`
}

type ModifyActionMatchApplicationRequest struct {
	MatchApplication *string `json:"matchApplication"`
}

type ModifyActionMatchProviderRequest struct {
	MatchProvider *string `json:"matchProvider"`
}

type ModifyActionTypeAndPayloadRequest struct {
	Type    ActionType  `json:"type" binding:"required,oneof=shoutrrr"`
	Payload interface{} `json:"payload" binding:"required"`
}

type ModifyActionEnabledRequest struct {
	Enabled bool `json:"enabled"`
}

type TestActionRequest struct {
	Application string `json:"application" binding:"required,min=1"`
	Provider    string `json:"provider" binding:"required,min=1"`
	Host        string `json:"host" binding:"required,min=1"`
	Version     string `json:"version" binding:"required,min=1"`
	State       string `json:"state" binding:"required,min=1"`
}

type PaginateUpdateRequest struct {
	PageSize   int    `form:"pageSize,default=5" binding:"numeric,gte=1"`
	Page       int    `form:"page,default=1" binding:"numeric,gte=1"`
	Order      string `form:"order,default=desc" binding:"oneof=asc desc"`
	OrderBy    string `form:"orderBy,default=updated_at" binding:"oneof=id application provider host version created_at updated_at"`
	SearchTerm string `form:"searchTerm"`
	SearchIn   string `form:"searchIn,default=application" binding:"oneof=application provider host version"`
}

type PaginateWebhookRequest struct {
	PageSize int    `form:"pageSize,default=5" binding:"numeric,gte=1"`
	Page     int    `form:"page,default=1" binding:"numeric,gte=1"`
	Order    string `form:"order,default=asc" binding:"oneof=asc desc"`
	OrderBy  string `form:"orderBy,default=label" binding:"oneof=id label type created_at updated_at"`
}

type PaginateActionRequest struct {
	PageSize int    `form:"pageSize,default=5" binding:"numeric,gte=1"`
	Page     int    `form:"page,default=1" binding:"numeric,gte=1"`
	Order    string `form:"order,default=asc" binding:"oneof=asc desc"`
	OrderBy  string `form:"orderBy,default=label" binding:"oneof=id label type created_at updated_at"`
}

type PaginateActionInvocationRequest struct {
	PageSize int    `form:"pageSize,default=5" binding:"numeric,gte=1"`
	Page     int    `form:"page,default=1" binding:"numeric,gte=1"`
	Order    string `form:"order,default=desc" binding:"oneof=asc desc"`
	OrderBy  string `form:"orderBy,default=created_at" binding:"oneof=id state retry_count created_at updated_at"`
}

type WebhookGenericRequest struct {
	Application string      `json:"application" binding:"required,min=1"`
	Provider    string      `json:"provider"`
	Host        string      `json:"host" binding:"required,min=1"`
	Version     string      `json:"version" binding:"required,min=1"`
	Metadata    interface{} `json:"metadata"`
}

type WebhookDiunMetadataRequest struct {
	Command   string `json:"ctn_command"`
	CreatedAt string `json:"ctn_createdat"`
	Id        string `json:"ctn_id"`
	Names     string `json:"ctn_names"`
	Size      string `json:"ctn_size"`
	State     string `json:"ctn_state"`
	Status    string `json:"ctn_status"`
}

type WebhookDiunRequest struct {
	DiunVersion string                     `json:"diun_version" binding:"required,min=1"`
	Hostname    string                     `json:"hostname" binding:"required,min=1"`
	Status      string                     `json:"status" binding:"required,min=1"`
	Provider    string                     `json:"provider" binding:"required,min=1"`
	Image       string                     `json:"image" binding:"required,min=1"`
	HubLink     string                     `json:"hub_link"`
	MimeType    string                     `json:"mime_type" binding:"required,min=1"`
	Digest      string                     `json:"digest" binding:"required,min=1"`
	Created     string                     `json:"created" binding:"required,min=1"`
	Platform    string                     `json:"platform" binding:"required,min=1"`
	Metadata    WebhookDiunMetadataRequest `json:"metadata"`
}

type EventWindowRequest struct {
	Size    int    `form:"size,default=10" binding:"numeric,gte=1"`
	Skip    int    `form:"skip,default=0" binding:"numeric"`
	Order   string `form:"order,default=desc" binding:"oneof=asc desc"`
	OrderBy string `form:"orderBy,default=created_at" binding:"oneof=id name created_at updated_at"`
}

// Responses

type Response struct {
}

type DataResponse struct {
	Response
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status string `json:"status,omitempty"`
	DataResponse
}

func NewDataResponseWithPayload(payload interface{}) *DataResponse {
	e := new(DataResponse)
	e.Data = payload
	return e
}

func NewErrorResponseWithStatusAndMessage(status string, message string) *ErrorResponse {
	e := new(ErrorResponse)
	e.Status = status
	e.Message = message
	return e
}

type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

func NewHealthResponse(b bool) *HealthResponse {
	r := new(HealthResponse)
	r.Healthy = b
	return r
}

type InfoResponse struct {
	Version  string `json:"version"`
	Name     string `json:"name"`
	TimeZone string `json:"timeZone"`
}

func NewInfoResponse(name string, version string, tz string) *InfoResponse {
	r := new(InfoResponse)
	r.Name = name
	r.Version = version
	r.TimeZone = tz
	return r
}

type UpdateResponse struct {
	ID          uuid.UUID   `json:"id"`
	Application string      `json:"application"`
	Provider    string      `json:"provider"`
	Host        string      `json:"host"`
	Version     string      `json:"version"`
	State       string      `json:"state"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

type UpdateSingleResponse struct {
	Data UpdateResponse `json:"data"`
}

func NewUpdateSingleResponse(id uuid.UUID, application string, provider string, host string, version string, state string, createdAt time.Time, updatedAt time.Time, metadata interface{}) *UpdateSingleResponse {
	e := new(UpdateSingleResponse)
	e.Data.ID = id
	e.Data.Application = application
	e.Data.Provider = provider
	e.Data.Host = host
	e.Data.Version = version
	e.Data.State = state
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	e.Data.Metadata = metadata
	return e
}

type UpdatePageResponse struct {
	Content       []*UpdateResponse `json:"content"`
	Page          int               `json:"page"`
	PageSize      int               `json:"pageSize"`
	OrderBy       string            `json:"orderBy"`
	Order         string            `json:"order"`
	TotalElements int64             `json:"totalElements"`
	TotalPages    int64             `json:"totalPages"`
}

type UpdateDataPageResponse struct {
	Data *UpdatePageResponse `json:"data"`
}

func NewUpdatePageResponse(content []*UpdateResponse, page int, pageSize int, orderBy string, order string, totalElements int64, totalPages int64) *UpdatePageResponse {
	e := new(UpdatePageResponse)
	e.Content = content
	e.Page = page
	e.PageSize = pageSize
	e.OrderBy = orderBy
	e.Order = order
	e.TotalElements = totalElements
	e.TotalPages = totalPages
	return e
}

type WebhookResponse struct {
	ID         uuid.UUID `json:"id"`
	Label      string    `json:"label"`
	Type       string    `json:"type"`
	IgnoreHost bool      `json:"ignoreHost"`
	Token      string    `json:"token,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type WebhookSingleResponse struct {
	Data WebhookResponse `json:"data"`
}

func NewWebhookSingleResponse(id uuid.UUID, label string, t string, ignoreHost bool, token string, createdAt time.Time, updatedAt time.Time) *WebhookSingleResponse {
	e := new(WebhookSingleResponse)
	e.Data.ID = id
	e.Data.Label = label
	e.Data.Type = t
	e.Data.IgnoreHost = ignoreHost
	e.Data.Token = token
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type WebhookPageResponse struct {
	Content       []*WebhookResponse `json:"content"`
	Page          int                `json:"page"`
	PageSize      int                `json:"pageSize"`
	OrderBy       string             `json:"orderBy"`
	Order         string             `json:"order"`
	TotalElements int64              `json:"totalElements"`
	TotalPages    int64              `json:"totalPages"`
}

func NewWebhookPageResponse(content []*WebhookResponse, page int, pageSize int, orderBy string, order string, totalElements int64, totalPages int64) *WebhookPageResponse {
	e := new(WebhookPageResponse)
	e.Content = content
	e.Page = page
	e.PageSize = pageSize
	e.OrderBy = orderBy
	e.Order = order
	e.TotalElements = totalElements
	e.TotalPages = totalPages
	return e
}

type EventResponse struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Payload   interface{} `json:"payload,omitempty"`
}

type EventSingleResponse struct {
	Data EventResponse `json:"data"`
}

func NewEventSingleResponse(id uuid.UUID, name string, createdAt time.Time, updatedAt time.Time, payload interface{}) *EventSingleResponse {
	e := new(EventSingleResponse)
	e.Data.ID = id
	e.Data.Name = name
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	e.Data.Payload = payload
	return e
}

type EventWindowResponse struct {
	Content []*EventResponse `json:"content"`
	Size    int              `json:"size"`
	Skip    int              `json:"skip"`
	HasNext bool             `json:"hasNext"`
	OrderBy string           `json:"orderBy"`
	Order   string           `json:"order"`
}

type EventPayloadUpdateCreatedDto struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Application string    `json:"application,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	Host        string    `json:"host,omitempty"`
	Version     string    `json:"version,omitempty"`
	State       string    `json:"state,omitempty"`
}

type EventPayloadUpdateUpdatedDto struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Application  string    `json:"application,omitempty"`
	Provider     string    `json:"provider,omitempty"`
	Host         string    `json:"host,omitempty"`
	VersionPrior string    `json:"versionPrior,omitempty"`
	Version      string    `json:"version,omitempty"`
	StatePrior   string    `json:"statePrior,omitempty"`
	State        string    `json:"state,omitempty"`
}

type EventPayloadUpdateDeletedDto struct {
	Application string `json:"application,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Host        string `json:"host,omitempty"`
	Version     string `json:"version,omitempty"`
	State       string `json:"state,omitempty"`
}

func NewEventWindowResponse(content []*EventResponse, size int, skip int, orderBy string, order string, hasNext bool) *EventWindowResponse {
	e := new(EventWindowResponse)
	e.Content = content
	e.Size = size
	e.Skip = skip
	e.HasNext = hasNext
	e.OrderBy = orderBy
	e.Order = order
	return e
}

type SecretResponse struct {
	ID        uuid.UUID `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SecretSingleResponse struct {
	Data SecretResponse `json:"data"`
}

func NewSecretSingleResponse(id uuid.UUID, key string, value string, createdAt time.Time, updatedAt time.Time) *SecretSingleResponse {
	e := new(SecretSingleResponse)
	e.Data.ID = id
	e.Data.Key = key
	e.Data.Value = value
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type SecretPageResponse struct {
	Content []*SecretResponse `json:"content"`
}

type SecretDataPageResponse struct {
	Data *SecretPageResponse `json:"data"`
}

func NewSecretPageResponse(content []*SecretResponse) *SecretPageResponse {
	e := new(SecretPageResponse)
	e.Content = content
	return e
}

type ConstantResponse struct {
	ID        uuid.UUID `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ConstantSingleResponse struct {
	Data ConstantResponse `json:"data"`
}

func NewConstantSingleResponse(id uuid.UUID, key string, value string, createdAt time.Time, updatedAt time.Time) *ConstantSingleResponse {
	e := new(ConstantSingleResponse)
	e.Data.ID = id
	e.Data.Key = key
	e.Data.Value = value
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type ConstantPageResponse struct {
	Content []*ConstantResponse `json:"content"`
}

type ConstantDataPageResponse struct {
	Data *ConstantPageResponse `json:"data"`
}

func NewConstantPageResponse(content []*ConstantResponse) *ConstantPageResponse {
	e := new(ConstantPageResponse)
	e.Content = content
	return e
}

type ActionResponse struct {
	ID               uuid.UUID   `json:"id"`
	Label            string      `json:"label"`
	Type             string      `json:"type"`
	MatchEvent       *string     `json:"matchEvent,omitempty"`
	MatchHost        *string     `json:"matchHost,omitempty"`
	MatchApplication *string     `json:"matchApplication,omitempty"`
	MatchProvider    *string     `json:"matchProvider,omitempty"`
	Payload          interface{} `json:"payload,omitempty"`
	Enabled          bool        `json:"enabled"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
}

type ActionSingleResponse struct {
	Data ActionResponse `json:"data"`
}

func NewActionSingleResponse(id uuid.UUID, label string, t string, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool, createdAt time.Time, updatedAt time.Time) *ActionSingleResponse {
	e := new(ActionSingleResponse)
	e.Data.ID = id
	e.Data.Label = label
	e.Data.Type = t
	e.Data.MatchEvent = matchEvent
	e.Data.MatchHost = matchHost
	e.Data.MatchApplication = matchApplication
	e.Data.MatchProvider = matchProvider
	e.Data.Payload = payload
	e.Data.Enabled = enabled
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type ActionPageResponse struct {
	Content       []*ActionResponse `json:"content"`
	Page          int               `json:"page"`
	PageSize      int               `json:"pageSize"`
	OrderBy       string            `json:"orderBy"`
	Order         string            `json:"order"`
	TotalElements int64             `json:"totalElements"`
	TotalPages    int64             `json:"totalPages"`
}

func NewActionPageResponse(content []*ActionResponse, page int, pageSize int, orderBy string, order string, totalElements int64, totalPages int64) *ActionPageResponse {
	e := new(ActionPageResponse)
	e.Content = content
	e.Page = page
	e.PageSize = pageSize
	e.OrderBy = orderBy
	e.Order = order
	e.TotalElements = totalElements
	e.TotalPages = totalPages
	return e
}

type ActionTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ActionTestSingleResponse struct {
	Data ActionTestResponse `json:"data"`
}

func NewActionTestSingleResponse(success bool, message string) *ActionTestSingleResponse {
	e := new(ActionTestSingleResponse)
	e.Data.Success = success
	e.Data.Message = message
	return e
}

type ActionInvocationResponse struct {
	ID         uuid.UUID `json:"id"`
	RetryCount int       `json:"retryCount"`
	State      string    `json:"state"`
	Message    *string   `json:"message,omitempty"`
	ActionID   string    `json:"actionId"`
	EventID    string    `json:"eventId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type ActionInvocationSingleResponse struct {
	Data ActionInvocationResponse `json:"data"`
}

func NewActionInvocationSingleResponse(id uuid.UUID, retryCount int, state string, message *string, actionId string, eventId string, createdAt time.Time, updatedAt time.Time) *ActionInvocationSingleResponse {
	e := new(ActionInvocationSingleResponse)
	e.Data.ID = id
	e.Data.RetryCount = retryCount
	e.Data.State = state
	e.Data.Message = message
	e.Data.ActionID = actionId
	e.Data.EventID = eventId
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type ActionInvocationPageResponse struct {
	Content       []*ActionInvocationResponse `json:"content"`
	Page          int                         `json:"page"`
	PageSize      int                         `json:"pageSize"`
	OrderBy       string                      `json:"orderBy"`
	Order         string                      `json:"order"`
	TotalElements int64                       `json:"totalElements"`
	TotalPages    int64                       `json:"totalPages"`
}

func NewActionInvocationPageResponse(content []*ActionInvocationResponse, page int, pageSize int, orderBy string, order string, totalElements int64, totalPages int64) *ActionInvocationPageResponse {
	e := new(ActionInvocationPageResponse)
	e.Content = content
	e.Page = page
	e.PageSize = pageSize
	e.OrderBy = orderBy
	e.Order = order
	e.TotalElements = totalElements
	e.TotalPages = totalPages
	return e
}
