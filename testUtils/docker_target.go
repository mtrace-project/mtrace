package testutils

import (
	"net/http"
	"testing"

	"github.com/moby/moby/client"
)

type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func NewMockDockerClient(t *testing.T, roundTripFunc func(req *http.Request) (*http.Response, error)) *client.Client {
	httpClient := &http.Client{
		Transport: &MockRoundTripper{RoundTripFunc: roundTripFunc},
	}
	cli, err := client.New(
		client.WithHTTPClient(httpClient),
	)
	if err != nil {
		t.Fatalf("failed to create mock docker client: %v", err)
	}
	return cli
}
