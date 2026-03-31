package mcp

import (
	"encoding/json"
	"fmt"
	"grubzo/internal/mcp/protocol"
	mcptrace "grubzo/internal/mcp/tracing"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func (h *Handlers) Describe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "grubzo-mcp",
		"description": "Grubzo MCP routes for menu discovery, order tracking, and agent chat.",
		"tools_url":   "/mcp/tools",
		"rpc_url":     "/mcp",
	})
}

func (h *Handlers) ListTools(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tools": h.Tools.ToolDefinitions()})
}

func (h *Handlers) InvokeTool(c *gin.Context) {
	ctx, span := mcptrace.Start(
		c.Request.Context(),
		"MCPRouter.InvokeTool",
		nil,
		attribute.String("mcp.tool.name", c.Param("tool")),
	)
	defer span.End()

	originalRequest := c.Request
	c.Request = c.Request.WithContext(ctx)
	defer func() {
		c.Request = originalRequest
	}()

	currentActor, err := actorFromGin(c)
	if err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	span.SetAttributes(
		attribute.Int("mcp.payload.size", len(body)),
		attribute.String("mcp.payload.preview", truncatePayloadPreview(body)),
	)

	result, err := h.Tools.Call(ctx, currentActor, c.Param("tool"), body)
	if err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HandleJSONRPC(c *gin.Context) {
	ctx, span := mcptrace.Start(
		c.Request.Context(),
		"MCPRouter.HandleJSONRPC",
		nil,
		semconv.RPCSystemKey.String("jsonrpc"),
		semconv.RPCService("grubzo-mcp"),
	)
	defer span.End()

	originalRequest := c.Request
	c.Request = c.Request.WithContext(ctx)
	defer func() {
		c.Request = originalRequest
	}()

	var request protocol.JSONRPCRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusBadRequest, protocol.Failure(nil, -32700, "invalid JSON-RPC payload", err.Error()))
		return
	}
	span.SetAttributes(
		semconv.RPCMethod(request.Method),
		semconv.RPCJsonrpcVersion(request.JSONRPC),
	)
	if request.ID != nil {
		span.SetAttributes(semconv.RPCJsonrpcRequestID(fmt.Sprint(request.ID)))
	}

	if request.JSONRPC != "2.0" {
		c.JSON(http.StatusBadRequest, protocol.Failure(request.ID, -32600, "jsonrpc must be 2.0", nil))
		return
	}

	currentActor, err := actorFromGin(c)
	if err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusUnauthorized, protocol.Failure(request.ID, -32001, err.Error(), nil))
		return
	}

	switch request.Method {
	case "initialize":
		c.JSON(http.StatusOK, protocol.Success(request.ID, protocol.InitializeResult{
			ProtocolVersion: "2025-11-25",
			ServerInfo: map[string]any{
				"name":    "grubzo-mcp",
				"version": "0.1.0",
			},
			Capabilities: map[string]any{
				"tools": map[string]any{},
			},
		}))
	case "tools/list":
		c.JSON(http.StatusOK, protocol.Success(request.ID, gin.H{
			"tools": h.Tools.ToolDefinitions(),
		}))
	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(request.Params, &params); err != nil {
			mcptrace.RecordError(span, err)
			c.JSON(http.StatusBadRequest, protocol.Failure(request.ID, -32602, "invalid tools/call params", err.Error()))
			return
		}
		result, err := h.Tools.Call(ctx, currentActor, params.Name, params.Arguments)
		if err != nil {
			mcptrace.RecordError(span, err)
			c.JSON(http.StatusBadRequest, protocol.Failure(request.ID, -32602, err.Error(), nil))
			return
		}
		c.JSON(http.StatusOK, protocol.Success(request.ID, result))
	default:
		c.JSON(http.StatusBadRequest, protocol.Failure(request.ID, -32601, "method not found", request.Method))
	}
}
