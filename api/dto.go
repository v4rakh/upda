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
	Order    string `form:"order,default=desc" binding:"oneof=asc desc"`
	OrderBy  string `form:"orderBy,default=updated_at" binding:"oneof=id label type created_at updated_at"`
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
	State     string      `json:"state"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Payload   interface{} `json:"payload,omitempty"`
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

type EventPayloadWebhookCreatedDto struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Label      string    `json:"label,omitempty"`
	Type       string    `json:"type,omitempty"`
	IgnoreHost bool      `json:"ignoreHost"`
}

type EventPayloadWebhookUpdatedDto struct {
	ID              uuid.UUID `json:"id,omitempty"`
	LabelPrior      string    `json:"labelPrior,omitempty"`
	Label           string    `json:"label,omitempty"`
	IgnoreHostPrior bool      `json:"ignoreHostPrior"`
	IgnoreHost      bool      `json:"ignoreHost"`
	Type            string    `json:"type,omitempty"`
}

type EventPayloadWebhookDeletedDto struct {
	Label      string `json:"label,omitempty"`
	Type       string `json:"type,omitempty"`
	IgnoreHost bool   `json:"ignoreHost"`
}
