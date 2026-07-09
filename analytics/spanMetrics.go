package analytics

import (
	"fmt"
	"math"
)

// Metrics based on span durations
type SpanMetricCalculator struct{}

func (c SpanMetricCalculator) Calculate(data *tracesData, traceAnalytics *TraceAnalytics) error {
	if traceAnalytics == nil {
		return fmt.Errorf("traceAnalytics is nil")
	}
	if data == nil {
		return fmt.Errorf("tracesData is nil")
	}

	for key, spanData := range data.spansData {
		if len(spanData.sortedDurations) == 0 {
			continue
		}

		spanAnalytics, ok := traceAnalytics.SpanAnalytics[key]
		if !ok {
			spanAnalytics = &SpanAnalytics{}
			traceAnalytics.SpanAnalytics[key] = spanAnalytics
		}

		spanAnalytics.ServiceName = spanData.serviceName
		spanAnalytics.OperationName = spanData.operationName
		spanAnalytics.MinDuration = spanData.sortedDurations[0]
		spanAnalytics.MaxDuration = spanData.sortedDurations[len(spanData.sortedDurations)-1]
		spanAnalytics.P50Duration = spanData.sortedDurations[getPercentileIndex(len(spanData.sortedDurations), 50)]
		spanAnalytics.P90Duration = spanData.sortedDurations[getPercentileIndex(len(spanData.sortedDurations), 90)]
		spanAnalytics.P99Duration = spanData.sortedDurations[getPercentileIndex(len(spanData.sortedDurations), 99)]
		spanAnalytics.AverageDuration = calculateAverage(spanData.sortedDurations)
		spanAnalytics.DurationStandardDeviation = calculateStandardDeviation(spanData.sortedDurations)
		spanAnalytics.ErrorRate = float64(spanData.errorCount) / float64(spanData.occurencies) * 100
		spanAnalytics.AveragePercentageOfTotalTrace = calculateAveragePercentageOfTotalTrace(spanData.durationsWithParentTrace)
		spanAnalytics.AveragePosition = math.MaxFloat64
		if len(spanData.spanPositionsInTrace) > 0 {
			spanAnalytics.AveragePosition = calculateAverage(spanData.spanPositionsInTrace)
		}
	}

	return nil
}

func (c SpanMetricCalculator) String() string {
	return "Span"
}
