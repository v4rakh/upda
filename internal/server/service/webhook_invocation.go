package service

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"strings"
)

const (
	providerDiun          = "oci"
	hostIgnoreReplacement = "global"
)

type WebhookInvocationService struct {
	updateService  *UpdateService
	webhookService *WebhookService
	webhookConfig  *config.Webhook
}

func NewWebhookInvocationService(w *WebhookService, u *UpdateService, c *config.Webhook) *WebhookInvocationService {
	return &WebhookInvocationService{
		updateService:  u,
		webhookService: w,
		webhookConfig:  c,
	}
}

func (s *WebhookInvocationService) ExecuteGeneric(id string, token string, req api.WebhookGenericRequest) error {
	if id == "" || token == "" {
		return service_error.ErrValidationNotBlank
	}

	var e *model.Webhook
	var err error

	if e, err = s.webhookService.Get(id); err != nil {
		return service_error.ErrResourceNotFound
	}

	if e.Token != token {
		return service_error.ErrResourceAccessDenied
	}

	host := req.Host
	if e.IgnoreHost {
		host = hostIgnoreReplacement
	}

	var provider string
	if req.Provider == "" {
		provider = e.Label
	} else {
		provider = req.Provider
	}

	if _, err = s.updateService.Upsert(req.Application, provider, host, req.Version, req); err != nil {
		return err
	}

	return nil
}

func (s *WebhookInvocationService) ExecuteDiun(id string, token string, req api.WebhookDiunRequest) error {
	if id == "" || token == "" {
		return service_error.ErrValidationNotBlank
	}

	var e *model.Webhook
	var err error

	if e, err = s.webhookService.Get(id); err != nil {
		return service_error.ErrResourceNotFound
	}

	if e.Token != token {
		return service_error.ErrResourceAccessDenied
	}

	host := req.Hostname
	if e.IgnoreHost {
		host = hostIgnoreReplacement
	}

	// assume the "image" attribute has a : separator at the end
	ss := strings.Split(req.Image, ":")
	version := ss[len(ss)-1]
	app := strings.Join(ss, "")
	app = strings.ReplaceAll(app, version, "")

	if version == "" {
		version = req.Digest
	}

	var provider string
	if e.Label == "" {
		provider = providerDiun
	} else {
		provider = e.Label
	}

	if _, err = s.updateService.Upsert(app, provider, host, version, req); err != nil {
		return err
	}

	return nil
}
