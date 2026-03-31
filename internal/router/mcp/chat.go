package mcp

import (
	"encoding/json"
	"fmt"
	"grubzo/internal/mcp/agent"
	mcptrace "grubzo/internal/mcp/tracing"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (h *Handlers) Chat(c *gin.Context) {
	ctx, span := mcptrace.Start(c.Request.Context(), "MCPRouter.Chat", nil)
	defer span.End()

	originalRequest := c.Request
	c.Request = c.Request.WithContext(ctx)
	defer func() {
		c.Request = originalRequest
	}()

	var request agent.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	span.SetAttributes(
		attribute.String("mcp.session.id", request.SessionID),
		attribute.Int("mcp.payload.size", len(request.Message)),
		attribute.String("mcp.payload.preview", truncatePayloadPreview(request.Message)),
	)

	currentActor, err := actorFromGin(c)
	if err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming is not supported by this server"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Status(http.StatusOK)

	emit := func(event agent.ChatEvent) error {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(c.Writer, "event: %s\n", event.Type); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", payload); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	response, err := h.Agent.Run(ctx, currentActor, request, emit)
	if err != nil {
		mcptrace.RecordError(span, err)
		h.Logger.Warn("chat request failed", zap.Error(err))
		_ = emit(agent.ChatEvent{Type: agent.ChatEventTypeError, Message: err.Error()})
		return
	}
	_ = emit(agent.ChatEvent{
		Type:      agent.ChatEventTypeDone,
		SessionID: response.SessionID,
		Payload:   response,
	})
}

func (h *Handlers) ListChatMessages(c *gin.Context) {
	ctx, span := mcptrace.Start(c.Request.Context(), "MCPRouter.ListChatMessages", nil)
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
	messages, err := h.Repository.ListChatMessages(ctx, c.Param("id"), currentActor.UserID, currentActor.TenantID)
	if err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *Handlers) ListChatSessions(c *gin.Context) {
	ctx, span := mcptrace.Start(c.Request.Context(), "MCPRouter.ListChatSessions", nil)
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
	sessions, err := h.Repository.ListChatSessions(ctx, currentActor.UserID, currentActor.TenantID)
	if err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *Handlers) DeleteChatSession(c *gin.Context) {
	ctx, span := mcptrace.StartServer(c.Request.Context(), "MCPRouter.DeleteChatSession", nil)
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

	if err := h.Repository.DeleteChatSession(ctx, c.Param("id"), currentActor.UserID, currentActor.TenantID); err != nil {
		mcptrace.RecordError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session deleted"})
}
