package analytics

import (
	"testing"
)

func TestSpanMetricCalculator_Calculate(t *testing.T) {
	calculator := SpanMetricCalculator{}

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

	t.Run("empty spans skip", func(t *testing.T) {
		data := &tracesData{
			spansData: map[string]*spansData{
				"srv-op": {
					sortedDurations: []int64{},
				},
			},
		}
		analytics := &TraceAnalytics{
			SpanAnalytics: make(map[string]*SpanAnalytics),
		}
		err := calculator.Calculate(data, analytics)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(analytics.SpanAnalytics) != 0 {
			t.Errorf("expected empty SpanAnalytics map, got %v", len(analytics.SpanAnalytics))
		}
	})

	t.Run("valid spans", func(t *testing.T) {
		data := &tracesData{
			spansData: map[string]*spansData{
				"srv-op": {
					serviceName:     "srv",
					operationName:   "op",
					sortedDurations: []int64{10, 20, 30, 40, 50},
					occurencies:     5,
					errorCount:      2,
					durationsWithParentTrace: []*spanTraceDuration{
						{duration: 10, traceDuration: 100},
						{duration: 20, traceDuration: 100},
						{duration: 30, traceDuration: 100},
						{duration: 40, traceDuration: 100},
						{duration: 50, traceDuration: 100},
					},
				},
			},
		}
		analytics := &TraceAnalytics{
			SpanAnalytics: make(map[string]*SpanAnalytics),
		}
		err := calculator.Calculate(data, analytics)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		spanAnalytics := analytics.SpanAnalytics["srv-op"]
		if spanAnalytics == nil {
			t.Fatalf("expected spanAnalytics to be populated")
		}

		if spanAnalytics.ServiceName != "srv" {
			t.Errorf("expected ServiceName 'srv', got %s", spanAnalytics.ServiceName)
		}
		if spanAnalytics.MinDuration != 10 {
			t.Errorf("expected MinDuration 10, got %d", spanAnalytics.MinDuration)
		}
		if spanAnalytics.MaxDuration != 50 {
			t.Errorf("expected MaxDuration 50, got %d", spanAnalytics.MaxDuration)
		}
		if spanAnalytics.P50Duration != 30 {
			t.Errorf("expected P50Duration 30, got %d", spanAnalytics.P50Duration)
		}
		if spanAnalytics.ErrorRate != 40.0 {
			t.Errorf("expected ErrorRate 40.0, got %f", spanAnalytics.ErrorRate)
		}
		if spanAnalytics.AveragePercentageOfTotalTrace != 30.0 {
			t.Errorf("expected AveragePercentageOfTotalTrace 30.0, got %f", spanAnalytics.AveragePercentageOfTotalTrace)
		}
	})
}
