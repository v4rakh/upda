package server

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/commons"
	"go.uber.org/zap"
)

type webhookService struct {
	repo          webhookRepository
	webhookConfig *webhookConfig
}

func newWebhookService(r webhookRepository, c *webhookConfig) *webhookService {
	return &webhookService{
		repo:          r,
		webhookConfig: c,
	}
}

func (s *webhookService) get(id string) (*Webhook, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	e, err := s.repo.find(id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *webhookService) create(label string, t api.WebhookType, ignoreHost bool) (*Webhook, error) {
	if label == "" || t == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var token string

	if token, err = commons.GenerateSecureRandomString(s.webhookConfig.tokenLength); err != nil {
		return nil, newServiceError(general, fmt.Errorf("token generation failed: %w", err))
	}

	var e *Webhook
	if e, err = s.repo.create(label, t, token, ignoreHost); err != nil {
		return nil, err
	} else {
		zap.L().Sugar().Info("Created webhook")
		return e, nil
	}
}

func (s *webhookService) updateLabel(id string, label string) (*Webhook, error) {
	if id == "" || label == "" {
		return nil, errorValidationNotBlank
	}

	var e *Webhook
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateLabel(id, label); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified webhook '%v'", id)
	return e, nil
}

func (s *webhookService) updateIgnoreHost(id string, ignoreHost bool) (*Webhook, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e *Webhook
	var err error

	if e, err = s.get(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.updateIgnoreHost(id, ignoreHost); err != nil {
		return nil, err
	}

	zap.L().Sugar().Infof("Modified webhook '%v'", id)
	return e, nil
}

func (s *webhookService) delete(id string) error {
	if id == "" {
		return errorValidationNotBlank
	}

	e, err := s.get(id)
	if err != nil {
		return err
	}

	if _, err = s.repo.delete(e.ID.String()); err != nil {
		return err
	}

	zap.L().Sugar().Infof("Deleted webhook '%v'", id)

	return nil
}

func (s *webhookService) paginate(page int, pageSize int, orderBy string, order string) ([]*Webhook, error) {
	return s.repo.paginate(page, pageSize, orderBy, order)
}

func (s *webhookService) count() (int64, error) {
	return s.repo.count()
}
