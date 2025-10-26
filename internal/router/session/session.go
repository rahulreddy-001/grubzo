package session

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

const (
	CookieName     = "SESSION_TOKEN"
	sessionMaxAge  = 60 * 60 * 24 * 14
	sessionKeepAge = 60 * 60 * 24 * 14
	cacheSize      = 2048
)

var ErrSessionNotFound = errors.New("session not found")

type Session interface {
	Token() string
	RefID() uuid.UUID
	UserID() uint
	CreatedAt() time.Time
	LoggedIn() bool

	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error

	Expired() bool
	Refreshable() bool
}

type Store interface {
	GetSession(c *gin.Context) (Session, error)
	GetSessionByToken(token string) (Session, error)
	GetSessionsByUserID(userID uint) ([]Session, error)
	RevokeSession(c *gin.Context) error
	RevokeSessionByRefID(refID uuid.UUID) error
	RevokeSessionsByUserID(userID uint) error
	RenewSession(c *gin.Context, userID uint) (Session, error)
	IssueSession(userID uint, data map[string]interface{}) (Session, error)
}
