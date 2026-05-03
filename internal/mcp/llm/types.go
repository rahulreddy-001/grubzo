package llm

import "context"

type ToolDefinition struct {
	Name        string
	Description string
	InputSchema map[string]any
}

type Message struct {
	Role       string
	Text       string
	ToolCallID string
	ToolName   string
	ToolInput  map[string]any
	Meta       map[string]any
	IsError    bool
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]any
	Meta      map[string]any
}

type Request struct {
	SystemPrompt string
	Messages     []Message
	Tools        []ToolDefinition
	Model        string
	Temperature  float64
	MaxTokens    int
}

type Response struct {
	Text       string
	ToolCalls  []ToolCall
	StopReason string
}

type Provider interface {
	Name() string
	DefaultModel() string
	Generate(ctx context.Context, request Request) (*Response, error)
}
