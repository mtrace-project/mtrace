package export

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewExporters(t *testing.T) { //nolint:gocyclo
	t.Run("all valid formats", func(t *testing.T) {
		formats := []string{JSON_FORMAT, JUNIT_FORMAT, MARKDOWN_FORMAT, STDOUT_FORMAT}
		exporters, err := NewExporters(formats, 1, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(exporters) != 4 {
			t.Errorf("expected 4 exporters, got %d", len(exporters))
		}

		foundJSON := false
		foundJunit := false
		foundMarkdown := false
		foundStdout := false

		for _, e := range exporters {
			switch e.Format() {
			case JSON_FORMAT:
				foundJSON = true
			case JUNIT_FORMAT:
				foundJunit = true
			case MARKDOWN_FORMAT:
				foundMarkdown = true
			case STDOUT_FORMAT:
				foundStdout = true
			}
		}

		if !foundJSON || !foundJunit || !foundMarkdown || !foundStdout {
			t.Errorf("missing some exporters: json=%v junit=%v markdown=%v stdout=%v", foundJSON, foundJunit, foundMarkdown, foundStdout)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		formats := []string{"invalid_format"}
		exporters, err := NewExporters(formats, 1, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if exporters != nil {
			t.Fatal("expected nil exporters on error")
		}
		if !strings.Contains(err.Error(), "unsupported export format") {
			t.Errorf("expected unsupported format error, got %v", err)
		}
	})

	t.Run("empty formats", func(t *testing.T) {
		exporters, err := NewExporters(nil, 1, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(exporters) != 0 {
			t.Errorf("expected 0 exporters, got %d", len(exporters))
		}
	})

	t.Run("with baseDir", func(t *testing.T) {
		formats := []string{JSON_FORMAT, JUNIT_FORMAT, MARKDOWN_FORMAT}
		baseDir := "/custom/base/dir"
		exporters, err := NewExporters(formats, 1, baseDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, e := range exporters {
			switch exp := e.(type) {
			case *jsonExporter:
				expected := filepath.Join(baseDir, "json_exports")
				if exp.outputFolder != expected {
					t.Errorf("expected json outputFolder %s, got %s", expected, exp.outputFolder)
				}
			case *junitExporter:
				expected := filepath.Join(baseDir, "junit_exports")
				if exp.outputFolder != expected {
					t.Errorf("expected junit outputFolder %s, got %s", expected, exp.outputFolder)
				}
			case *markdownExporter:
				expected := filepath.Join(baseDir, "markdown_exports")
				if exp.outputFolder != expected {
					t.Errorf("expected markdown outputFolder %s, got %s", expected, exp.outputFolder)
				}
			}
		}
	})
}
