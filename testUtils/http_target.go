package testutils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// HTTPRequestRecord represents a recorded request received by the mock target server.
type HTTPRequestRecord struct {
	Method  string
	URL     string
	Headers http.Header
	Body    string
}

// HTTPTargetServer represents a mock HTTP target server.
type HTTPTargetServer struct {
	Server   *httptest.Server
	Requests []HTTPRequestRecord
}

// StartHTTPTargetServer starts a mock HTTP target server that records requests.
func StartHTTPTargetServer(t *testing.T, responseBody string, statusCode int) *HTTPTargetServer {
	t.Helper()

	target := &HTTPTargetServer{
		Requests: []HTTPRequestRecord{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		defer r.Body.Close() // nolint:errcheck

		target.Requests = append(target.Requests, HTTPRequestRecord{
			Method:  r.Method,
			URL:     r.URL.String(),
			Headers: r.Header,
			Body:    string(bodyBytes),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	t.Cleanup(func() {
		server.Close()
	})

	target.Server = server
	return target
}
