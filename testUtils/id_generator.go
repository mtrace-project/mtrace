package testutils

type MockIdGenerator struct {
	TraceID string
	SpanID  string
	Err     error
}

func (m *MockIdGenerator) Generate(length int) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	if length == 32 {
		return m.TraceID, nil
	}
	return m.SpanID, nil
}
