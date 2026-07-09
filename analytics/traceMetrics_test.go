package analytics

import (
	"testing"
)

func TestTraceDurationMetricCalculator_Calculate(t *testing.T) {
	calculator := TraceDurationMetricCalculator{}

	t.Run("nil traceAnalytics", func(t *testing.T) {
		err := calculator.Calculate(&tracesData{}, nil)
		if err == nil {
			t.Errorf("expected error when traceAnalytics is nil")
		}
	})

	t.Run("nil data", func(t *testing.T) {
		err := calculator.Calculate(nil, &TraceAnalytics{})
		if err == nil {
			t.Errorf("expected error when data is nil")
		}
	})

	t.Run("empty sortedDurations", func(t *testing.T) {
		data := &tracesData{
			sortedDurations: []int64{},
		}
		err := calculator.Calculate(data, &TraceAnalytics{})
		if err == nil || err.Error() != "no trace durations available for calculation" {
			t.Errorf("expected specific error for empty durations, got: %v", err)
		}
	})

	t.Run("valid durations", func(t *testing.T) {
		data := &tracesData{
			sortedDurations: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		}
		analytics := &TraceAnalytics{}
		err := calculator.Calculate(data, analytics)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if analytics.MinDuration != 10 {
			t.Errorf("expected MinDuration 10, got %d", analytics.MinDuration)
		}
		if analytics.MaxDuration != 100 {
			t.Errorf("expected MaxDuration 100, got %d", analytics.MaxDuration)
		}
		if analytics.P50Duration != 50 {
			t.Errorf("expected P50Duration 50, got %d", analytics.P50Duration)
		}
		if analytics.P90Duration != 90 {
			t.Errorf("expected P90Duration 90, got %d", analytics.P90Duration)
		}
		if analytics.P99Duration != 100 {
			t.Errorf("expected P99Duration 100, got %d", analytics.P99Duration)
		}
		if analytics.AverageDuration != 55.0 {
			t.Errorf("expected AverageDuration 55, got %f", analytics.AverageDuration)
		}
	})
}

func TestTraceCountsMetricCalculator_Calculate(t *testing.T) {
	calculator := TraceCountsMetricCalculator{}

	t.Run("nil traceAnalytics", func(t *testing.T) {
		err := calculator.Calculate(&tracesData{}, nil)
		if err == nil {
			t.Errorf("expected error when traceAnalytics is nil")
		}
	})

	t.Run("nil data", func(t *testing.T) {
		err := calculator.Calculate(nil, &TraceAnalytics{})
		if err == nil {
			t.Errorf("expected error when data is nil")
		}
	})

	t.Run("valid counts", func(t *testing.T) {
		data := &tracesData{
			spanCountPerTrace:  []int{5, 10, 15},
			errorCountPerTrace: []int{0, 2, 0},
		}
		analytics := &TraceAnalytics{}
		err := calculator.Calculate(data, analytics)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if analytics.AverageSpanCount != 10.0 {
			t.Errorf("expected AverageSpanCount 10, got %f", analytics.AverageSpanCount)
		}
		if analytics.AverageSpanErrorCount != float64(2)/3 {
			t.Errorf("expected AverageSpanErrorCount 0.666, got %f", analytics.AverageSpanErrorCount)
		}
		// 1 trace with error out of 3 total traces = 33.333%
		if analytics.ErrorRate < 33.3 || analytics.ErrorRate > 33.4 {
			t.Errorf("expected ErrorRate around 33.33, got %f", analytics.ErrorRate)
		}
	})
}
