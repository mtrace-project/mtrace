package assertion

import (
	"fmt"
	"strings"

	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/parser"
	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/trace"
)

type Assertion interface {
	Assert(t *trace.Trace) (bool, error)
}

func NewAssertion(dto *parser.AssertionDTO) (Assertion, error) {
	switch strings.ToLower(dto.Type) {
	case "cel":
		return NewCelAssertion(dto)
	default:
		return nil, fmt.Errorf("unsupported assertion type: %s", dto.Type)
	}
}

func NewAssertions(dtos []*parser.AssertionDTO) ([]Assertion, error) {
	var assertions []Assertion
	for _, dto := range dtos {
		assertion, err := NewAssertion(dto)
		if err != nil {
			return nil, fmt.Errorf("error creating assertion: %w", err)
		}
		assertions = append(assertions, assertion)
	}
	return assertions, nil
}
