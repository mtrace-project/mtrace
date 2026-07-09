package analytics

import (
	"strings"
	"testing"
)

func TestNewAnalyticsExporters(t *testing.T) {
	t.Run("all valid formats", func(t *testing.T) {
		formats := []string{JSON_FORMAT, HTML_FORMAT}
		exporters, err := NewAnalyticsExporters(formats, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(exporters) != 2 {
			t.Errorf("expected 2 exporters, got %d", len(exporters))
		}

		foundJSON := false
		foundHTML := false

		for _, e := range exporters {
			switch e.Format() {
			case JSON_FORMAT:
				foundJSON = true
			case HTML_FORMAT:
				foundHTML = true
			}
		}

		if !foundJSON || !foundHTML {
			t.Errorf("missing some exporters: json=%v html=%v", foundJSON, foundHTML)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		formats := []string{"invalid_format"}
		exporters, err := NewAnalyticsExporters(formats, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if exporters != nil {
			t.Fatal("expected nil exporters on error")
		}
		if !strings.Contains(err.Error(), "unsupported analytics format") {
			t.Errorf("expected unsupported format error, got %v", err)
		}
	})

	t.Run("empty formats", func(t *testing.T) {
		exporters, err := NewAnalyticsExporters(nil, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(exporters) != 0 {
			t.Errorf("expected 0 exporters, got %d", len(exporters))
		}
	})
}
