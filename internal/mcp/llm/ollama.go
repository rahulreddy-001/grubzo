package llm

import (
	"context"
	"encoding/json"
	"fmt"
)

type ollamaProvider struct {
	httpProvider
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Tools    []openAITool    `json:"tools,omitempty"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content,omitempty"`
	ToolCalls []openAIToolCall `json:"tool_calls,omitempty"`
}

type ollamaResponse struct {
	Message struct {
		Content   string           `json:"content"`
		ToolCalls []openAIToolCall `json:"tool_calls"`
	} `json:"message"`
}

func (p *ollamaProvider) Generate(ctx context.Context, request Request) (*Response, error) {
	model := request.Model
	if model == "" {
		model = p.model
	}
	if model == "" {
		return nil, fmt.Errorf("ollama model is not configured")
	}

	payload := ollamaRequest{
		Model:    model,
		Messages: ollamaMessagesFromHistory(request.SystemPrompt, request.Messages),
		Stream:   false,
	}
	for _, tool := range request.Tools {
		payload.Tools = append(payload.Tools, openAITool{
			Type: "function",
			Function: openAIFunctionTool{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	var resp ollamaResponse
	if err := p.postJSON(
		ctx,
		p.baseURL(),
		map[string]string{},
		payload,
		&resp,
	); err != nil {
		return nil, err
	}
	out := &Response{Text: resp.Message.Content}
	for _, toolCall := range resp.Message.ToolCalls {
		var args map[string]any
		if toolCall.Function.Arguments != "" {
			_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		}
		out.ToolCalls = append(out.ToolCalls, ToolCall{
			ID:        toolCall.ID,
			Name:      toolCall.Function.Name,
			Arguments: args,
		})
	}
	return out, nil
}

func (p *ollamaProvider) baseURL() string {
	if p.base == "" {
		return "http://localhost:11434/api/chat"
	}
	return p.base
}

func ollamaMessagesFromHistory(systemPrompt string, history []Message) []ollamaMessage {
	out := make([]ollamaMessage, 0, len(history)+1)
	if systemPrompt != "" {
		out = append(out, ollamaMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	for _, message := range history {
		switch {
		case message.Role == "user":
			out = append(out, ollamaMessage{Role: "user", Content: message.Text})
		case message.Role == "assistant" && message.ToolName == "":
			out = append(out, ollamaMessage{Role: "assistant", Content: message.Text})
		case message.Role == "assistant" && message.ToolName != "":
			args, _ := json.Marshal(message.ToolInput)
			out = append(out, ollamaMessage{
				Role: "assistant",
				ToolCalls: []openAIToolCall{{
					ID:   message.ToolCallID,
					Type: "function",
					Function: openAIToolCallFunction{
						Name:      message.ToolName,
						Arguments: string(args),
					},
				}},
			})
		case message.Role == "tool":
			out = append(out, ollamaMessage{
				Role:    "tool",
				Content: message.Text,
			})
		}
	}
	return out
}
