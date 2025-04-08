package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"strings"
)

const (
	providerDiun          = "oci"
	hostIgnoreReplacement = "global"
)

type webhookInvocationService struct {
	updateService  *updateService
	webhookService *webhookService
	webhookConfig  *webhookConfig
}

func newWebhookInvocationService(w *webhookService, u *updateService, c *webhookConfig) *webhookInvocationService {
	return &webhookInvocationService{
		updateService:  u,
		webhookService: w,
		webhookConfig:  c,
	}
}

func (s *webhookInvocationService) executeGeneric(id string, token string, req api.WebhookGenericRequest) error {
	if id == "" || token == "" {
		return errorValidationNotBlank
	}

	var e *Webhook
	var err error

	if e, err = s.webhookService.get(id); err != nil {
		return errorResourceNotFound
	}

	if e.Token != token {
		return errorResourceAccessDenied
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

	if _, err = s.updateService.upsert(req.Application, provider, host, req.Version, req); err != nil {
		return err
	}

	return nil
}

func (s *webhookInvocationService) executeDiun(id string, token string, req api.WebhookDiunRequest) error {
	if id == "" || token == "" {
		return errorValidationNotBlank
	}

	var e *Webhook
	var err error

	if e, err = s.webhookService.get(id); err != nil {
		return errorResourceNotFound
	}

	if e.Token != token {
		return errorResourceAccessDenied
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

	if _, err = s.updateService.upsert(app, provider, host, version, req); err != nil {
		return err
	}

	return nil
}
