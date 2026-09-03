package agent_test

import (
	"testing"
	"context"
	"orbit/apps/api/internal/agent"
	"orbit/apps/api/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestNewAgentService(t *testing.T) {
	svc := agent.NewAgentService()
	assert.NotNil(t, svc)
}

func TestExecuteToolUnknown(t *testing.T) {
	svc := agent.NewAgentService()
	result, err := svc.ExecuteTool(context.Background(), 1, agent.ToolCall{Name: "unknown_tool"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown tool")
	assert.Contains(t, result.Error, "unknown tool")
}

func TestExecuteToolGetMetrics(t *testing.T) {
	svc := agent.NewAgentService()
	if database.DB == nil {
		t.Skip("database not connected, skipping get_metrics test")
	}
	result, err := svc.ExecuteTool(context.Background(), 1, agent.ToolCall{Name: "get_metrics"})
	assert.NoError(t, err)
	assert.Equal(t, "get_metrics", result.ToolName)
	assert.Contains(t, result.Result, "total_spend")
}

func TestGetToolDefinitions(t *testing.T) {
	defs := agent.GetToolDefinitions()
	assert.NotEmpty(t, defs)
	assert.Equal(t, "get_metrics", defs[0]["name"])
}
