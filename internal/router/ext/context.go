package ext

import (
	"grubzo/internal/router/session"
	"grubzo/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	SessionKey      = "session"
	SessionStoreKey = "session_store"
)

type Context struct {
	ctx *gin.Context
}

func Ctx(ctx *gin.Context) *Context {
	return &Context{ctx: ctx}
}

func (c *Context) SetSession(sess session.Session) *Context {
	c.ctx.Set(SessionKey, sess)
	return c
}

func (c *Context) GetSession() (session.Session, error) {
	val, ok := c.ctx.Get(SessionKey)
	if !ok {
		return nil, session.ErrNoUserSession
	}
	return utils.AsType[session.Session](val)
}

func (c *Context) GetUserSession() (session.UserSession, error) {
	if sess, err := c.GetSession(); err != nil || sess == nil {
		return session.UserSession{}, session.ErrNoUserSession
	} else {
		return sess.GetUserSession()
	}
}

func (c *Context) IsLoggedIn() bool {
	sess, err := c.GetUserSession()
	if err != nil {
		return false
	}
	return sess.UserID != 0
}

func (c *Context) TenantID() uint {
	sess, err := c.GetUserSession()
	if err != nil {
		return 0
	}
	return sess.TenantID
}

func (c *Context) UserID() uint {
	sess, err := c.GetUserSession()
	if err != nil {
		return 0
	}
	return sess.UserID
}

func (c *Context) LocationID() uint {
	sess, err := c.GetUserSession()
	if err != nil {
		return 0
	}
	return sess.Location
}
