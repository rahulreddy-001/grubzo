package session

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"grubzo/internal/utils"
	"grubzo/internal/utils/random"
)

type memorySession struct {
	t         string
	userID    uint64
	tenantID  uint64
	createdAt time.Time
	data      map[string]interface{}
	sync.Mutex
}

func newMemorySession(t string, userID, tenantID uint64, createdAt time.Time, data map[string]interface{}) *memorySession {
	return &memorySession{
		t:         t,
		userID:    userID,
		tenantID:  tenantID,
		createdAt: createdAt,
		data:      data,
	}
}

func (s *memorySession) Token() string        { return s.t }
func (s *memorySession) UserID() uint64       { return s.userID }
func (s *memorySession) TenantID() uint64     { return s.tenantID }
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

func (s *memorySession) SetUserSession(us UserSession) error {
	usRaw, err := json.Marshal(us)
	if err != nil {
		return err
	}
	return s.Set(UserSessionKey, usRaw)
}

func (s *memorySession) GetUserSession() (UserSession, error) {
	if val, err := s.Get(UserSessionKey); err != nil {
		return UserSession{}, err
	} else {
		return utils.AsType[UserSession](val)
	}
}

type memoryStore struct {
	sessions     map[string]*memorySession
	cookieDomain string
	sync.RWMutex
}

func NewMemorySessionStore(cookieDomain ...string) Store {
	domain := ""
	if len(cookieDomain) > 0 {
		domain = normalizeCookieDomain(cookieDomain[0])
	}
	return &memoryStore{
		sessions:     map[string]*memorySession{},
		cookieDomain: domain,
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
			return ms.RenewSession(c, s.UserID(), s.TenantID())
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

func (ms *memoryStore) GetSessionsByUserID(userID, tenantID uint64) ([]Session, error) {
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

	c.SetCookie(CookieName, "", -1, "/", ms.cookieDomain, false, true)
	return nil
}

func (ms *memoryStore) RenewSession(c *gin.Context, userID, tenantID uint64) (Session, error) {
	oldToken, _ := c.Cookie(CookieName)
	if len(oldToken) > 0 {
		ms.Lock()
		delete(ms.sessions, oldToken)
		ms.Unlock()
	}

	s, err := ms.IssueSession(userID, tenantID, nil)
	if err != nil {
		return nil, err
	}

	c.SetCookie(
		CookieName,
		s.Token(),
		sessionMaxAge+sessionKeepAge,
		"/",
		ms.cookieDomain,
		false,
		true,
	)
	return s, nil
}

func (ms *memoryStore) IssueSession(userID, tenantID uint64, data map[string]interface{}) (Session, error) {
	if data == nil {
		data = map[string]interface{}{}
	}

	s := newMemorySession(
		random.SecureAlphaNumeric(50),
		userID,
		tenantID,
		time.Now(),
		data,
	)

	ms.Lock()
	ms.sessions[s.Token()] = s
	ms.Unlock()

	return s, nil
}
