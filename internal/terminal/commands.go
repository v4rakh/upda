package terminal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli/v3"
)

const (
	envServerUrl    = "UPDA_SERVER_URL"
	envWebhookId    = "UPDA_WEBHOOK_ID"
	envWebhookToken = "UPDA_WEBHOOK_TOKEN" //nolint:gosec // env var name, not a credential

	flagUrl                = "url"
	flagUser               = "user"
	flagPass               = "pass"
	flagWebhookId          = "webhook-id"
	flagWebhookToken       = "webhook-token"
	flagWebhookApplication = "application"
	flagWebhookHost        = "host"
	flagWebhookProvider    = "provider"
	flagWebhookVersion     = "application-version"
	flagWebhookMetadata    = "metadata"

	flagTimeout = "timeout"

	webhooksUrlPath = "/api/v1/webhooks"
)

var (
	timeout            time.Duration
	serverUrl          string
	webhookId          string
	webhookToken       string
	webhookApplication string
	webhookHost        string
	webhookProvider    string
	webhookVersion     string
	webhookMetadata    map[string]string

	timeoutFlag = &cli.DurationFlag{
		Name:        flagTimeout,
		Usage:       "optional flag to determine maximum timeout to query upda",
		Aliases:     []string{"to"},
		Required:    false,
		Value:       10 * time.Second,
		Destination: &timeout,
	}
	urlFlag = &cli.StringFlag{
		Name:        flagUrl,
		Usage:       "the server url (FQDN without context path)",
		Required:    true,
		Aliases:     []string{"s"},
		Sources:     cli.EnvVars(envServerUrl),
		Destination: &serverUrl,
	}
	webhookIdFlag = &cli.StringFlag{
		Name:        flagWebhookId,
		Usage:       "webhook id",
		Required:    true,
		Aliases:     []string{"i"},
		Sources:     cli.EnvVars(envWebhookId),
		Destination: &webhookId,
	}
	webhookTokenFlag = &cli.StringFlag{
		Name:        flagWebhookToken,
		Usage:       "webhook token",
		Required:    true,
		Aliases:     []string{"t"},
		Sources:     cli.EnvVars(envWebhookToken),
		Destination: &webhookToken,
	}
	webhookApplicationFlag = &cli.StringFlag{
		Name:        flagWebhookApplication,
		Usage:       "Application name",
		Required:    false,
		Destination: &webhookApplication,
	}
	webhookHostFlag = &cli.StringFlag{
		Name:        flagWebhookHost,
		Usage:       "Host",
		Required:    false,
		Destination: &webhookHost,
	}
	webhookProviderFlag = &cli.StringFlag{
		Name:        flagWebhookProvider,
		Usage:       "Provider",
		Required:    false,
		Destination: &webhookProvider,
	}
	webhookVersionFlag = &cli.StringFlag{
		Name:        flagWebhookVersion,
		Usage:       "Version",
		Required:    false,
		Destination: &webhookVersion,
	}
	webhookMetadataFlag = &cli.StringMapFlag{
		Name:        flagWebhookMetadata,
		Usage:       "Metadata to add, you can provide multiple --metadata",
		Destination: &webhookMetadata,
	}

	WebhookSendCmd = &cli.Command{
		Name:  "send",
		Usage: "Sends data to a webhook, prefers the flag based approach, if not all required properties given, tries to parse 1st argument as JSON",
		Flags: []cli.Flag{
			urlFlag,
			webhookIdFlag,
			webhookTokenFlag,
			webhookApplicationFlag,
			webhookHostFlag,
			webhookProviderFlag,
			webhookVersionFlag,
			webhookMetadataFlag,
			timeoutFlag,
		},
		ArgsUsage: "[<full json payload if no --application --host --application-version has been provided>]",
		Action:    webhookSend,
	}
)

func webhookSend(_ context.Context, cmd *cli.Command) error {
	if err := failIfFlagsNotPresent(cmd, []string{flagUrl, flagWebhookId, flagWebhookToken}); err != nil {
		return cli.Exit(err, 1)
	}

	application := cmd.String(flagWebhookApplication)
	host := cmd.String(flagWebhookHost)
	applicationVersion := cmd.String(flagWebhookVersion)

	var payload interface{}

	if application != "" && host != "" && applicationVersion != "" {
		structuredPayload := api.WebhookGenericRequest{}
		structuredPayload.Application = application
		structuredPayload.Host = host
		structuredPayload.Version = applicationVersion

		provider := cmd.String(flagWebhookProvider)
		metaData := cmd.StringMap(flagWebhookMetadata)

		if provider != "" {
			structuredPayload.Provider = provider
		}

		if len(metaData) > 0 {
			structuredPayload.Metadata = metaData
		}

		payload = structuredPayload
	} else {
		// no flags given, fallback to plain JSON argument
		if !cmd.Args().Present() || cmd.Args().Len() < 1 {
			return cli.Exit(errors.New("args required - try 'webhook send help'"), 1)
		}

		// validate payload is valid json
		unstructuredPayload := cmd.Args().First()
		if unstructuredPayload == "" {
			return cli.Exit(errors.New("payload cannot be blank"), 1)
		}
		if !json.Valid([]byte(unstructuredPayload)) {
			return cli.Exit(errors.New("payload is not valid JSON"), 1)
		}

		payload = unstructuredPayload
	}

	var errorRes api.ErrorResponse
	url := fmt.Sprintf("%s/%s", webhooksUrlPath, cmd.String(flagWebhookId))
	client := newClient(cmd)
	res, err := client.R().
		SetHeader(api.HeaderContentType, api.HeaderContentTypeApplicationJson).
		SetHeader(api.HeaderWebhookToken, cmd.String(flagWebhookToken)).
		SetBody(payload).
		SetError(&errorRes).
		Post(url)

	if err != nil {
		return cli.Exit(fmt.Errorf("error during webhook invocation: %w", err), 1)
	}
	if !res.IsSuccess() {
		return cli.Exit(fmt.Sprintf("error during webhook invocation: (%d) %+v", res.StatusCode(), errorRes), 1)
	}

	return nil
}

func failIfFlagsNotPresent(cmd *cli.Command, flagKeys []string) error {
	if flagKeys == nil {
		return errors.New("flagKeys cannot be null")
	}

	for _, key := range flagKeys {
		if cmd.String(key) == "" {
			return fmt.Errorf("'%v' is required but blank", key)
		}
	}

	return nil
}

func newClient(cmd *cli.Command) *resty.Client {
	client := resty.New()
	client.SetHeader("User-Agent", fmt.Sprintf("%s/%s", meta.Name, meta.Version))
	client.SetDisableWarn(true)
	client.SetTimeout(cmd.Duration(flagTimeout))
	client.SetBaseURL(cmd.String(flagUrl))

	username := cmd.String(flagUser)
	pass := cmd.String(flagPass)

	if username != "" && pass != "" {
		client.SetBasicAuth(cmd.String(flagUser), cmd.String(flagPass))
	}

	return client
}
