package service

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

const (
	metricUpdatesTotal     = "upda_updates_all"
	metricUpdatesTotalHelp = "amount of all updates"

	metricUpdatesPending     = "upda_updates_pending"
	metricUpdatesPendingHelp = "amount of all updates in pending state"

	metricUpdatesIgnored     = "upda_updates_ignored"
	metricUpdatesIgnoredHelp = "amount of all updates in ignored state"

	metricUpdatesApproved     = "upda_updates_approved"
	metricUpdatesApprovedHelp = "amount of all updates in approved state"

	metricWebhooks     = "upda_webhooks"
	metricWebhooksHelp = "amount of all webhooks"

	metricEvents     = "upda_events"
	metricEventsHelp = "amount of all events"

	metricActions     = "upda_actions"
	metricActionsHelp = "amount of all actions"
)

type PrometheusService struct {
	router           *gin.Engine
	prometheus       *ginprom.Prometheus
	prometheusConfig *config.Prometheus
}

func NewPrometheusService(e *gin.Engine, c *config.Prometheus) *PrometheusService {
	var p *ginprom.Prometheus

	if !c.Enabled {
		return &PrometheusService{
			prometheus:       p,
			prometheusConfig: c,
		}
	}

	path := fmt.Sprintf("%s%s", c.BasePath, c.Path)
	if c.SecureTokenEnabled {
		p = ginprom.New(
			ginprom.Engine(e),
			ginprom.Namespace(""),
			ginprom.Subsystem(""),
			ginprom.Path(path),
			ginprom.Ignore(path),
			ginprom.Token(c.SecureToken),
		)
	} else {
		p = ginprom.New(
			ginprom.Engine(e),
			ginprom.Namespace(""),
			ginprom.Subsystem(""),
			ginprom.Ignore(path),
			ginprom.Path(path),
		)
	}

	return &PrometheusService{
		prometheus:       p,
		prometheusConfig: c,
	}
}

func (s *PrometheusService) GetProm() *ginprom.Prometheus {
	return s.prometheus
}

func (s *PrometheusService) Init() error {
	if !s.prometheusConfig.Enabled {
		return nil
	}

	if err := s.RegisterGaugeNoLabels(metricUpdatesTotal, metricUpdatesTotalHelp); err != nil {
		return err
	}
	if err := s.RegisterGaugeNoLabels(metricUpdatesPending, metricUpdatesPendingHelp); err != nil {
		return err
	}
	if err := s.RegisterGaugeNoLabels(metricUpdatesIgnored, metricUpdatesIgnoredHelp); err != nil {
		return err
	}
	if err := s.RegisterGaugeNoLabels(metricUpdatesApproved, metricUpdatesApprovedHelp); err != nil {
		return err
	}
	if err := s.RegisterGaugeNoLabels(metricWebhooks, metricWebhooksHelp); err != nil {
		return err
	}
	if err := s.RegisterGaugeNoLabels(metricEvents, metricEventsHelp); err != nil {
		return err
	}
	if err := s.RegisterGaugeNoLabels(metricActions, metricActionsHelp); err != nil {
		return err
	}

	return nil
}

func (s *PrometheusService) RegisterGaugeNoLabels(name string, help string) error {
	return s.RegisterGauge(name, help, make([]string, 0))
}

func (s *PrometheusService) RegisterGauge(name string, help string, labels []string) error {
	if !s.prometheusConfig.Enabled {
		return nil
	}

	if name == "" || help == "" {
		return service_error.ErrValidationNotBlank
	}

	s.prometheus.AddCustomGauge(name, help, labels)

	return nil
}

func (s *PrometheusService) RegisterCounterNoLabels(name string, help string) error {
	return s.RegisterCounter(name, help, make([]string, 0))
}

func (s *PrometheusService) RegisterCounter(name string, help string, labels []string) error {
	if !s.prometheusConfig.Enabled {
		return nil
	}

	if name == "" || help == "" {
		return service_error.ErrValidationNotBlank
	}

	s.prometheus.AddCustomCounter(name, help, labels)

	return nil
}

func (s *PrometheusService) SetGaugeNoLabels(name string, value float64) error {
	return s.SetGauge(name, make([]string, 0), value)
}

func (s *PrometheusService) SetGauge(name string, labelValues []string, value float64) error {
	if !s.prometheusConfig.Enabled {
		return nil
	}

	if name == "" {
		return service_error.ErrValidationNotBlank
	}

	if err := s.prometheus.SetGaugeValue(name, labelValues, value); err != nil {
		return err
	}

	return nil
}

func (s *PrometheusService) IncreaseCounterNoLabels(name string) error {
	return s.IncreaseCounter(name, make([]string, 0))
}

func (s *PrometheusService) IncreaseCounter(name string, labelValues []string) error {
	if !s.prometheusConfig.Enabled {
		return nil
	}

	if name == "" {
		return service_error.ErrValidationNotBlank
	}

	if err := s.prometheus.IncrementCounterValue(name, labelValues); err != nil {
		return err
	}

	return nil
}
