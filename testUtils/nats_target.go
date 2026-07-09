package testutils

import (
	"strings"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

// StartMockNATSServer starts an in-memory NATS server on a random port.
// If enableJetstream is true, JetStream is initialized using a temp directory.
// Returns the connection address (without the "nats://" scheme).
func StartMockNATSServer(t *testing.T, enableJetstream bool) string {
	t.Helper()

	opts := &natsserver.Options{
		Host:   "127.0.0.1",
		Port:   -1, // Use random port
		NoLog:  true,
		NoSigs: true,
	}

	if enableJetstream {
		opts.JetStream = true
		opts.StoreDir = t.TempDir()
	}

	return StartMockNATSServerWithOptions(t, opts)
}

// StartMockNATSServerWithOptions starts an in-memory NATS server with the given Options.
// Returns the connection address (without the "nats://" scheme).
func StartMockNATSServerWithOptions(t *testing.T, opts *natsserver.Options) string {
	t.Helper()

	if opts.Host == "" {
		opts.Host = "127.0.0.1"
	}
	if opts.Port == 0 {
		opts.Port = -1
	}

	ns, err := natsserver.NewServer(opts)
	if err != nil {
		t.Fatalf("Failed to create NATS server: %v", err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		t.Fatalf("NATS server failed to start")
	}

	t.Cleanup(func() {
		ns.Shutdown()
		ns.WaitForShutdown()
	})

	url := ns.ClientURL()
	return strings.TrimPrefix(url, "nats://")
}
