package trigger_test

import (
	"context"
	"testing"

	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/parser"
	testutils "gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/testUtils"
	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/trigger"
)

func TestUnsupportedTrigger(t *testing.T) {
	dto := &parser.TriggerDTO{
		Type: "unknownType",
	}
	_, err := trigger.NewTrigger(dto, &testutils.MockIdGenerator{}, "", context.Background())
	if err == nil {
		t.Error("Expected error for unsupported trigger type")
	}
}
