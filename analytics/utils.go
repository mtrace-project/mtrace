package analytics

import "math"

func getPercentileIndex(n int, percentile float64) int {
	// Nearest-Rank: ceil((p / 100) * n) - 1
	idx := int(math.Ceil((percentile/100.0)*float64(n))) - 1

	if idx < 0 {
		return 0
	}
	if idx >= n {
		return n - 1
	}
	return idx
}

func calculateStandardDeviation(durations []int64) float64 {
	mean := calculateAverage(durations)
	var variance float64
	// Sqrt((Sum((x_i - mean)^2)) / N), where i is the index of the duration in the slice and N is the total number of durations
	for _, duration := range durations {
		diff := float64(duration) - mean
		variance += diff * diff
	}
	variance /= float64(len(durations))
	return math.Sqrt(variance)
}

type number interface {
	~int | ~int64 | ~int32 | ~float64
}

func calculateAverage[T number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum T
	for _, value := range values {
		sum += value
	}

	return float64(sum) / float64(len(values))
}

// Percentage of traces with at least one span with an error status
func calculateErrorRate(errorCounts []int) float64 {
	if len(errorCounts) == 0 {
		return 0
	}

	var errorSamples int
	for _, errorCount := range errorCounts {
		if errorCount > 0 {
			errorSamples++
		}
	}

	return (float64(errorSamples) / float64(len(errorCounts))) * 100
}

// Average percentage of the span duration over the relative total trace duration
func calculateAveragePercentageOfTotalTrace(durationsWithParentTrace []*spanTraceDuration) float64 {
	if len(durationsWithParentTrace) == 0 {
		return 0
	}

	var percentageSum float64
	for _, d := range durationsWithParentTrace {
		if d.traceDuration == 0 {
			continue
		}
		percentageSum += float64(d.duration) / float64(d.traceDuration) * 100
	}

	return percentageSum / float64(len(durationsWithParentTrace))
}
