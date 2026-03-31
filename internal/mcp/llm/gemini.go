package llm

import (
	"context"
	"fmt"
	"strings"
)

type geminiProvider struct {
	httpProvider
}

type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"system_instruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	Tools             []geminiTool    `json:"tools,omitempty"`
	GenerationConfig  *geminiConfig   `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text             string                  `json:"text,omitempty"`
	FunctionCall     *geminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *geminiFunctionResponse `json:"functionResponse,omitempty"`
	ThoughtSignature string                  `json:"thoughtSignature,omitempty"`
}

type geminiFunctionCall struct {
	ID   string         `json:"id,omitempty"`
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

type geminiFunctionResponse struct {
	ID       string         `json:"id,omitempty"`
	Name     string         `json:"name"`
	Response map[string]any `json:"response,omitempty"`
}

type geminiTool struct {
	FunctionDeclarations []geminiFunctionDeclaration `json:"functionDeclarations"`
}

type geminiFunctionDeclaration struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type geminiConfig struct {
	Temperature     *float64 `json:"temperature,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
}

func (p *geminiProvider) Generate(ctx context.Context, request Request) (*Response, error) {
	model := strings.TrimSpace(request.Model)
	if model == "" {
		model = strings.TrimSpace(p.model)
	}
	if model == "" {
		return nil, fmt.Errorf("gemini model is not configured")
	}

	payload := geminiRequest{
		Contents: geminiContentsFromHistory(request.Messages),
	}
	if prompt := strings.TrimSpace(request.SystemPrompt); prompt != "" {
		payload.SystemInstruction = &geminiContent{
			Parts: []geminiPart{{Text: prompt}},
		}
	}
	if request.Temperature > 0 || request.MaxTokens > 0 {
		payload.GenerationConfig = &geminiConfig{}
		if request.Temperature > 0 {
			payload.GenerationConfig.Temperature = &request.Temperature
		}
		if request.MaxTokens > 0 {
			payload.GenerationConfig.MaxOutputTokens = request.MaxTokens
		}
	}
	if len(request.Tools) > 0 {
		declarations := make([]geminiFunctionDeclaration, 0, len(request.Tools))
		for _, tool := range request.Tools {
			declarations = append(declarations, geminiFunctionDeclaration{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			})
		}
		payload.Tools = []geminiTool{{
			FunctionDeclarations: declarations,
		}}
	}

	var resp geminiResponse
	if err := p.postJSON(
		ctx,
		p.baseURL(model),
		map[string]string{"x-goog-api-key": p.apiKey},
		payload,
		&resp,
	); err != nil {
		return nil, err
	}
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("gemini returned no candidates")
	}

	candidate := resp.Candidates[0]
	out := &Response{StopReason: candidate.FinishReason}
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			out.Text += part.Text
		}
		if part.FunctionCall != nil {
			meta := map[string]any{}
			if part.ThoughtSignature != "" {
				meta["thoughtSignature"] = part.ThoughtSignature
			}
			out.ToolCalls = append(out.ToolCalls, ToolCall{
				ID:        part.FunctionCall.ID,
				Name:      part.FunctionCall.Name,
				Arguments: part.FunctionCall.Args,
				Meta:      meta,
			})
		}
	}
	return out, nil
}

func (p *geminiProvider) baseURL(model string) string {
	base := strings.TrimSpace(p.base)
	if base == "" {
		base = "https://generativelanguage.googleapis.com/v1beta"
	}
	if strings.Contains(base, ":generateContent") {
		return base
	}
	if strings.HasSuffix(base, "/models") {
		return fmt.Sprintf("%s/%s:generateContent", base, model)
	}
	return fmt.Sprintf("%s/models/%s:generateContent", strings.TrimRight(base, "/"), model)
}

func geminiContentsFromHistory(history []Message) []geminiContent {
	out := make([]geminiContent, 0, len(history))
	for index := 0; index < len(history); {
		message := history[index]
		switch {
		case message.Role == "user":
			out = append(out, geminiContent{
				Role: "user",
				Parts: []geminiPart{{
					Text: message.Text,
				}},
			})
			index++
		case message.Role == "assistant" && message.ToolName == "":
			part := geminiPart{
				Text: message.Text,
			}
			if signature := geminiThoughtSignature(message.Meta); signature != "" {
				part.ThoughtSignature = signature
			}
			out = append(out, geminiContent{
				Role:  "model",
				Parts: []geminiPart{part},
			})
			index++
		case message.Role == "assistant" && message.ToolName != "":
			parts := make([]geminiPart, 0, 1)
			for index < len(history) {
				current := history[index]
				if current.Role != "assistant" || current.ToolName == "" {
					break
				}
				part := geminiPart{
					FunctionCall: &geminiFunctionCall{
						ID:   current.ToolCallID,
						Name: current.ToolName,
						Args: current.ToolInput,
					},
				}
				if signature := geminiThoughtSignature(current.Meta); signature != "" {
					part.ThoughtSignature = signature
				} else if len(parts) == 0 {
					part.ThoughtSignature = "skip_thought_signature_validator"
				}
				parts = append(parts, part)
				index++
			}
			out = append(out, geminiContent{
				Role:  "model",
				Parts: parts,
			})
		case message.Role == "tool":
			parts := make([]geminiPart, 0, 1)
			for index < len(history) {
				current := history[index]
				if current.Role != "tool" {
					break
				}
				responsePayload := map[string]any{
					"result": current.ToolInput,
					"text":   current.Text,
				}
				if current.IsError {
					responsePayload["error"] = true
				}
				parts = append(parts, geminiPart{
					FunctionResponse: &geminiFunctionResponse{
						ID:       current.ToolCallID,
						Name:     current.ToolName,
						Response: responsePayload,
					},
				})
				index++
			}
			out = append(out, geminiContent{
				Role:  "user",
				Parts: parts,
			})
		}
	}
	return out
}

func geminiThoughtSignature(meta map[string]any) string {
	if len(meta) == 0 {
		return ""
	}
	if value, ok := meta["thoughtSignature"].(string); ok {
		return value
	}
	if value, ok := meta["thought_signature"].(string); ok {
		return value
	}
	if provider, ok := meta["provider"].(map[string]any); ok {
		if value, ok := provider["thoughtSignature"].(string); ok {
			return value
		}
		if value, ok := provider["thought_signature"].(string); ok {
			return value
		}
	}
	return ""
}
