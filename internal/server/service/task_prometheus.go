package service

import (
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

type PrometheusTask struct {
	updateService     *UpdateService
	eventService      *EventService
	actionService     *ActionService
	webhookService    *WebhookService
	prometheusService *PrometheusService
	taskService       *TaskService
	prometheusConfig  *config.Prometheus
}

const (
	jobNamePrometheusRefresh = "PROMETHEUS_REFRESH"
)

func NewPrometheusTask(u *UpdateService, e *EventService, w *WebhookService, a *ActionService, p *PrometheusService, t *TaskService, c *config.Prometheus) *PrometheusTask {
	return &PrometheusTask{
		updateService:     u,
		eventService:      e,
		actionService:     a,
		webhookService:    w,
		prometheusService: p,
		taskService:       t,
		prometheusConfig:  c,
	}
}

// Init initializes background tasks for the service, should be called directly after NewPrometheusTask
func (s *PrometheusTask) Init() error {
	return s.configurePrometheusRefreshTask()
}

func (s *PrometheusTask) configurePrometheusRefreshTask() error {
	if !s.prometheusConfig.Enabled {
		return nil
	}

	if err := s.prometheusService.RegisterGaugeNoLabels(constant.MetricUpdatesTotal, constant.MetricUpdatesTotalHelp); err != nil {
		return err
	}
	if err := s.prometheusService.RegisterGaugeNoLabels(constant.MetricUpdatesPending, constant.MetricUpdatesPendingHelp); err != nil {
		return err
	}
	if err := s.prometheusService.RegisterGaugeNoLabels(constant.MetricUpdatesIgnored, constant.MetricUpdatesIgnoredHelp); err != nil {
		return err
	}
	if err := s.prometheusService.RegisterGaugeNoLabels(constant.MetricUpdatesApproved, constant.MetricUpdatesApprovedHelp); err != nil {
		return err
	}
	if err := s.prometheusService.RegisterGaugeNoLabels(constant.MetricWebhooks, constant.MetricWebhooksHelp); err != nil {
		return err
	}
	if err := s.prometheusService.RegisterGaugeNoLabels(constant.MetricEvents, constant.MetricEventsHelp); err != nil {
		return err
	}
	if err := s.prometheusService.RegisterGaugeNoLabels(constant.MetricActions, constant.MetricActionsHelp); err != nil {
		return err
	}

	runnable := func() {
		updates, updatesError := s.updateService.GetAll()

		if updatesError = s.prometheusService.SetGaugeNoLabels(constant.MetricUpdatesTotal, float64(len(updates))); updatesError != nil {
			log.Error().Msgf("Could not refresh updates all prometheus metric. Reason: %s", updatesError.Error())
		}

		var pendingTotal int64
		var ignoredTotal int64
		var ackTotal int64

		for _, update := range updates {
			if constant.UpdateStatePending.String() == update.State {
				pendingTotal += 1
			} else if constant.UpdateStateIgnored.String() == update.State {
				ignoredTotal += 1
			} else if constant.UpdateStateApproved.String() == update.State {
				ackTotal += 1
			}
		}

		if updatesError = s.prometheusService.SetGaugeNoLabels(constant.MetricUpdatesPending, float64(pendingTotal)); updatesError != nil {
			log.Error().Msgf("Could not refresh updates pending prometheus metric. Reason: %s", updatesError.Error())
		}
		if updatesError = s.prometheusService.SetGaugeNoLabels(constant.MetricUpdatesIgnored, float64(ignoredTotal)); updatesError != nil {
			log.Error().Msgf("Could not refresh updates ignored prometheus metric. Reason: %s", updatesError.Error())
		}
		if updatesError = s.prometheusService.SetGaugeNoLabels(constant.MetricUpdatesApproved, float64(ackTotal)); updatesError != nil {
			log.Error().Msgf("Could not refresh updates approved prometheus metric. Reason: %s", updatesError.Error())
		}

		var webhooksTotal int64
		var webhooksError error
		webhooksTotal, webhooksError = s.webhookService.Count()
		if webhooksError = s.prometheusService.SetGaugeNoLabels(constant.MetricWebhooks, float64(webhooksTotal)); webhooksError != nil {
			log.Error().Msgf("Could not refresh webhooks prometheus metric. Reason: %s", webhooksError.Error())
		}

		var eventsTotal int64
		var eventsError error
		eventsTotal, eventsError = s.eventService.Count()
		if eventsError = s.prometheusService.SetGaugeNoLabels(constant.MetricEvents, float64(eventsTotal)); eventsError != nil {
			log.Error().Msgf("Could not refresh events prometheus metric. Reason: %s", eventsError.Error())
		}

		var actionsTotal int64
		var actionsError error
		actionsTotal, actionsError = s.actionService.Count()
		if actionsError = s.prometheusService.SetGaugeNoLabels(constant.MetricActions, float64(actionsTotal)); actionsError != nil {
			log.Error().Msgf("Could not refresh actions prometheus metric. Reason: %s", actionsError.Error())
		}
	}

	_, err := s.taskService.Enqueue(gocron.DurationJob(s.prometheusConfig.RefreshInterval), gocron.NewTask(runnable), jobNamePrometheusRefresh)
	return err
}
