package terminal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli/v3"
	"os"
	"slices"
	"strconv"
	"text/tabwriter"
	"time"
)

const (
	envServerUrl    = "UPDA_SERVER_URL"
	envUser         = "UPDA_USER"
	envPassword     = "UPDA_PASSWORD"
	envWebhookId    = "UPDA_WEBHOOK_ID"
	envWebhookToken = "UPDA_WEBHOOK_TOKEN"

	flagUrl                = "url"
	flagUser               = "user"
	flagPass               = "pass"
	flagUpdatePageSize     = "page-size"
	flagWebhookId          = "webhook-id"
	flagWebhookToken       = "webhook-token"
	flagWebhookApplication = "application"
	flagWebhookHost        = "host"
	flagWebhookProvider    = "provider"
	flagWebhookVersion     = "application-version"
	flagWebhookMetadata    = "metadata"

	flagRaw     = "raw"
	flagTimeout = "timeout"

	webhooksUrlPath = "/api/v1/webhooks"
	updatesUrlPath  = "/api/v1/updates"
)

var (
	raw                bool
	timeout            time.Duration
	serverUrl          string
	user               string
	password           string
	updatePageSize     int
	webhookId          string
	webhookToken       string
	webhookApplication string
	webhookHost        string
	webhookProvider    string
	webhookVersion     string
	webhookMetadata    map[string]string

	rawFlag = &cli.BoolFlag{
		Name:        flagRaw,
		Usage:       "on success raw JSON data from response is returned",
		Aliases:     []string{"r"},
		Value:       false,
		Destination: &raw,
	}
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
	userFlag = &cli.StringFlag{
		Name:        flagUser,
		Usage:       "user",
		Required:    true,
		Aliases:     []string{"u"},
		Sources:     cli.EnvVars(envUser),
		Destination: &user,
	}
	passwordFlag = &cli.StringFlag{
		Name:        flagPass,
		Usage:       "password",
		Required:    true,
		Aliases:     []string{"p"},
		Sources:     cli.EnvVars(envPassword),
		Destination: &password,
	}
	updatePageSizeFlag = &cli.IntFlag{
		Name:        flagUpdatePageSize,
		Usage:       "update show page size",
		Value:       10000,
		Required:    false,
		Aliases:     []string{"ps"},
		Destination: &updatePageSize,
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

	WebhookCreateCmd = &cli.Command{
		Name:  "create",
		Usage: "Creates a webhook",
		Flags: []cli.Flag{
			urlFlag,
			userFlag,
			passwordFlag,
			rawFlag,
			timeoutFlag,
		},
		ArgsUsage: "<label> [<type (generic|diun, default: generic)>] [<ignore-host (true|false, default: false)>]",
		Action:    webhookCreate,
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

	UpdateShowCmd = &cli.Command{
		Name:  "show",
		Usage: "Shows updates",
		Flags: []cli.Flag{
			urlFlag,
			userFlag,
			passwordFlag,
			updatePageSizeFlag,
			rawFlag,
			timeoutFlag,
		},
		Action: updateShow,
	}
)

func webhookCreate(_ context.Context, cmd *cli.Command) error {
	if err := failIfFlagsNotPresent(cmd, []string{flagUrl, flagUser, flagPass}); err != nil {
		return cli.Exit(err, 1)
	}
	if !cmd.Args().Present() {
		return cli.Exit(errors.New("args required - try 'webhook create help'"), 1)
	}
	// validate label
	label := cmd.Args().First()
	if label == "" || len(label) > 255 {
		return cli.Exit(errors.New("label cannot be blank or only be 255 characters long"), 1)
	}

	// validate type
	t := cmd.Args().Get(1)
	validTypes := []string{constant.WebhookTypeGeneric.String(), constant.WebhookTypeDiun.String()}
	if t == "" {
		t = constant.WebhookTypeGeneric.String()
	}

	if !slices.Contains(validTypes, t) {
		return cli.Exit(errors.New(fmt.Sprintf("type must be one of %v", validTypes)), 1)
	}

	ignoreHost := cmd.Args().Get(2) == "true"

	// fully constructed payload
	payload := api.CreateWebhookRequest{
		Label:      label,
		Type:       t,
		IgnoreHost: ignoreHost,
	}

	var successRes api.WebhookSingleResponse
	var errorRes api.ErrorResponse
	url := webhooksUrlPath
	client := newClient(cmd)
	res, err := client.R().
		SetHeader(api.HeaderContentType, api.HeaderContentTypeApplicationJson).
		SetBody(&payload).
		SetResult(&successRes).
		SetError(&errorRes).
		Post(url)

	if err != nil {
		return cli.Exit(fmt.Errorf("error during webhook creation: %w", err), 1)
	}
	if !res.IsSuccess() {
		return cli.Exit(fmt.Sprintf("error during webhook creation: (%d) %+v", res.StatusCode(), errorRes), 1)
	}

	if cmd.Bool(flagRaw) {
		fmt.Println(string(res.Body()))
		return nil
	}

	fmt.Printf("ID\t%v\n", successRes.Data.ID)
	fmt.Printf("URL\t%v\n", fmt.Sprintf("%s/%s", url, successRes.Data.ID))
	fmt.Printf("Token\t%v\n", successRes.Data.Token)
	fmt.Printf("Type\t%v\n", successRes.Data.Type)

	return nil
}

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

		if metaData != nil && len(metaData) > 0 {
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

func updateShow(_ context.Context, cmd *cli.Command) error {
	if err := failIfFlagsNotPresent(cmd, []string{flagUrl, flagUser, flagPass}); err != nil {
		return cli.Exit(err, 1)
	}

	var successRes api.UpdateDataPageResponse
	var errorRes api.ErrorResponse
	url := updatesUrlPath
	client := newClient(cmd)
	res, err := client.R().
		SetHeader(api.HeaderContentType, api.HeaderContentTypeApplicationJson).
		SetQueryParam("pageSize", strconv.Itoa(cmd.Int(flagUpdatePageSize))).
		SetResult(&successRes).
		SetError(&errorRes).
		Get(url)

	if err != nil {
		return cli.Exit(fmt.Errorf("error during showing updates: %w", err), 1)
	}
	if !res.IsSuccess() {
		return cli.Exit(fmt.Sprintf("error during showing updates: (%d) %+v", res.StatusCode(), errorRes), 1)
	}

	if cmd.Bool(flagRaw) {
		fmt.Println(string(res.Body()))
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.Debug)

	for _, u := range successRes.Data.Content {
		if _, err = fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", u.Application, u.Host, u.Provider, u.Version, u.State); err != nil {
			return cli.Exit(fmt.Sprintf("error during showing updates: %+v", errorRes), 1)
		}
	}
	if err = w.Flush(); err != nil {
		return cli.Exit(fmt.Sprintf("error during showing updates: %+v", errorRes), 1)
	}

	return nil
}

func failIfFlagsNotPresent(cmd *cli.Command, flagKeys []string) error {
	if flagKeys == nil {
		return errors.New("flagKeys cannot be null")
	}

	for _, key := range flagKeys {
		if cmd.String(key) == "" {
			return errors.New(fmt.Sprintf("'%v' is required but blank", key))
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
