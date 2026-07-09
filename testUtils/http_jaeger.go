package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// StartMockJaegerServer starts a mock HTTP server that simulates Jaeger trace search endpoint.
func StartMockJaegerServer(t *testing.T, handler func(traceID string) (any, int)) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		const prefix = "/api/traces/"
		if !strings.HasPrefix(path, prefix) {
			http.Error(w, "invalid endpoint", http.StatusNotFound)
			return
		}

		traceID := strings.TrimPrefix(path, prefix)
		data, statusCode := handler(traceID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if statusCode == http.StatusOK {
			resp := map[string]any{
				"data":   data,
				"total":  0,
				"limit":  0,
				"offset": 0,
				"errors": nil,
			}
			_ = json.NewEncoder(w).Encode(resp)
		} else {
			_, _ = w.Write([]byte("mock error"))
		}
	}))

	t.Cleanup(func() {
		server.Close()
	})

	return server
}
