package terminal

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/commons"
	"git.myservermanager.com/varakh/upda/util"
	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"text/tabwriter"
)

const (
	name    = "upda-cli"
	desc    = "a commandline helper for upda"
	version = commons.Version

	envServerUrl    = "UPDA_SERVER_URL"
	envUser         = "UPDA_USER"
	envPassword     = "UPDA_PASSWORD"
	envWebhookId    = "UPDA_WEBHOOK_ID"
	envWebhookToken = "UPDA_WEBHOOK_TOKEN"

	flagServerUrl      = "server-url"
	flagUser           = "user"
	flagPass           = "pass"
	flagWebhookId      = "webhook-id"
	flagWebhookToken   = "webhook-token"
	flagUpdatePageSize = "page-size"

	flagRaw = "raw"

	webhooksUrlPath = "/api/v1/webhooks"
	updatesUrlPath  = "/api/v1/updates"
)

func Start() {
	var raw bool
	var serverUrl string
	var user string
	var password string
	var webhookId string
	var webhookToken string
	var updatePageSize int64

	rawFlag := &cli.BoolFlag{
		Name:        flagRaw,
		Usage:       "on success raw JSON data from response is returned",
		Aliases:     []string{"r"},
		Value:       false,
		Destination: &raw,
	}

	serverUrlFlag := &cli.StringFlag{
		Name:        flagServerUrl,
		Usage:       "the server url (FQDN without context path)",
		Required:    true,
		Aliases:     []string{"s"},
		EnvVars:     []string{envServerUrl},
		Destination: &serverUrl,
	}
	userFlag := &cli.StringFlag{
		Name:        flagUser,
		Usage:       "user",
		Required:    true,
		Aliases:     []string{"u"},
		EnvVars:     []string{envUser},
		Destination: &user,
	}
	passwordFlag := &cli.StringFlag{
		Name:        flagPass,
		Usage:       "password",
		Required:    true,
		Aliases:     []string{"p"},
		EnvVars:     []string{envPassword},
		Destination: &password,
	}
	webhookIdFlag := &cli.StringFlag{
		Name:        flagWebhookId,
		Usage:       "webhook id",
		Required:    true,
		Aliases:     []string{"i"},
		EnvVars:     []string{envWebhookId},
		Destination: &webhookId,
	}
	webhookTokenFlag := &cli.StringFlag{
		Name:        flagWebhookToken,
		Usage:       "webhook token",
		Required:    true,
		Aliases:     []string{"t"},
		EnvVars:     []string{envWebhookToken},
		Destination: &webhookToken,
	}
	updatePageSizeFlag := &cli.Int64Flag{
		Name:        flagUpdatePageSize,
		Usage:       "update show page size",
		Value:       10000,
		Required:    false,
		Aliases:     []string{"ps"},
		Destination: &updatePageSize,
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "show version",
	}
	app := &cli.App{
		Name:                 name,
		Usage:                desc,
		Version:              version,
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "webhook",
				Aliases: []string{"w"},
				Usage:   "Options for webhook",
				Subcommands: []*cli.Command{
					{
						Name:    "create",
						Usage:   "Creates a webhook",
						Aliases: []string{"c"},
						Flags: []cli.Flag{
							serverUrlFlag,
							userFlag,
							passwordFlag,
							rawFlag,
						},
						ArgsUsage: "<label> [<type (generic|diun, default: generic)>] [<ignore-host (true|false, default: false)>]",
						Action:    webhookCreate,
					},
					{
						Name:    "send",
						Usage:   "Sends data to a webhook",
						Aliases: []string{"s"},
						Flags: []cli.Flag{
							serverUrlFlag,
							webhookIdFlag,
							webhookTokenFlag,
						},
						ArgsUsage: "<json payload>",
						Action:    webhookSend,
					},
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Options for update",
				Subcommands: []*cli.Command{
					{
						Name:    "show",
						Usage:   "Shows updates",
						Aliases: []string{"s"},
						Flags: []cli.Flag{
							serverUrlFlag,
							userFlag,
							passwordFlag,
							updatePageSizeFlag,
							rawFlag,
						},
						Action: updateShow,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func webhookCreate(cCtx *cli.Context) error {
	if err := failIfFlagsNotPresent(cCtx, []string{flagServerUrl, flagUser, flagPass}); err != nil {
		return cli.Exit(err, 1)
	}
	if !cCtx.Args().Present() {
		return cli.Exit(errors.New("args required - try 'webhook create help'"), 1)
	}
	// validate label
	label := cCtx.Args().First()
	if label == "" || len(label) > 255 {
		return cli.Exit(errors.New("label cannot be blank or only be 255 characters long"), 1)
	}

	// validate type
	t := cCtx.Args().Get(1)
	validTypes := []string{api.WebhookTypeGeneric.Value(), api.WebhookTypeDiun.Value()}
	if t == "" {
		t = api.WebhookTypeGeneric.Value()
	}

	if !util.FindInSlice(validTypes, t) {
		return cli.Exit(errors.New(fmt.Sprintf("type must be one of %v", validTypes)), 1)
	}

	ignoreHost := cCtx.Args().Get(2) == "true"

	// fully constructed payload
	payload := api.CreateWebhookRequest{
		Label:      label,
		Type:       t,
		IgnoreHost: ignoreHost,
	}

	var successRes api.WebhookSingleResponse
	var errorRes api.ErrorResponse
	url := cCtx.String(flagServerUrl) + webhooksUrlPath
	client := resty.New()
	client.SetDisableWarn(true)
	res, err := client.R().
		SetBasicAuth(cCtx.String(flagUser), cCtx.String(flagPass)).
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

	if cCtx.Bool(flagRaw) {
		fmt.Println(string(res.Body()))
		return nil
	}

	fmt.Printf("ID\t%v\n", successRes.Data.ID)
	fmt.Printf("URL\t%v\n", fmt.Sprintf("%s/%s", url, successRes.Data.ID))
	fmt.Printf("Token\t%v\n", successRes.Data.Token)
	fmt.Printf("Type\t%v\n", successRes.Data.Type)

	return nil
}

func webhookSend(cCtx *cli.Context) error {
	if err := failIfFlagsNotPresent(cCtx, []string{flagServerUrl, flagWebhookId, flagWebhookToken}); err != nil {
		return cli.Exit(err, 1)
	}

	if !cCtx.Args().Present() || cCtx.Args().Len() < 1 {
		return cli.Exit(errors.New("args required - try 'webhook send help'"), 1)
	}
	// validate payload is valid json
	payloadArg := cCtx.Args().First()
	if payloadArg == "" {
		return cli.Exit(errors.New("payload cannot be blank"), 1)
	}
	if !json.Valid([]byte(payloadArg)) {
		return cli.Exit(errors.New("payload is not valid JSON"), 1)
	}

	var errorRes api.ErrorResponse
	url := cCtx.String(flagServerUrl) + webhooksUrlPath + "/" + cCtx.String(flagWebhookId)
	client := resty.New()
	client.SetDisableWarn(true)
	res, err := client.R().
		SetHeader(api.HeaderContentType, api.HeaderContentTypeApplicationJson).
		SetHeader(api.HeaderWebhookToken, cCtx.String(flagWebhookToken)).
		SetBody(payloadArg).
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

func updateShow(cCtx *cli.Context) error {
	if err := failIfFlagsNotPresent(cCtx, []string{flagServerUrl, flagUser, flagPass}); err != nil {
		return cli.Exit(err, 1)
	}

	var successRes api.UpdateDataPageResponse
	var errorRes api.ErrorResponse
	url := cCtx.String(flagServerUrl) + updatesUrlPath + fmt.Sprintf("?pageSize=%d", cCtx.Int64(flagUpdatePageSize))
	client := resty.New()
	client.SetDisableWarn(true)
	res, err := client.R().
		SetBasicAuth(cCtx.String(flagUser), cCtx.String(flagPass)).
		SetHeader(api.HeaderContentType, api.HeaderContentTypeApplicationJson).
		SetResult(&successRes).
		SetError(&errorRes).
		Get(url)

	if err != nil {
		return cli.Exit(fmt.Errorf("error during showing updates: %w", err), 1)
	}
	if !res.IsSuccess() {
		return cli.Exit(fmt.Sprintf("error during showing updates: (%d) %+v", res.StatusCode(), errorRes), 1)
	}

	if cCtx.Bool(flagRaw) {
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

func failIfFlagsNotPresent(cCtx *cli.Context, flagKeys []string) error {
	if flagKeys == nil {
		return errors.New("flagKeys cannot be null")
	}

	for _, key := range flagKeys {
		if cCtx.String(key) == "" {
			return errors.New(fmt.Sprintf("'%v' is required but blank", key))
		}
	}

	return nil
}
