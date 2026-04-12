package server

import (
	"context"
	"encoding/json"
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/str"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	envTestDir    = "TEST_DIR"
	envTestBinary = "TEST_BINARY"

	appHost = "localhost"

	appStartupReadyTimeout  = 10 * time.Second
	appStartupRetryInterval = 500 * time.Millisecond
	appStopTimeout          = 10 * time.Second

	containerStartupTimeout = 20 * time.Second
	containerLogTimeout     = 10 * time.Second

	postgresContainerName = "db"
	postgresPort          = "5432"
	postgresImage         = "postgres:17-alpine"
	postgresUser          = "upda"
	postgresPass          = "upda"
	postgresDb            = "upda"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	os.Exit(m.Run())
}

func createDatabaseContainer(ctx context.Context, t *testing.T) testcontainers.Container {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:          postgresImage,
		ExposedPorts:   []string{postgresPort},
		LogConsumerCfg: newContainerVerboseLogConsumerConfig(postgresContainerName, containerLogTimeout),
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPass,
			"POSTGRES_DB":       postgresDb,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithStartupTimeout(containerStartupTimeout),
			wait.ForSQL("5432/tcp", "postgres",
				func(host string, port string) string {
					strippedPort := stripPortSuffix(port)
					return fmt.Sprintf(
						"postgres://%s:%s@%s:%s/%s?sslmode=disable",
						postgresUser, postgresPass, host, strippedPort, postgresDb,
					)
				},
			).WithStartupTimeout(containerStartupTimeout),
		),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start database container: %v", err)
	}

	return c
}

func TestApplication(t *testing.T) {
	ctx := context.Background()
	var err error

	// test prerequisites
	testBinaryPath := os.Getenv(envTestBinary)
	if testBinaryPath == "" {
		t.Skipf("'%s' environment not set", envTestBinary)
	}
	coverageDir := os.Getenv(envTestDir)
	if coverageDir == "" {
		t.Skipf("'%s' environment not set", envTestDir)
	}
	if _, err := os.Stat(testBinaryPath); os.IsNotExist(err) {
		t.Skipf("test binary not found at '%s'", testBinaryPath)
	}

	// setup prerequisites
	dbContainer := createDatabaseContainer(ctx, t)
	defer func() {
		if err := dbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate database container: %v", err)
		}
	}()
	dbMappedHost, err := dbContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get database container host: %v", err)
	}
	dbMappedPort, err := dbContainer.MappedPort(ctx, postgresPort)
	if err != nil {
		t.Fatalf("failed to get database container mapped port: %v", err)
	}
	t.Logf("Using database container '%s:%s'", dbMappedHost, dbMappedPort.Port())

	// setup actual application
	adminUser, _ := str.GenerateSecureRandomString(8)
	adminPass, _ := str.GenerateSecureRandomString(8)
	secret, _ := str.GenerateSecureRandomString(32)
	sessionSecret, _ := str.GenerateSecureRandomString(16)

	appPort := randomUserPort()
	appURL := fmt.Sprintf("http://%s:%d", appHost, appPort)
	checkReadyFn := testAppHTTPReadyCheck(fmt.Sprintf("%s/api/v1/health", appURL), http.StatusOK)

	app := newTestBinaryApplication(
		meta.Name,
		testBinaryPath,
		[]string{
			fmt.Sprintf("GOCOVERDIR=%s", coverageDir),
			fmt.Sprintf("SERVER_PORT=%d", appPort),
			fmt.Sprintf("DB_POSTGRES_HOST=%s", dbMappedHost),
			fmt.Sprintf("DB_POSTGRES_PORT=%s", dbMappedPort.Port()),
			fmt.Sprintf("DB_POSTGRES_USER=%s", postgresUser),
			fmt.Sprintf("DB_POSTGRES_PASSWORD=%s", postgresPass),
			fmt.Sprintf("DB_POSTGRES_NAME=%s", postgresDb),
			fmt.Sprintf("AUTH_SESSION_USER=%s", adminUser),
			fmt.Sprintf("AUTH_SESSION_PASSWORD=%s", adminPass),
			fmt.Sprintf("AUTH_SESSION_SECRET=%s", sessionSecret),
			fmt.Sprintf("SECRET=%s", secret),
			"WEB_INTERFACE_ENABLED=false",
			"PROMETHEUS_ENABLED=true",
			"PROMETHEUS_SECURE_TOKEN_ENABLED=false",
			fmt.Sprintf("PROMETHEUS_PORT=%d", appPort),
			"LOGGING_LEVEL=debug",
			"LOGGING_LEVEL_REQUESTS=debug",
			"LOGGING_ENCODING=console",
		},
		checkReadyFn,
		"server",
		"serve",
	)
	if err = app.Start(appStartupReadyTimeout, appStartupRetryInterval, t); err != nil {
		t.Fatalf("failed to start application: %v", err)
	}
	defer func(app testAppSpawner) {
		if errStop := app.Stop(appStopTimeout); errStop != nil {
			t.Logf("failed to stop application: %s", errStop)
		}
	}(app)

	// Run tests
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/info", appURL))
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
	assert.NotNil(t, i.Version)
}
