package agent

import "context"

type AgentBackend interface {
	Chat(ctx context.Context, userID, tenantID uint64, conversationID uint64, message string) ([]AgentMessage, error)
	StreamChat(ctx context.Context, tenantID uint64, prompt string) (string, error)
	ExecuteTool(ctx context.Context, tenantID uint64, call ToolCall) (*ToolResult, error)
	ListTools() []ToolDefinition
}
