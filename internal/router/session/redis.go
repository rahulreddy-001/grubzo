package session

import (
	"context"
	"encoding/json"
	"fmt"
	"grubzo/internal/utils"
	"grubzo/internal/utils/random"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type sessionRecord struct {
	Token    string    `json:"token"`
	UserID   uint      `json:"user_id"`
	TenantID uint      `json:"tenant_id"`
	Data     []byte    `json:"data"`
	Created  time.Time `json:"created_at"`
}

func (r *sessionRecord) toMap() (map[string]any, error) {
	data := map[string]any{}
	err := json.Unmarshal(r.Data, &data)
	return data, err
}

type redisSession struct {
	t         string
	userID    uint
	tenantID  uint
	createdAt time.Time
	data      map[string]any
	db        *redis.Client
	ctx       context.Context
	sync.Mutex
}

func createRedisSession(ctx context.Context, token string, userID uint, tenantID uint, data map[string]any, db *redis.Client) *redisSession {
	if data == nil {
		data = map[string]any{}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &redisSession{
		t:         token,
		userID:    userID,
		tenantID:  tenantID,
		createdAt: time.Now(),
		data:      data,
		db:        db,
		ctx:       ctx,
	}
}

func (rs *redisSession) sync() error {
	rawData, err := json.Marshal(rs.data)
	if err != nil {
		return err
	}
	record := sessionRecord{
		Token:    rs.t,
		UserID:   rs.userID,
		TenantID: rs.tenantID,
		Data:     rawData,
		Created:  rs.createdAt,
	}
	err = rs.db.JSONSet(rs.ctx, record.Token, ".", record).Err()
	if err != nil {
		return err
	}
	ttl := time.Duration(sessionMaxAge+sessionKeepAge)*time.Second - time.Since(rs.CreatedAt())
	return rs.db.Expire(rs.ctx, record.Token, ttl).Err()
}

func (rs *redisSession) Token() string        { return rs.t }
func (rs *redisSession) UserID() uint         { return rs.userID }
func (rs *redisSession) TenantID() uint       { return rs.tenantID }
func (rs *redisSession) CreatedAt() time.Time { return rs.createdAt }
func (rs *redisSession) LoggedIn() bool       { return rs.userID != 0 }

func (rs *redisSession) Expired() bool {
	return time.Since(rs.CreatedAt()) > time.Duration(sessionMaxAge)*time.Second
}

func (rs *redisSession) Refreshable() bool {
	return time.Since(rs.createdAt) <= time.Duration(sessionMaxAge+sessionKeepAge)*time.Second
}

func (rs *redisSession) Get(key string) (interface{}, error) {
	rs.Lock()
	defer rs.Unlock()
	return rs.data[key], nil
}

func (rs *redisSession) Set(key string, value interface{}) error {
	rs.Lock()
	defer rs.Unlock()
	rs.data[key] = value
	return rs.sync()
}

func (rs *redisSession) Delete(key string) error {
	rs.Lock()
	defer rs.Unlock()
	delete(rs.data, key)
	return rs.sync()
}
func (rs *redisSession) SetUserSession(us UserSession) error {
	return rs.Set(UserSessionKey, us)
}

func (rs *redisSession) GetUserSession() (UserSession, error) {
	raw, err := rs.Get(UserSessionKey)
	if err != nil {
		return UserSession{}, err
	}
	return utils.AsType[UserSession](raw)

}

type redisStore struct {
	db *redis.Client
}

func NewRedisSessionStore(db *redis.Client) Store {
	return &redisStore{
		db: db,
	}
}

func (rs *redisStore) IssueSession(userID uint, tenantID uint, data map[string]any) (Session, error) {
	return rs.issueSession(context.Background(), userID, tenantID, data)
}

func (rs *redisStore) issueSession(ctx context.Context, userID uint, tenantID uint, data map[string]any) (Session, error) {
	token := fmt.Sprintf("%s:%s", "session", random.SecureAlphaNumeric(50))
	session := createRedisSession(ctx, token, userID, tenantID, data, rs.db)
	if err := session.sync(); err != nil {
		return nil, err
	}
	return session, nil
}

func (rs *redisStore) GetSession(c *gin.Context) (Session, error) {
	token, err := c.Cookie(CookieName)
	if err != nil {
		return nil, ErrNoUserSession
	}
	session, err := rs.getSessionByToken(c.Request.Context(), token)
	if err != nil {
		return nil, err
	}
	if session != nil {
		if !session.Expired() {
			return session, nil
		}
		if session.Refreshable() {
			return rs.RenewSession(c, session.UserID(), session.TenantID())
		}
	}
	_ = rs.RevokeSession(c)
	return nil, ErrSessionNotFound
}

func (rs *redisStore) RevokeSession(c *gin.Context) error {
	token, err := c.Cookie(CookieName)
	if err != nil {
		return ErrNoUserSession
	}
	if err := rs.db.JSONDel(c.Request.Context(), token, ".").Err(); err != nil {
		return err
	}
	c.SetCookie(CookieName, "", -1, "/", "", false, true)
	return nil
}

func (rs *redisStore) GetSessionsByUserID(userID, tenantID uint) ([]Session, error) {
	//TODO: create FTSearch index `idx_session_user_id_tenant_id`
	qs := fmt.Sprintf("@user_id[%d %d] @tenant_id[%d %d] ", userID, userID, tenantID, tenantID)
	searchResult, err := rs.db.FTSearch(context.Background(), "idx_session_user_id_tenant_id", qs).Result()
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, 0, len(searchResult.Docs))
	for _, doc := range searchResult.Docs {
		var rec sessionRecord
		if err := json.Unmarshal([]byte(doc.Fields["$"]), &rec); err != nil {
			return nil, err
		}
		data, err := rec.toMap()
		if err != nil {
			return nil, err
		}
		session := createRedisSession(
			context.Background(),
			rec.Token,
			rec.UserID,
			rec.TenantID,
			data,
			rs.db,
		)
		if session.Refreshable() {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (rs *redisStore) RenewSession(c *gin.Context, userID, tenantID uint) (Session, error) {
	rs.RevokeSession(c)
	newSession, err := rs.issueSession(c.Request.Context(), userID, tenantID, nil)
	if err != nil {
		return nil, err
	}

	c.SetCookie(
		CookieName,
		newSession.Token(),
		sessionMaxAge+sessionKeepAge,
		"/",
		"",
		false,
		true,
	)
	return newSession, nil
}

func (rs *redisStore) GetSessionByToken(token string) (Session, error) {
	return rs.getSessionByToken(context.Background(), token)
}

func (rs *redisStore) getSessionByToken(ctx context.Context, token string) (Session, error) {
	res, err := rs.db.JSONGet(ctx, token, ".").Result()
	if err != nil {
		return nil, err
	}

	record := sessionRecord{}
	if err := json.Unmarshal([]byte(res), &record); err != nil {
		return nil, err
	}
	data, err := record.toMap()
	if err != nil {
		return nil, err
	}
	return createRedisSession(
		ctx,
		record.Token,
		record.UserID,
		record.TenantID,
		data,
		rs.db,
	), nil
}
