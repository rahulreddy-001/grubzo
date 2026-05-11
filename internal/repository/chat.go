package repository

//go:generate go run ../../cmd/injecttrace -file chat.go -receiver Repository -service Repository

import (
	"context"
	"encoding/json"
	"errors"
	"grubzo/internal/models/entity"
	"time"

	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type ChatRepository interface {
	EnsureChatSession(ctx context.Context, input *ChatSessionInput) (*entity.ChatSession, error)
	AppendChatMessage(ctx context.Context, sessionID string, input StoredChatMessage) error
	ListChatMessages(ctx context.Context, sessionID string, userID, tenantID uint64) ([]entity.ChatMessage, error)
	ListChatSessions(ctx context.Context, userID, tenantID uint64) ([]ChatSessionSummary, error)
	DeleteChatSession(ctx context.Context, sessionID string, userID, tenantID uint64) error
}

type ChatSessionInput struct {
	SessionID  string
	UserID     uint64
	TenantID   uint64
	LocationID uint64
	Provider   string
	Model      string
}

type ChatMessageRole string

const (
	ChatMessageRoleUser      ChatMessageRole = "user"
	ChatMessageRoleAssistant ChatMessageRole = "assistant"
	ChatMessageRoleTool      ChatMessageRole = "tool"
)

type ChatMessageKind string

const (
	ChatMessageKindText       ChatMessageKind = "text"
	ChatMessageKindToolCall   ChatMessageKind = "tool_call"
	ChatMessageKindToolResult ChatMessageKind = "tool_result"
)

type StoredChatMessage struct {
	Role       ChatMessageRole
	Kind       ChatMessageKind
	Content    string
	ToolName   string
	ToolCallID string
	Meta       any
	IsError    bool
}

type ChatSessionSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Provider    string `json:"provider"`
	Model       string `json:"model"`
	LastMessage string `json:"last_message"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (r *Repository) EnsureChatSession(ctx context.Context, input *ChatSessionInput) (*entity.ChatSession, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.EnsureChatSession")
	defer span.End()

	sessionID := input.SessionID
	if sessionID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		sessionID = id.String()
	}

	var session entity.ChatSession
	err := r.dbWithContext(ctx).
		Where("id = ? AND user_id = ? AND tenant_id = ?", sessionID, input.UserID, input.TenantID).
		First(&session).Error
	if err == nil {
		session.LocationID = input.LocationID
		session.Provider = input.Provider
		session.Model = input.Model
		if saveErr := r.dbWithContext(ctx).Save(&session).Error; saveErr != nil {
			return nil, saveErr
		}
		return &session, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	session = entity.ChatSession{
		ID:         sessionID,
		UserID:     input.UserID,
		TenantID:   input.TenantID,
		LocationID: input.LocationID,
		Provider:   input.Provider,
		Model:      input.Model,
	}
	if err := r.dbWithContext(ctx).Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) AppendChatMessage(ctx context.Context, sessionID string, input StoredChatMessage) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.AppendChatMessage")
	defer span.End()

	metaJSON := "{}"
	if input.Meta != nil {
		if payload, err := json.Marshal(input.Meta); err == nil {
			metaJSON = string(payload)
		}
	}

	message := entity.ChatMessage{
		SessionID:  sessionID,
		Role:       string(input.Role),
		Kind:       string(input.Kind),
		Content:    input.Content,
		ToolName:   input.ToolName,
		ToolCallID: input.ToolCallID,
		MetaJSON:   metaJSON,
		IsError:    input.IsError,
	}
	if err := r.dbWithContext(ctx).Create(&message).Error; err != nil {
		return err
	}
	return r.dbWithContext(ctx).
		Model(&entity.ChatSession{}).
		Where("id = ?", sessionID).
		Update("updated_at", time.Now().UTC()).
		Error
}

func (r *Repository) ListChatMessages(ctx context.Context, sessionID string, userID, tenantID uint64) ([]entity.ChatMessage, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.ListChatMessages")
	defer span.End()

	var session entity.ChatSession
	if err := r.dbWithContext(ctx).
		Where("id = ? AND user_id = ? AND tenant_id = ?", sessionID, userID, tenantID).
		First(&session).Error; err != nil {
		return nil, err
	}

	var messages []entity.ChatMessage
	err := r.dbWithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("id asc").
		Find(&messages).Error
	return messages, err
}

func (r *Repository) ListChatSessions(ctx context.Context, userID, tenantID uint64) ([]ChatSessionSummary, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.ListChatSessions")
	defer span.End()

	var sessions []entity.ChatSession
	if err := r.dbWithContext(ctx).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Order("updated_at desc").
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	summaries := make([]ChatSessionSummary, 0, len(sessions))
	for _, session := range sessions {
		var latest entity.ChatMessage
		var firstUserMessage entity.ChatMessage
		preview := ""
		title := "New chat"

		err := r.dbWithContext(ctx).
			Where("session_id = ?", session.ID).
			Order("id desc").
			First(&latest).Error
		if err == nil {
			preview = latest.Content
			if len(preview) > 512 {
				preview = preview[:512]
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		err = r.dbWithContext(ctx).
			Where(
				"session_id = ? AND role = ? AND kind = ?",
				session.ID,
				string(ChatMessageRoleUser),
				string(ChatMessageKindText),
			).
			Order("id asc").
			First(&firstUserMessage).Error
		if err == nil {
			title = firstUserMessage.Content
			if len(title) > 120 {
				title = title[:120]
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		summaries = append(summaries, ChatSessionSummary{
			ID:          session.ID,
			Title:       title,
			Provider:    session.Provider,
			Model:       session.Model,
			LastMessage: preview,
			CreatedAt:   session.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   session.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	return summaries, nil
}

func (r *Repository) DeleteChatSession(ctx context.Context, sessionID string, userID, tenantID uint64) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.DeleteChatSession")
	defer span.End()

	return r.dbWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var session entity.ChatSession
		if err := tx.
			Where("id = ? AND user_id = ? AND tenant_id = ?", sessionID, userID, tenantID).
			First(&session).Error; err != nil {
			return err
		}

		if err := tx.Where("session_id = ?", sessionID).Delete(&entity.ChatMessage{}).Error; err != nil {
			return err
		}

		return tx.Delete(&session).Error
	})
}
