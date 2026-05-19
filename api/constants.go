package api

const (
	HeaderAppName    = "X-App-Name"
	HeaderAppVersion = "X-App-Version"

	HeaderWebhookToken = "X-Webhook-Token" //nolint:gosec // header name, not a credential

	HeaderContentType                = "Content-Type"
	HeaderContentTypeApplicationJson = "application/json"
)
