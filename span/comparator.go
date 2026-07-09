package span

import (
	"fmt"
	"strings"
)

func (s *Span) Equal(expected *ExpectedSpan) (bool, string) {
	if expected == nil {
		return false, "expected span is nil"
	}

	if s == nil {
		return false, "actual span is nil"
	}

	if expected.ServiceName != "" && s.ServiceName != "" && !strings.EqualFold(expected.ServiceName, s.ServiceName) {
		return false, fmt.Sprintf("service name does not match, expected: %s, actual: %s.\n Expected: %v, Actual: %v", expected.ServiceName, s.ServiceName, expected, s)
	}
	if expected.OperationName != nil && s.OperationName != "" && !strings.EqualFold(*expected.OperationName, s.OperationName) {
		return false, fmt.Sprintf("operation name does not match, expected: %s, actual: %s.\n Expected: %v, Actual: %v", *expected.OperationName, s.OperationName, expected, s)
	}
	if expected.SpanKind != nil && s.SpanKind != "" && !strings.EqualFold(*expected.SpanKind, s.SpanKind) {
		return false, fmt.Sprintf("span kind does not match, expected: %s, actual: %s.\n Expected: %v, Actual: %v", *expected.SpanKind, s.SpanKind, expected, s)
	}
	if expected.SpanStatus != nil && s.SpanStatus != "" && !strings.EqualFold(*expected.SpanStatus, s.SpanStatus) {
		return false, fmt.Sprintf("span status does not match, expected: %s, actual: %s.\n Expected: %v, Actual: %v", *expected.SpanStatus, s.SpanStatus, expected, s)
	}
	if expected.maxDuration != nil && s.Duration > *expected.maxDuration {
		return false, fmt.Sprintf("duration exceeds maximum, expected max: %s, actual: %s.\n Expected: %v, Actual: %v", expected.maxDuration.String(), s.Duration.String(), expected, s)
	}
	if expected.minDuration != nil && s.Duration < *expected.minDuration {
		return false, fmt.Sprintf("duration is less than minimum, expected min: %s, actual: %s.\n Expected: %v, Actual: %v", expected.minDuration.String(), s.Duration.String(), expected, s)
	}
	return true, ""
}
