package service

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"git.myservermanager.com/varakh/upda/internal/str"
	"github.com/rs/zerolog/log"
)

type WebhookService struct {
	repo          repository.WebhookRepository
	webhookConfig *config.Webhook
}

func NewWebhookService(r repository.WebhookRepository, c *config.Webhook) *WebhookService {
	return &WebhookService{
		repo:          r,
		webhookConfig: c,
	}
}

func (s *WebhookService) Get(id string) (*model.Webhook, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e, err := s.repo.Find(id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *WebhookService) Create(label string, t constant.WebhookType, ignoreHost bool, ignoreHostReplacement string) (*model.Webhook, error) {
	if label == "" || ignoreHostReplacement == "" || t == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var token string

	if token, err = str.GenerateSecureRandomString(s.webhookConfig.TokenLength); err != nil {
		return nil, service_error.NewServiceError(service_error.ErrCodeGeneral, fmt.Errorf("token generation failed: %w", err))
	}

	var e *model.Webhook
	if e, err = s.repo.Create(label, t.String(), token, ignoreHost, ignoreHostReplacement); err != nil {
		return nil, err
	} else {
		log.Info().Msg("Created webhook")
		return e, nil
	}
}

func (s *WebhookService) UpdateLabel(id string, label string) (*model.Webhook, error) {
	if id == "" || label == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Webhook
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateLabel(id, label); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified webhook '%v'", id)
	return e, nil
}

func (s *WebhookService) UpdateIgnoreHostReplacement(id string, ignoreHostReplacement string) (*model.Webhook, error) {
	if id == "" || ignoreHostReplacement == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Webhook
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateIgnoreHostReplacement(id, ignoreHostReplacement); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified webhook '%v'", id)
	return e, nil
}

func (s *WebhookService) UpdateIgnoreHost(id string, ignoreHost bool) (*model.Webhook, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Webhook
	var err error

	if e, err = s.Get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateIgnoreHost(id, ignoreHost); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified webhook '%v'", id)
	return e, nil
}

func (s *WebhookService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	e, err := s.Get(id)
	if err != nil {
		return err
	}

	if _, err = s.repo.Delete(e.ID.String()); err != nil {
		return err
	}

	log.Info().Msgf("Deleted webhook '%v'", id)

	return nil
}

func (s *WebhookService) Paginate(page int, pageSize int, orderBy string, order string) ([]*model.Webhook, error) {
	return s.repo.Paginate(page, pageSize, orderBy, order)
}

func (s *WebhookService) Count() (int64, error) {
	return s.repo.Count()
}
