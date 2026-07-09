package testutils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// StartMockOpenObserveServer starts a mock HTTP server that simulates OpenObserve trace search endpoint.
// The handler receives the SQL query string and returns the list of hits and HTTP status code.
func StartMockOpenObserveServer(t *testing.T, handler func(sqlQuery string) ([]map[string]any, int)) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close() // nolint:errcheck

		var req struct {
			Query struct {
				SQL string `json:"sql"`
			} `json:"query"`
		}

		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}

		hits, statusCode := handler(req.Query.SQL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if statusCode == http.StatusOK {
			resp := map[string]any{
				"hits": hits,
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
