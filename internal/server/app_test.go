//go:build integration

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/commons"
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
	postgresNetworkAlias     = "db"
	postgresPort             = "5432"
	postgresImage            = "postgres:17-alpine"
	postgresUser             = "upda"
	postgresPass             = "upda"
	postgresDb               = "upda"
	appPort                  = "8080"
	imageEnvKey              = "IMAGE"
	containerStartupTimeout  = 20 * time.Second
	containerDeadlineTimeout = 120 * time.Second
	containerLogTimeout      = 10 * time.Second
)

type VerboseLogConsumer struct{}

func (g *VerboseLogConsumer) Accept(l testcontainers.Log) {
	fmt.Printf("Container Log: %s", string(l.Content))
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

	// Create loggers for the containers
	logConsumer := &VerboseLogConsumer{}
	logConsumerConfig := &testcontainers.LogConsumerConfig{
		Consumers: []testcontainers.LogConsumer{logConsumer},
		Opts: []testcontainers.LogProductionOption{
			testcontainers.WithLogProductionTimeout(containerLogTimeout),
		},
	}

	// Create database container
	const natPortFormat = "%s/tcp"
	dbContainerReq := testcontainers.ContainerRequest{
		Image:        postgresImage,
		ExposedPorts: []string{fmt.Sprintf(natPortFormat, postgresPort)},
		Networks:     []string{n.Name},
		NetworkAliases: map[string][]string{
			n.Name: {postgresNetworkAlias},
		},
		LogConsumerCfg: logConsumerConfig,
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

	var dbContainer testcontainers.Container
	if dbContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: dbContainerReq,
		Started:          true,
	}); err != nil {
		t.Fatalf("failed to start database container: %v", err)
	}
	testcontainers.CleanupContainer(t, dbContainer)

	// Create application container
	appPortMapped := fmt.Sprintf(natPortFormat, appPort)
	appNatPort := nat.Port(appPortMapped)

	adminUser, _ := commons.GenerateSecureRandomString(8)
	adminPass, _ := commons.GenerateSecureRandomString(8)
	secret, _ := commons.GenerateSecureRandomString(32)

	appContainerReq := testcontainers.ContainerRequest{
		Image:          image,
		ExposedPorts:   []string{appPortMapped},
		Networks:       []string{n.Name},
		LogConsumerCfg: logConsumerConfig,
		Env: map[string]string{
			"DB_POSTGRES_HOST":               postgresNetworkAlias,
			"DB_POSTGRES_PORT":               postgresPort,
			"DB_POSTGRES_USER":               postgresUser,
			"DB_POSTGRES_PASSWORD":           postgresPass,
			"DB_POSTGRES_NAME":               postgresDb,
			"SECRET":                         secret,
			"EMBEDDED_WEB_INTERFACE_ENABLED": "false",
			"BASIC_AUTH_USER":                adminUser,
			"BASIC_AUTH_PASSWORD":            adminPass,
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
	defer resp.Body.Close()

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
	assert.Equal(t, commons.Name, i.Name)
	assert.NotEmpty(t, i.TimeZone)
	assert.NotEmpty(t, i.Version)
}
