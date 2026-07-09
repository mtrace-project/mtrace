package analytics

import (
	"fmt"
)

// Metrics based on trace durations
type TraceDurationMetricCalculator struct{}

func (c TraceDurationMetricCalculator) Calculate(data *tracesData, traceAnalytics *TraceAnalytics) error {
	if traceAnalytics == nil {
		return fmt.Errorf("traceAnalytics is nil")
	}
	if data == nil {
		return fmt.Errorf("tracesData is nil")
	}
	if len(data.sortedDurations) == 0 {
		return fmt.Errorf("no trace durations available for calculation")
	}

	traceAnalytics.MinDuration = data.sortedDurations[0]
	traceAnalytics.MaxDuration = data.sortedDurations[len(data.sortedDurations)-1]
	traceAnalytics.P50Duration = data.sortedDurations[getPercentileIndex(len(data.sortedDurations), 50)]
	traceAnalytics.P90Duration = data.sortedDurations[getPercentileIndex(len(data.sortedDurations), 90)]
	traceAnalytics.P99Duration = data.sortedDurations[getPercentileIndex(len(data.sortedDurations), 99)]
	traceAnalytics.DurationStandardDeviation = calculateStandardDeviation(data.sortedDurations)
	traceAnalytics.AverageDuration = calculateAverage(data.sortedDurations)

	return nil
}

func (c TraceDurationMetricCalculator) String() string {
	return "TraceDuration"
}

// Metrics based on span/trace counts
type TraceCountsMetricCalculator struct{}

func (c TraceCountsMetricCalculator) Calculate(data *tracesData, traceAnalytics *TraceAnalytics) error {
	if traceAnalytics == nil {
		return fmt.Errorf("traceAnalytics is nil")
	}
	if data == nil {
		return fmt.Errorf("tracesData is nil")
	}

	traceAnalytics.AverageSpanCount = calculateAverage(data.spanCountPerTrace)
	traceAnalytics.AverageSpanErrorCount = calculateAverage(data.errorCountPerTrace)
	traceAnalytics.ErrorRate = calculateErrorRate(data.errorCountPerTrace)

	return nil
}

func (c TraceCountsMetricCalculator) String() string {
	return "TraceCounts"
}
