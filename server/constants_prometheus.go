package server

const (
	metricUpdatesTotal     = "updates_all"
	metricUpdatesTotalHelp = "amount of all updates"

	metricUpdatesPending     = "updates_pending"
	metricUpdatesPendingHelp = "amount of all updates in pending state"

	metricUpdatesIgnored     = "updates_ignored"
	metricUpdatesIgnoredHelp = "amount of all updates in ignored state"

	metricUpdatesApproved     = "updates_approved"
	metricUpdatesApprovedHelp = "amount of all updates in approved state"

	metricUpdates     = "updates"
	metricUpdatesHelp = "details for all updates, 0=pending, 1=approved, 2=ignored"

	metricWebhooks     = "webhooks"
	metricWebhooksHelp = "amount of all webhooks"

	metricEvents     = "events"
	metricEventsHelp = "amount of all events"
)
