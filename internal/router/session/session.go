package session

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CookieName     = "SESSION_TOKEN"
	sessionMaxAge  = 60 * 60 * 24 * 14
	sessionKeepAge = 60 * 60 * 24 * 14
)

const (
	UserSessionKey = "user"
)

type UserSession struct {
	UserID   uint64   `json:"user_id"`
	TenantID uint64   `json:"tenant_id"`
	Email    string   `json:"email"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
	Location uint64   `json:"location"`
}

var (
	ErrNoSessionStore     = errors.New("no session store found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrNoUserSession      = errors.New("no active user session found")
	ErrInvalidUserSession = errors.New("invalid user session")
)

type Session interface {
	Token() string
	UserID() uint64
	TenantID() uint64
	CreatedAt() time.Time
	LoggedIn() bool
	Expired() bool
	Refreshable() bool

	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
	SetUserSession(UserSession) error
	GetUserSession() (UserSession, error)
}

type Store interface {
	IssueSession(userID, tenantID uint64, data map[string]interface{}) (Session, error)

	GetSession(c *gin.Context) (Session, error)
	RenewSession(c *gin.Context, userID, tenantID uint64) (Session, error)
	RevokeSession(c *gin.Context) error

	GetSessionByToken(token string) (Session, error)
	GetSessionsByUserID(userID, tenantID uint64) ([]Session, error)
}
