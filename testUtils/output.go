package testutils

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"
)

// CaptureStdout captures the standard output generated during the execution of f.
// It forces the slog logger to LevelDebug for the duration of the capture.
func CaptureStdout(t *testing.T, f func()) string {
	t.Helper()
	return CaptureStdoutWithLevel(t, slog.LevelDebug, f)
}

// CaptureStdoutWithLevel captures the standard output generated during the execution of f,
// and sets the slog logger to the specified level for the duration of the capture.
func CaptureStdoutWithLevel(t *testing.T, level slog.Level, f func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stdout capture: %v", err)
	}
	os.Stdout = w

	// Set slog default to write to the pipe w
	oldSlogLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})))

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	f()
	_ = w.Close()
	os.Stdout = old

	// Restore slog default
	slog.SetDefault(oldSlogLogger)

	return <-outChan
}
