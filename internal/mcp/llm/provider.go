package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"grubzo/internal/config"
	mcptrace "grubzo/internal/mcp/tracing"
	"io"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type httpProvider struct {
	name   string
	model  string
	client *http.Client
	apiKey string
	base   string
}

type unavailableProvider struct {
	reason string
}

func NewProvider(cfg *config.Config) (Provider, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.MCP.LLM.Provider))
	timeout := time.Duration(cfg.MCP.LLM.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	client := &http.Client{Timeout: timeout}

	switch provider {
	case "":
		return unavailableProvider{reason: "mcp.llm.provider is not configured"}, nil
	case "anthropic":
		return &anthropicProvider{
			httpProvider: httpProvider{
				name:   "anthropic",
				model:  cfg.MCP.LLM.Model,
				client: client,
				apiKey: cfg.MCP.LLM.APIKey,
				base:   strings.TrimRight(cfg.MCP.LLM.BaseURL, "/"),
			},
		}, nil
	case "openai", "openai-compatible", "openai_compatible":
		return &openAIProvider{
			httpProvider: httpProvider{
				name:   "openai-compatible",
				model:  cfg.MCP.LLM.Model,
				client: client,
				apiKey: cfg.MCP.LLM.APIKey,
				base:   strings.TrimRight(cfg.MCP.LLM.BaseURL, "/"),
			},
		}, nil
	case "ollama":
		return &ollamaProvider{
			httpProvider: httpProvider{
				name:   "ollama",
				model:  cfg.MCP.LLM.Model,
				client: client,
				apiKey: cfg.MCP.LLM.APIKey,
				base:   strings.TrimRight(cfg.MCP.LLM.BaseURL, "/"),
			},
		}, nil
	case "gemini", "google", "google-gemini":
		return &geminiProvider{
			httpProvider: httpProvider{
				name:   "gemini",
				model:  cfg.MCP.LLM.Model,
				client: client,
				apiKey: cfg.MCP.LLM.APIKey,
				base:   strings.TrimRight(cfg.MCP.LLM.BaseURL, "/"),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported mcp.llm.provider %q", cfg.MCP.LLM.Provider)
	}
}

func (p unavailableProvider) Name() string {
	return "unconfigured"
}

func (p unavailableProvider) DefaultModel() string {
	return ""
}

func (p unavailableProvider) Generate(_ context.Context, _ Request) (*Response, error) {
	return nil, fmt.Errorf("llm provider unavailable: %s", p.reason)
}

func (p *httpProvider) Name() string {
	return p.name
}

func (p *httpProvider) DefaultModel() string {
	return p.model
}

func (p *httpProvider) postJSON(
	ctx context.Context,
	url string,
	headers map[string]string,
	payload any,
	out any,
) (err error) {
	ctx, span := mcptrace.StartClient(
		ctx,
		"MCPLLM.postJSON",
		http.MethodPost,
		url,
		payload,
		attribute.String("llm.provider", p.name),
	)
	defer mcptrace.End(span, &err)

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	mcptrace.InjectHTTPHeaders(ctx, req.Header)
	for key, value := range headers {
		if value != "" {
			req.Header.Set(key, value)
		}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	span.SetAttributes(
		semconv.HTTPResponseStatusCode(resp.StatusCode),
		semconv.HTTPResponseBodySize(len(respBody)),
	)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s returned %d: %s", p.name, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(respBody, out)
}
