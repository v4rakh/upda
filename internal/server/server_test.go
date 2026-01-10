package server

import (
	"context"
	"encoding/json"
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/str"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	natPortFormat            = "%s/tcp"
	containerStartupTimeout  = 20 * time.Second
	containerDeadlineTimeout = 120 * time.Second
	containerLogTimeout      = 10 * time.Second
	imageEnvKey              = "IMAGE"

	postgresNetworkAlias = "db"
	postgresPort         = "5432"
	postgresImage        = "postgres:17-alpine"
	postgresUser         = "upda"
	postgresPass         = "upda"
	postgresDb           = "upda"
	appPort              = "8080"
)

type VerboseLogConsumer struct {
	name string
}

func NewVerboseLogConsumer(name string) *VerboseLogConsumer {
	return &VerboseLogConsumer{name: name}
}

func (g *VerboseLogConsumer) Accept(l testcontainers.Log) {
	fmt.Printf("Container[%s]:\t%s", g.name, string(l.Content))
}

func NewVerboseLogConsumerConfig(name string) *testcontainers.LogConsumerConfig {
	return &testcontainers.LogConsumerConfig{
		Consumers: []testcontainers.LogConsumer{NewVerboseLogConsumer(name)},
		Opts: []testcontainers.LogProductionOption{
			testcontainers.WithLogProductionTimeout(containerLogTimeout),
		},
	}
}

func createDatabaseContainer(ctx context.Context, t *testing.T, networkName string) testcontainers.Container {
	var err error
	req := testcontainers.ContainerRequest{
		Image:        postgresImage,
		ExposedPorts: []string{fmt.Sprintf(natPortFormat, postgresPort)},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {postgresNetworkAlias},
		},
		LogConsumerCfg: NewVerboseLogConsumerConfig(postgresNetworkAlias),
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPass,
			"POSTGRES_DB":       postgresDb,
		},
		WaitingFor: wait.ForAll(
			wait.ForExposedPort(),
			wait.ForLog("database system is ready to accept connections"),
		).WithStartupTimeoutDefault(containerStartupTimeout).
			WithDeadline(containerDeadlineTimeout),
	}

	var c testcontainers.Container
	if c, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}); err != nil {
		t.Fatalf("failed to start database container: %v", err)
	}
	testcontainers.CleanupContainer(t, c)

	return c
}

func setupTestEnvironment(t *testing.T) (testcontainers.Container, testcontainers.Container, string) {
	image := os.Getenv(imageEnvKey)
	if image == "" {
		t.Fatalf("%s environment variable is not set for tests", imageEnvKey)
	}

	var err error
	ctx := context.Background()

	// Create a network
	var n *testcontainers.DockerNetwork
	if n, err = network.New(ctx, network.WithAttachable(), network.WithDriver("bridge")); err != nil {
		t.Fatalf("failed to create docker network: %v", err)
	}
	testcontainers.CleanupNetwork(t, n)

	// Create database container
	dbContainer := createDatabaseContainer(ctx, t, n.Name)

	// Create application container
	appPortMapped := fmt.Sprintf(natPortFormat, appPort)
	appNatPort := nat.Port(appPortMapped)

	adminUser, _ := str.GenerateSecureRandomString(8)
	adminPass, _ := str.GenerateSecureRandomString(8)
	secret, _ := str.GenerateSecureRandomString(32)
	sessionSecret, _ := str.GenerateSecureRandomString(16)

	appContainerReq := testcontainers.ContainerRequest{
		Image:          image,
		ExposedPorts:   []string{appPortMapped},
		Networks:       []string{n.Name},
		LogConsumerCfg: NewVerboseLogConsumerConfig(meta.Name),
		Env: map[string]string{
			"AUTH_SESSION_USER":     adminUser,
			"AUTH_SESSION_PASSWORD": adminPass,
			"AUTH_SESSION_SECRET":   sessionSecret,
			"DB_POSTGRES_HOST":      postgresNetworkAlias,
			"DB_POSTGRES_PORT":      postgresPort,
			"DB_POSTGRES_USER":      postgresUser,
			"DB_POSTGRES_PASSWORD":  postgresPass,
			"DB_POSTGRES_NAME":      postgresDb,
			"SECRET":                secret,
			"WEB_INTERFACE_ENABLED": "false",
		},
		WaitingFor: wait.ForAll(
			wait.ForHTTP("/api/v1/health").WithPort(appNatPort),
			wait.ForExposedPort(),
		).WithStartupTimeoutDefault(containerStartupTimeout).
			WithDeadline(containerDeadlineTimeout),
	}

	var appContainer testcontainers.Container
	if appContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: appContainerReq,
		Started:          true,
	}); err != nil {
		t.Fatalf("failed to start application container: %v", err)
	}
	testcontainers.CleanupContainer(t, appContainer)

	// get the host and port for the application container
	var appHost string
	if appHost, err = appContainer.Host(ctx); err != nil {
		t.Fatalf("failed to get application container host: %v", err)
	}

	var mappedAppPort nat.Port
	if mappedAppPort, err = appContainer.MappedPort(ctx, appPort); err != nil {
		t.Fatalf("failed to get application container port: %v", err)
	}
	appURL := fmt.Sprintf("http://%s:%s", appHost, mappedAppPort.Port())

	return dbContainer, appContainer, appURL
}

func TestRequestInfoEndpoint(t *testing.T) {
	ctx := context.Background()

	dbContainer, appContainer, appURL := setupTestEnvironment(t)
	defer func() {
		if err := dbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate database container: %v", err)
		}
		if err := appContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate application container: %v", err)
		}
	}()

	resp, err := http.Get(appURL + "/api/v1/info")
	defer func(b io.ReadCloser) {
		if bodyCloseErr := b.Close(); bodyCloseErr != nil {
			t.Fatalf("cannot close body %v", bodyCloseErr)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	i := &api.InfoResponse{}
	p := &api.DataResponse{Data: i}
	err = json.Unmarshal(body, p)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, meta.Name, i.Name)
	assert.NotEmpty(t, i.TimeZone)
	assert.NotEmpty(t, i.Version)
}
