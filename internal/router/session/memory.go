package session

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"grubzo/internal/utils/random"
)

type memorySession struct {
	t         string
	refID     uuid.UUID
	userID    uint
	createdAt time.Time
	data      map[string]interface{}
	sync.Mutex
}

func newMemorySession(t string, refID uuid.UUID, userID uint, createdAt time.Time, data map[string]interface{}) *memorySession {
	return &memorySession{
		t:         t,
		refID:     refID,
		userID:    userID,
		createdAt: createdAt,
		data:      data,
	}
}

func (s *memorySession) Token() string        { return s.t }
func (s *memorySession) RefID() uuid.UUID     { return s.refID }
func (s *memorySession) UserID() uint         { return s.userID }
func (s *memorySession) CreatedAt() time.Time { return s.createdAt }
func (s *memorySession) LoggedIn() bool       { return s.userID != 0 }
func (s *memorySession) Expired() bool {
	return time.Since(s.createdAt) > time.Duration(sessionMaxAge)*time.Second
}
func (s *memorySession) Refreshable() bool {
	return time.Since(s.createdAt) <= time.Duration(sessionMaxAge+sessionKeepAge)*time.Second
}
func (s *memorySession) Get(k string) (interface{}, error) {
	s.Lock()
	defer s.Unlock()
	return s.data[k], nil
}
func (s *memorySession) Set(k string, v interface{}) error {
	s.Lock()
	defer s.Unlock()
	s.data[k] = v
	return nil
}
func (s *memorySession) Delete(k string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.data, k)
	return nil
}

type memoryStore struct {
	sessions map[string]*memorySession
	sync.RWMutex
}

func NewMemorySessionStore() Store {
	return &memoryStore{
		sessions: map[string]*memorySession{},
	}
}

func (ms *memoryStore) GetSession(c *gin.Context) (Session, error) {
	token, err := c.Cookie(CookieName)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	s, err := ms.GetSessionByToken(token)
	if err != nil {
		if err != ErrSessionNotFound {
			return nil, err
		}
	}

	if s != nil {
		if !s.Expired() {
			return s, nil
		}
		if s.Refreshable() {
			return ms.RenewSession(c, s.UserID())
		}
	}

	_ = ms.RevokeSession(c)
	return nil, ErrSessionNotFound
}

func (ms *memoryStore) GetSessionByToken(token string) (Session, error) {
	if len(token) == 0 {
		return nil, ErrSessionNotFound
	}

	ms.RLock()
	defer ms.RUnlock()
	s, ok := ms.sessions[token]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return s, nil
}

func (ms *memoryStore) GetSessionsByUserID(userID uint) ([]Session, error) {
	if userID == 0 {
		return []Session{}, nil
	}

	ms.RLock()
	defer ms.RUnlock()

	result := make([]Session, 0)
	for _, s := range ms.sessions {
		if s.UserID() == userID && s.Refreshable() {
			result = append(result, s)
		}
	}
	return result, nil
}

func (ms *memoryStore) RevokeSession(c *gin.Context) error {
	token, err := c.Cookie(CookieName)
	if err != nil || token == "" {
		return nil
	}

	ms.Lock()
	delete(ms.sessions, token)
	ms.Unlock()

	c.SetCookie(CookieName, "", -1, "/", "", false, true)
	return nil
}

func (ms *memoryStore) RevokeSessionByRefID(refID uuid.UUID) error {
	if refID == uuid.Nil {
		return nil
	}
	ms.Lock()
	defer ms.Unlock()
	for k, s := range ms.sessions {
		if s.RefID() == refID {
			delete(ms.sessions, k)
			return nil
		}
	}
	return nil
}

func (ms *memoryStore) RevokeSessionsByUserID(userID uint) error {
	if userID == 0 {
		return nil
	}
	ms.Lock()
	defer ms.Unlock()
	for k, s := range ms.sessions {
		if s.UserID() == userID {
			delete(ms.sessions, k)
		}
	}
	return nil
}

func (ms *memoryStore) RenewSession(c *gin.Context, userID uint) (Session, error) {
	oldToken, _ := c.Cookie(CookieName)
	if len(oldToken) > 0 {
		ms.Lock()
		delete(ms.sessions, oldToken)
		ms.Unlock()
	}

	s, err := ms.IssueSession(userID, nil)
	if err != nil {
		return nil, err
	}

	c.SetCookie(
		CookieName,
		s.Token(),
		sessionMaxAge+sessionKeepAge,
		"/",
		"",
		false,
		true,
	)
	return s, nil
}

func (ms *memoryStore) IssueSession(userID uint, data map[string]interface{}) (Session, error) {
	if data == nil {
		data = map[string]interface{}{}
	}

	s := newMemorySession(
		random.SecureAlphaNumeric(50),
		uuid.Must(uuid.NewV7()),
		userID,
		time.Now(),
		data,
	)

	ms.Lock()
	ms.sessions[s.Token()] = s
	ms.Unlock()

	return s, nil
}
