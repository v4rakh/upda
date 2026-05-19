package api

import (
	"time"
)

// Requests

// json/body

type ModifyUpdateStateRequest struct {
	State string `binding:"required,min=1,max=50" json:"state"`
}

type ModifyWebhookLabelRequest struct {
	Label string `binding:"required,min=1,max=255" json:"label"`
}

type ModifyWebhookIgnoreHostRequest struct {
	IgnoreHost bool `json:"ignoreHost"`
}

type ModifyWebhookIgnoreHostReplacementRequest struct {
	IgnoreHostReplacement string `binding:"required,min=1,max=255" json:"ignoreHostReplacement"`
}

type CreateWebhookRequest struct {
	Label                 string `binding:"required,min=1,max=255"      json:"label"`
	Type                  string `binding:"required,oneof=generic diun" json:"type"`
	IgnoreHost            bool   `json:"ignoreHost"`
	IgnoreHostReplacement string `binding:"required,min=1,max=255"      json:"ignoreHostReplacement"`
}

type CreateSecretRequest struct {
	Key   string `binding:"required,min=1" json:"key"`
	Value string `binding:"required,min=1" json:"value"`
}

type CreateConstantRequest struct {
	Key   string `binding:"required,min=1" json:"key"`
	Value string `binding:"required,min=1" json:"value"`
}

type CreateActionRequest struct {
	Label            string      `binding:"required,min=1,max=255"  json:"label"`
	Type             string      `binding:"required,oneof=shoutrrr" json:"type"`
	MatchEvent       *string     `json:"matchEvent"`
	MatchHost        *string     `json:"matchHost"`
	MatchApplication *string     `json:"matchApplication"`
	MatchProvider    *string     `json:"matchProvider"`
	Payload          interface{} `json:"payload"`
	Enabled          bool        `json:"enabled"`
}
type CreateFilterPresetRequest struct {
	Type       string  `binding:"required,oneof=update"  json:"type"`
	Label      string  `binding:"required,min=1,max=255" json:"label"`
	Parameters string  `binding:"required"               json:"parameters"`
	Color      *string `json:"color"`
}

type CreateCommentRequest struct {
	Content string `binding:"required,min=1" json:"content"`
}

type ModifySecretValueRequest struct {
	Value string `binding:"required,min=1" json:"value"`
}

type ModifyConstantValueRequest struct {
	Value string `binding:"required,min=1" json:"value"`
}

type ModifyActionLabelRequest struct {
	Label string `binding:"required,min=1,max=255" json:"label"`
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
	Type    string      `binding:"required,oneof=shoutrrr" json:"type"`
	Payload interface{} `binding:"required"                json:"payload"`
}

type ModifyActionEnabledRequest struct {
	Enabled bool `json:"enabled"`
}

type ModifyCommentContentRequest struct {
	Content string `binding:"required,min=1" json:"content"`
}

type TestActionRequest struct {
	Application string `binding:"required,min=1" json:"application"`
	Provider    string `binding:"required,min=1" json:"provider"`
	Host        string `binding:"required,min=1" json:"host"`
	Version     string `binding:"required,min=1" json:"version"`
	State       string `binding:"required,min=1" json:"state"`
}

type WebhookGenericRequest struct {
	Application string      `binding:"required,min=1" json:"application"`
	Provider    string      `json:"provider"`
	Host        string      `binding:"required,min=1" json:"host"`
	Version     string      `binding:"required,min=1" json:"version"`
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
	DiunVersion string                     `binding:"required,min=1" json:"diun_version"`
	Hostname    string                     `binding:"required,min=1" json:"hostname"`
	Status      string                     `binding:"required,min=1" json:"status"`
	Provider    string                     `binding:"required,min=1" json:"provider"`
	Image       string                     `binding:"required,min=1" json:"image"`
	HubLink     string                     `json:"hub_link"`
	MimeType    string                     `binding:"required,min=1" json:"mime_type"`
	Digest      string                     `binding:"required,min=1" json:"digest"`
	Created     string                     `binding:"required,min=1" json:"created"`
	Platform    string                     `binding:"required,min=1" json:"platform"`
	Metadata    WebhookDiunMetadataRequest `json:"metadata"`
}

// query parameters

type PaginateUpdateRequest struct {
	PageSize   int    `binding:"numeric,gte=1"                                                    form:"pageSize,default=5"`
	Page       int    `binding:"numeric,gte=1"                                                    form:"page,default=1"`
	Order      string `binding:"oneof=asc desc"                                                   form:"order,default=desc"`
	OrderBy    string `binding:"oneof=id application provider host version created_at updated_at" form:"orderBy,default=updated_at"`
	SearchTerm string `form:"searchTerm"`
	SearchIn   string `binding:"oneof=application provider host version"                          form:"searchIn,default=application"`
}

type PaginateWebhookRequest struct {
	PageSize int    `binding:"numeric,gte=1"                             form:"pageSize,default=5"`
	Page     int    `binding:"numeric,gte=1"                             form:"page,default=1"`
	Order    string `binding:"oneof=asc desc"                            form:"order,default=asc"`
	OrderBy  string `binding:"oneof=id label type created_at updated_at" form:"orderBy,default=label"`
}

type PaginateActionRequest struct {
	PageSize int    `binding:"numeric,gte=1"                             form:"pageSize,default=5"`
	Page     int    `binding:"numeric,gte=1"                             form:"page,default=1"`
	Order    string `binding:"oneof=asc desc"                            form:"order,default=asc"`
	OrderBy  string `binding:"oneof=id label type created_at updated_at" form:"orderBy,default=label"`
}

type PaginateActionInvocationRequest struct {
	PageSize int    `binding:"numeric,gte=1"                                    form:"pageSize,default=5"`
	Page     int    `binding:"numeric,gte=1"                                    form:"page,default=1"`
	Order    string `binding:"oneof=asc desc"                                   form:"order,default=desc"`
	OrderBy  string `binding:"oneof=id state retry_count created_at updated_at" form:"orderBy,default=created_at"`
}

type EventWindowRequest struct {
	Size     int     `binding:"numeric,gte=1"                       form:"size,default=10"`
	Skip     int     `binding:"numeric"                             form:"skip,default=0"`
	Order    string  `binding:"oneof=asc desc"                      form:"order,default=desc"`
	OrderBy  string  `binding:"oneof=id name created_at updated_at" form:"orderBy,default=created_at"`
	UpdateID *string `form:"updateId"`
}

type PaginateCommentRequest struct {
	PageSize int `binding:"numeric,gte=1" form:"pageSize,default=5"`
	Page     int `binding:"numeric,gte=1" form:"page,default=1"`
}

// uri parameters

type FilterPresetUriRequest struct {
	Type string `binding:"required,oneof=update" uri:"type"`
}

type IDUriRequest struct {
	ID string `binding:"required,uuid4" uri:"id"`
}

type UpdateIDUriRequest struct {
	ID string `binding:"required,uuid4" uri:"updateId"`
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
	ID          string      `json:"id"`
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

func NewUpdateSingleResponse(id string, application string, provider string, host string, version string, state string, createdAt time.Time, updatedAt time.Time, metadata interface{}) *UpdateSingleResponse {
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
	ID                    string    `json:"id"`
	Label                 string    `json:"label"`
	Type                  string    `json:"type"`
	IgnoreHost            bool      `json:"ignoreHost"`
	IgnoreHostReplacement string    `json:"ignoreHostReplacement"`
	Token                 string    `json:"token,omitempty"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type WebhookSingleResponse struct {
	Data WebhookResponse `json:"data"`
}

func NewWebhookSingleResponse(id string, label string, t string, ignoreHost bool, ignoreHostReplacement string, token string, createdAt time.Time, updatedAt time.Time) *WebhookSingleResponse {
	e := new(WebhookSingleResponse)
	e.Data.ID = id
	e.Data.Label = label
	e.Data.Type = t
	e.Data.IgnoreHost = ignoreHost
	e.Data.IgnoreHostReplacement = ignoreHostReplacement
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
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Payload   interface{} `json:"payload,omitempty"`
}

type EventSingleResponse struct {
	Data EventResponse `json:"data"`
}

func NewEventSingleResponse(id string, name string, createdAt time.Time, updatedAt time.Time, payload interface{}) *EventSingleResponse {
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
	ID          string `json:"id,omitempty"`
	Application string `json:"application,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Host        string `json:"host,omitempty"`
	Version     string `json:"version,omitempty"`
	State       string `json:"state,omitempty"`
	StateLabel  string `json:"stateLabel,omitempty"`
}

type EventPayloadUpdateUpdatedDto struct {
	ID              string `json:"id,omitempty"`
	Application     string `json:"application,omitempty"`
	Provider        string `json:"provider,omitempty"`
	Host            string `json:"host,omitempty"`
	VersionPrior    string `json:"versionPrior,omitempty"`
	Version         string `json:"version,omitempty"`
	StatePrior      string `json:"statePrior,omitempty"`
	StatePriorLabel string `json:"statePriorLabel,omitempty"`
	State           string `json:"state,omitempty"`
	StateLabel      string `json:"stateLabel,omitempty"`
}

type EventPayloadUpdateDeletedDto struct {
	Application string `json:"application,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Host        string `json:"host,omitempty"`
	Version     string `json:"version,omitempty"`
	State       string `json:"state,omitempty"`
	StateLabel  string `json:"stateLabel,omitempty"`
}

type EventPayloadCommentCreatedDto struct {
	CommentID   string `json:"commentId,omitempty"`
	Author      string `json:"author,omitempty"`
	Content     string `json:"content,omitempty"`
	UpdateID    string `json:"updateId,omitempty"`
	Application string `json:"application,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Host        string `json:"host,omitempty"`
	Version     string `json:"version,omitempty"`
	State       string `json:"state,omitempty"`
	StateLabel  string `json:"stateLabel,omitempty"`
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
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SecretSingleResponse struct {
	Data SecretResponse `json:"data"`
}

func NewSecretSingleResponse(id string, key string, value string, createdAt time.Time, updatedAt time.Time) *SecretSingleResponse {
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
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ConstantSingleResponse struct {
	Data ConstantResponse `json:"data"`
}

func NewConstantSingleResponse(id string, key string, value string, createdAt time.Time, updatedAt time.Time) *ConstantSingleResponse {
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
	ID               string      `json:"id"`
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

func NewActionSingleResponse(id string, label string, t string, matchEvent *string, matchHost *string, matchApplication *string, matchProvider *string, payload interface{}, enabled bool, createdAt time.Time, updatedAt time.Time) *ActionSingleResponse {
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
	ID         string    `json:"id"`
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

func NewActionInvocationSingleResponse(id string, retryCount int, state string, message *string, actionId string, eventId string, createdAt time.Time, updatedAt time.Time) *ActionInvocationSingleResponse {
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

type FilterPresetResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Label      string    `json:"label"`
	Parameters string    `json:"parameters"`
	Color      *string   `json:"color,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type FilterPresetSingleResponse struct {
	Data FilterPresetResponse `json:"data"`
}

func NewFilterPresetSingleResponse(id string, t string, label string, parameters string, color *string, createdAt time.Time, updatedAt time.Time) *FilterPresetSingleResponse {
	e := new(FilterPresetSingleResponse)
	e.Data.ID = id
	e.Data.Type = t
	e.Data.Label = label
	e.Data.Parameters = parameters
	e.Data.Color = color
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type FilterPresetPageResponse struct {
	Content []*FilterPresetResponse `json:"content"`
}

type FilterPresetDataPageResponse struct {
	Data *FilterPresetDataPageResponse `json:"data"`
}

func NewFilterPresetPageResponse(content []*FilterPresetResponse) *FilterPresetPageResponse {
	e := new(FilterPresetPageResponse)
	e.Content = content
	return e
}

type CommentResponse struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	UpdateID  string    `json:"updateId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommentSingleResponse struct {
	Data CommentResponse `json:"data"`
}

func NewCommentSingleResponse(id string, author string, content string, updateId string, createdAt time.Time, updatedAt time.Time) *CommentSingleResponse {
	e := new(CommentSingleResponse)
	e.Data.ID = id
	e.Data.Author = author
	e.Data.Content = content
	e.Data.UpdateID = updateId
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type CommentPageResponse struct {
	Content       []*CommentResponse `json:"content"`
	Page          int                `json:"page"`
	PageSize      int                `json:"pageSize"`
	TotalElements int64              `json:"totalElements"`
	TotalPages    int64              `json:"totalPages"`
}

type CommentDataPageResponse struct {
	Data *CommentDataPageResponse `json:"data"`
}

func NewCommentPageResponse(content []*CommentResponse, page int, pageSize int, totalElements int64, totalPages int64) *CommentPageResponse {
	e := new(CommentPageResponse)
	e.Content = content
	e.Page = page
	e.PageSize = pageSize
	e.TotalElements = totalElements
	e.TotalPages = totalPages
	return e
}

// Update State Definition DTOs

type CreateUpdateStateDefinitionRequest struct {
	Name             string  `binding:"required,min=1,max=50,alphanum" json:"name"`
	Label            string  `binding:"required,min=1,max=100"         json:"label"`
	Color            string  `binding:"required,min=1,max=50"          json:"color"`
	Icon             string  `binding:"required,min=1,max=100"         json:"icon"`
	Description      *string `json:"description"`
	IsInitial        bool    `json:"isInitial"`
	SkipOnNewVersion bool    `json:"skipOnNewVersion"`
}

type ModifyUpdateStateDefinitionRequest struct {
	Name             string  `binding:"required,min=1,max=50,alphanum" json:"name"`
	Label            string  `binding:"required,min=1,max=100"         json:"label"`
	Color            string  `binding:"required,min=1,max=50"          json:"color"`
	Icon             string  `binding:"required,min=1,max=100"         json:"icon"`
	Description      *string `json:"description"`
	IsInitial        bool    `json:"isInitial"`
	SkipOnNewVersion bool    `json:"skipOnNewVersion"`
	SortOrder        int     `json:"sortOrder"`
}

type ReorderUpdateStateDefinitionItem struct {
	ID        string `binding:"required,uuid4" json:"id"`
	SortOrder int    `binding:"min=0"          json:"sortOrder"`
}

type ReorderUpdateStateDefinitionsRequest struct {
	Items []ReorderUpdateStateDefinitionItem `binding:"required,min=1,dive" json:"items"`
}

type UpdateStateDefinitionResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Label            string    `json:"label"`
	Color            string    `json:"color"`
	Icon             string    `json:"icon"`
	Description      *string   `json:"description,omitempty"`
	IsInitial        bool      `json:"isInitial"`
	SkipOnNewVersion bool      `json:"skipOnNewVersion"`
	SortOrder        int       `json:"sortOrder"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type UpdateStateDefinitionSingleResponse struct {
	Data UpdateStateDefinitionResponse `json:"data"`
}

func NewUpdateStateDefinitionSingleResponse(id string, name string, label string, color string, icon string, description *string, isInitial bool, skipOnNewVersion bool, sortOrder int, createdAt time.Time, updatedAt time.Time) *UpdateStateDefinitionSingleResponse {
	e := new(UpdateStateDefinitionSingleResponse)
	e.Data.ID = id
	e.Data.Name = name
	e.Data.Label = label
	e.Data.Color = color
	e.Data.Icon = icon
	e.Data.Description = description
	e.Data.IsInitial = isInitial
	e.Data.SkipOnNewVersion = skipOnNewVersion
	e.Data.SortOrder = sortOrder
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type UpdateStateDefinitionPageResponse struct {
	Content []*UpdateStateDefinitionResponse `json:"content"`
}

type UpdateStateDefinitionDataPageResponse struct {
	Data *UpdateStateDefinitionPageResponse `json:"data"`
}

func NewUpdateStateDefinitionPageResponse(content []*UpdateStateDefinitionResponse) *UpdateStateDefinitionPageResponse {
	e := new(UpdateStateDefinitionPageResponse)
	e.Content = content
	return e
}

// Update State Transition DTOs

type CreateUpdateStateTransitionRequest struct {
	FromStateId string `binding:"required,uuid4" json:"fromStateId"`
	ToStateId   string `binding:"required,uuid4" json:"toStateId"`
}

type StateIdUriRequest struct {
	StateID string `binding:"required,uuid4" uri:"stateId"`
}

type UpdateStateTransitionResponse struct {
	ID        string                        `json:"id"`
	FromState UpdateStateDefinitionResponse `json:"fromState"`
	ToState   UpdateStateDefinitionResponse `json:"toState"`
	CreatedAt time.Time                     `json:"createdAt"`
	UpdatedAt time.Time                     `json:"updatedAt"`
}

type UpdateStateTransitionSingleResponse struct {
	Data UpdateStateTransitionResponse `json:"data"`
}

func NewUpdateStateTransitionSingleResponse(id string, fromState UpdateStateDefinitionResponse, toState UpdateStateDefinitionResponse, createdAt time.Time, updatedAt time.Time) *UpdateStateTransitionSingleResponse {
	e := new(UpdateStateTransitionSingleResponse)
	e.Data.ID = id
	e.Data.FromState = fromState
	e.Data.ToState = toState
	e.Data.CreatedAt = createdAt
	e.Data.UpdatedAt = updatedAt
	return e
}

type UpdateStateTransitionPageResponse struct {
	Content []*UpdateStateTransitionResponse `json:"content"`
}

type UpdateStateTransitionDataPageResponse struct {
	Data *UpdateStateTransitionPageResponse `json:"data"`
}

func NewUpdateStateTransitionPageResponse(content []*UpdateStateTransitionResponse) *UpdateStateTransitionPageResponse {
	e := new(UpdateStateTransitionPageResponse)
	e.Content = content
	return e
}
