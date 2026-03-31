package llm

import (
	"context"
	"encoding/json"
	"fmt"
)

type openAIProvider struct {
	httpProvider
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Tools       []openAITool    `json:"tools,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type openAITool struct {
	Type     string             `json:"type"`
	Function openAIFunctionTool `json:"function"`
}

type openAIFunctionTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type openAIMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
	Name       string           `json:"name,omitempty"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
}

type openAIToolCall struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Function openAIToolCallFunction `json:"function"`
}

type openAIToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}

func (p *openAIProvider) Generate(ctx context.Context, request Request) (*Response, error) {
	model := request.Model
	if model == "" {
		model = p.model
	}
	if model == "" {
		return nil, fmt.Errorf("openai-compatible model is not configured")
	}

	payload := openAIRequest{
		Model:       model,
		Messages:    openAIMessagesFromHistory(request.SystemPrompt, request.Messages),
		Temperature: request.Temperature,
		MaxTokens:   request.MaxTokens,
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

	var resp openAIResponse
	if err := p.postJSON(
		ctx,
		p.baseURL(),
		map[string]string{"Authorization": "Bearer " + p.apiKey},
		payload,
		&resp,
	); err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai-compatible provider returned no choices")
	}

	choice := resp.Choices[0].Message
	out := &Response{Text: choice.Content}
	for _, toolCall := range choice.ToolCalls {
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

func (p *openAIProvider) baseURL() string {
	if p.base == "" {
		return "https://api.openai.com/v1/chat/completions"
	}
	if p.base == "https://api.openai.com/v1" {
		return p.base + "/chat/completions"
	}
	return p.base
}

func openAIMessagesFromHistory(systemPrompt string, history []Message) []openAIMessage {
	out := make([]openAIMessage, 0, len(history)+1)
	if systemPrompt != "" {
		out = append(out, openAIMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	for _, message := range history {
		switch {
		case message.Role == "user":
			out = append(out, openAIMessage{
				Role:    "user",
				Content: message.Text,
			})
		case message.Role == "assistant" && message.ToolName == "":
			out = append(out, openAIMessage{
				Role:    "assistant",
				Content: message.Text,
			})
		case message.Role == "assistant" && message.ToolName != "":
			args, _ := json.Marshal(message.ToolInput)
			out = append(out, openAIMessage{
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
			out = append(out, openAIMessage{
				Role:       "tool",
				Content:    message.Text,
				ToolCallID: message.ToolCallID,
				Name:       message.ToolName,
			})
		}
	}
	return out
}
