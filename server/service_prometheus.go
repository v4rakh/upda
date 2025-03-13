package server

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/commons"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

type prometheusService struct {
	router           *gin.Engine
	prometheus       *ginprom.Prometheus
	prometheusConfig *prometheusConfig
	serverConfig     *serverConfig
}

func newPrometheusService(r *gin.Engine, pc *prometheusConfig, sc *serverConfig) *prometheusService {
	var p *ginprom.Prometheus

	if pc.enabled {
		path := fmt.Sprintf("%s%s", sc.basePath, pc.path)
		if pc.secureTokenEnabled {
			p = ginprom.New(
				ginprom.Engine(r),
				ginprom.Namespace(commons.Name),
				ginprom.Subsystem(""),
				ginprom.Path(path),
				ginprom.Ignore(path),
				ginprom.Token(pc.secureToken),
			)
		} else {
			p = ginprom.New(
				ginprom.Engine(r),
				ginprom.Namespace(commons.Name),
				ginprom.Subsystem(""),
				ginprom.Ignore(path),
				ginprom.Path(path),
			)
		}
	}

	return &prometheusService{
		prometheus:       p,
		prometheusConfig: pc,
	}
}

func (s *prometheusService) init() error {
	if !s.prometheusConfig.enabled {
		return nil
	}

	if err := s.registerGaugeNoLabels(metricUpdatesTotal, metricUpdatesTotalHelp); err != nil {
		return err
	}
	if err := s.registerGaugeNoLabels(metricUpdatesPending, metricUpdatesPendingHelp); err != nil {
		return err
	}
	if err := s.registerGaugeNoLabels(metricUpdatesIgnored, metricUpdatesIgnoredHelp); err != nil {
		return err
	}
	if err := s.registerGaugeNoLabels(metricUpdatesApproved, metricUpdatesApprovedHelp); err != nil {
		return err
	}
	if err := s.registerGaugeNoLabels(metricWebhooks, metricWebhooksHelp); err != nil {
		return err
	}
	if err := s.registerGaugeNoLabels(metricEvents, metricEventsHelp); err != nil {
		return err
	}
	if err := s.registerGaugeNoLabels(metricActions, metricActionsHelp); err != nil {
		return err
	}

	return nil
}

func (s *prometheusService) registerGaugeNoLabels(name string, help string) error {
	return s.registerGauge(name, help, make([]string, 0))
}

func (s *prometheusService) registerGauge(name string, help string, labels []string) error {
	if !s.prometheusConfig.enabled {
		return nil
	}

	if name == "" || help == "" {
		return errorValidationNotBlank
	}

	s.prometheus.AddCustomGauge(name, help, labels)

	return nil
}

func (s *prometheusService) registerCounterNoLabels(name string, help string) error {
	return s.registerCounter(name, help, make([]string, 0))
}

func (s *prometheusService) registerCounter(name string, help string, labels []string) error {
	if !s.prometheusConfig.enabled {
		return nil
	}

	if name == "" || help == "" {
		return errorValidationNotBlank
	}

	s.prometheus.AddCustomCounter(name, help, labels)

	return nil
}

func (s *prometheusService) setGaugeNoLabels(name string, value float64) error {
	return s.setGauge(name, make([]string, 0), value)
}

func (s *prometheusService) setGauge(name string, labelValues []string, value float64) error {
	if !s.prometheusConfig.enabled {
		return nil
	}

	if name == "" {
		return errorValidationNotBlank
	}

	if err := s.prometheus.SetGaugeValue(name, labelValues, value); err != nil {
		return err
	}

	return nil
}

func (s *prometheusService) increaseCounterNoLabels(name string) error {
	return s.increaseCounter(name, make([]string, 0))
}

func (s *prometheusService) increaseCounter(name string, labelValues []string) error {
	if !s.prometheusConfig.enabled {
		return nil
	}

	if name == "" {
		return errorValidationNotBlank
	}

	if err := s.prometheus.IncrementCounterValue(name, labelValues); err != nil {
		return err
	}

	return nil
}
