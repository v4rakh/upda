package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

var (
	errAppNotReady = errors.New("application did not become ready in time")
)

type testAppReadyFunc func() bool

// testAppHTTPReadyCheck returns an testAppReadyFunc that checks the a endpoint to determine if testAppSpawner is ready
func testAppHTTPReadyCheck(healthCheckUrl string, expectedStatus int) testAppReadyFunc {
	return func() bool {
		resp, err := http.Get(healthCheckUrl)
		if err != nil {
			return false
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		return resp.StatusCode == expectedStatus
	}
}

// testAppSpawner defines an interface for starting and stopping applications in tests, as well as checking their readiness.
type testAppSpawner interface {
	Name() string
	Start(maxWait time.Duration, interval time.Duration, t *testing.T) error
	Stop(maxWait time.Duration) error
}

// testBinaryApplication is a simple implementation of the testAppSpawner interface that starts a binary application as a subprocess and checks its readiness using a provided testAppReadyFunc.
type testBinaryApplication struct {
	appName      string
	cmd          *exec.Cmd
	cmdArgs      []string
	env          []string
	checkReadyFn testAppReadyFunc
	ctx          context.Context
	cancel       context.CancelFunc
}

// Name returns the name of the application.
func (a *testBinaryApplication) Name() string {
	return a.appName
}

// newTestBinaryApplication creates a new testBinaryApplication with the specified name, binary path, environment variables, and readiness check function.
func newTestBinaryApplication(name string, binaryPath string, env []string, checkReadyFn testAppReadyFunc, cmdArgs ...string) testAppSpawner {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, binaryPath, cmdArgs...)

	return &testBinaryApplication{
		appName:      name,
		cmd:          cmd,
		cmdArgs:      cmdArgs,
		env:          env,
		checkReadyFn: checkReadyFn,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start starts the application as a subprocess, sets up logging to the testing.T output, and waits for the application to become ready using the provided testAppReadyFunc.
func (a *testBinaryApplication) Start(maxWait time.Duration, interval time.Duration, t *testing.T) error {
	a.cmd.Env = append(os.Environ(), a.env...)

	logger := &appVerboseLogWriter{Name: a.appName, t: t}
	a.cmd.Stdout = logger
	a.cmd.Stderr = logger

	if err := a.cmd.Start(); err != nil {
		return err
	}

	// wait for ready
	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		if a.checkReadyFn() {
			return nil
		}
		time.Sleep(interval)
	}

	return errAppNotReady
}

// Stop attempts to gracefully stop the application by sending an interrupt signal.
func (a *testBinaryApplication) Stop(maxWait time.Duration) error {
	if a.cmd.Process != nil {
		_ = a.cmd.Process.Signal(os.Interrupt)
		done := make(chan error, 1)
		go func() {
			done <- a.cmd.Wait()
		}()
		select {
		case <-done:
		case <-time.After(maxWait):
			_ = a.cmd.Process.Kill()
		}
	}
	a.cancel()
	return nil
}

// appVerboseLogWriter is an io.Writer that writes logs to the testing.T log output, prefixed with the application name.
type appVerboseLogWriter struct {
	Name string
	t    *testing.T
}

// Write implements the io.Writer interface, writing the log output to the testing.T log with the application name as a prefix.
func (w *appVerboseLogWriter) Write(p []byte) (n int, err error) {
	w.t.Logf("App[%s]:\t%s", w.Name, p)
	return len(p), nil
}

// containerVerboseLogConsumer is a testcontainers.LogConsumer that writes container logs to the standard output, prefixed with the container name.
type containerVerboseLogConsumer struct {
	name string
}

// newContainerVerboseLogConsumer creates a new containerVerboseLogConsumer with the specified container name.
func newContainerVerboseLogConsumer(name string) *containerVerboseLogConsumer {
	return &containerVerboseLogConsumer{name: name}
}

// Accept implements the testcontainers.LogConsumer interface, writing the log output to the standard output with the container name as a prefix.
func (g *containerVerboseLogConsumer) Accept(l testcontainers.Log) {
	fmt.Printf("Container[%s]:\t%s", g.name, string(l.Content))
}

// newContainerVerboseLogConsumerConfig creates a new testcontainers.LogConsumerConfig that uses a containerVerboseLogConsumer for the specified container name and log production timeout.
func newContainerVerboseLogConsumerConfig(name string, timeout time.Duration) *testcontainers.LogConsumerConfig {
	return &testcontainers.LogConsumerConfig{
		Consumers: []testcontainers.LogConsumer{newContainerVerboseLogConsumer(name)},
		Opts: []testcontainers.LogProductionOption{
			testcontainers.WithLogProductionTimeout(timeout),
		},
	}
}

// randomUserPort generates a random port number in the user port range (1024-65535) to avoid conflicts with well-known ports during testing.
func randomUserPort() int {
	return rand.Intn(65535-1024+1) + 1024
}
