package llm

import (
	"context"
	"encoding/json"
	"fmt"
)

type anthropicProvider struct {
	httpProvider
}

type anthropicRequest struct {
	Model       string             `json:"model"`
	System      string             `json:"system,omitempty"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature,omitempty"`
	Messages    []anthropicMessage `json:"messages"`
	Tools       []anthropicTool    `json:"tools,omitempty"`
}

type anthropicTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

type anthropicMessage struct {
	Role    string                  `json:"role"`
	Content []anthropicContentBlock `json:"content"`
}

type anthropicContentBlock struct {
	Type      string         `json:"type"`
	Text      string         `json:"text,omitempty"`
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Input     map[string]any `json:"input,omitempty"`
	ToolUseID string         `json:"tool_use_id,omitempty"`
	IsError   bool           `json:"is_error,omitempty"`
	Content   string         `json:"content,omitempty"`
}

type anthropicResponse struct {
	StopReason string                  `json:"stop_reason"`
	Content    []anthropicContentBlock `json:"content"`
}

func (p *anthropicProvider) Generate(ctx context.Context, request Request) (*Response, error) {
	model := request.Model
	if model == "" {
		model = p.model
	}
	if model == "" {
		return nil, fmt.Errorf("anthropic model is not configured")
	}

	payload := anthropicRequest{
		Model:       model,
		System:      request.SystemPrompt,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		Messages:    anthropicMessagesFromHistory(request.Messages),
	}
	for _, tool := range request.Tools {
		payload.Tools = append(payload.Tools, anthropicTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}
	if payload.MaxTokens <= 0 {
		payload.MaxTokens = 1024
	}

	var resp anthropicResponse
	if err := p.postJSON(
		ctx,
		p.baseURL(),
		map[string]string{
			"x-api-key":         p.apiKey,
			"anthropic-version": "2023-06-01",
		},
		payload,
		&resp,
	); err != nil {
		return nil, err
	}

	out := &Response{StopReason: resp.StopReason}
	for _, block := range resp.Content {
		switch block.Type {
		case "text":
			out.Text += block.Text
		case "tool_use":
			out.ToolCalls = append(out.ToolCalls, ToolCall{
				ID:        block.ID,
				Name:      block.Name,
				Arguments: block.Input,
			})
		}
	}
	return out, nil
}

func (p *anthropicProvider) baseURL() string {
	if p.base == "" {
		return "https://api.anthropic.com/v1/messages"
	}
	if p.base == "https://api.anthropic.com/v1/messages" {
		return p.base
	}
	if p.base == "https://api.anthropic.com/v1" {
		return p.base + "/messages"
	}
	return p.base
}

func anthropicMessagesFromHistory(history []Message) []anthropicMessage {
	out := make([]anthropicMessage, 0, len(history))
	for _, message := range history {
		switch {
		case message.Role == "user":
			out = append(out, anthropicMessage{
				Role: "user",
				Content: []anthropicContentBlock{{
					Type: "text",
					Text: message.Text,
				}},
			})
		case message.Role == "assistant" && message.ToolName == "":
			out = append(out, anthropicMessage{
				Role: "assistant",
				Content: []anthropicContentBlock{{
					Type: "text",
					Text: message.Text,
				}},
			})
		case message.Role == "assistant" && message.ToolName != "":
			out = append(out, anthropicMessage{
				Role: "assistant",
				Content: []anthropicContentBlock{{
					Type:  "tool_use",
					ID:    message.ToolCallID,
					Name:  message.ToolName,
					Input: message.ToolInput,
				}},
			})
		case message.Role == "tool":
			content := message.Text
			if content == "" {
				bytes, _ := json.Marshal(message.ToolInput)
				content = string(bytes)
			}
			out = append(out, anthropicMessage{
				Role: "user",
				Content: []anthropicContentBlock{{
					Type:      "tool_result",
					ToolUseID: message.ToolCallID,
					Content:   content,
					IsError:   message.IsError,
				}},
			})
		}
	}
	return out
}
