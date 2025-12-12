package ext

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/blendle/zapdriver"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CError struct {
	Err string
}

func Error(err string) *CError {
	return &CError{
		Err: err,
	}
}

func (e *CError) Error() string {
	return e.Err
}

type InternalError struct {
	Err    error
	Stack  string
	Fields []zap.Field
	Panic  bool
}

func (i *InternalError) Error() string {
	if i.Panic {
		return fmt.Sprintf("[Panic] %s\n%s", i.Err.Error(), i.Stack)
	}
	return fmt.Sprintf("%s\n%s", i.Err.Error(), i.Stack)
}

func (i *InternalError) JSON() map[string]any {
	return map[string]any{
		"panic":  i.Panic,
		"stack":  i.Stack,
		"fields": i.Fields,
		"error":  i.Err.Error(),
	}
}

func InternalServerError(err error) *InternalError {
	return &InternalError{
		Err:   err,
		Stack: string(debug.Stack()),
		Fields: []zap.Field{
			zapdriver.ErrorReport(runtime.Caller(1)),
			zap.String("error", err.Error()),
		},
		Panic: false,
	}
}

func Panic(err error) *InternalError {
	return &InternalError{
		Err:   err,
		Stack: string(debug.Stack()),
		Fields: []zap.Field{
			zapdriver.ErrorReport(runtime.Caller(1)),
			zap.String("error", err.Error()),
		},
		Panic: true,
	}
}

func (c *Context) RespondWithError(err error) {
	if err, ok := err.(*CError); ok {
		c.ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Err,
		})
		return
	}
	if err, ok := err.(*InternalError); ok {
		c.ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
			"debug": err.JSON(),
		})
		return
	}
	c.ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
		"debug": InternalServerError(err).JSON(),
	})
}

func (c *Context) BadRequestParams() {
	c.ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "invalid request parameters",
	})
}

func (c *Context) BadRequestBody() {
	c.ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "invalid request body",
	})
}

func (c *Context) Unauthorized() {
	c.ctx.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}

func (c *Context) BadRequestWith(err string) {
	c.ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err,
	})
}

func (c *Context) RespondWithOK(obj any) {
	c.ctx.JSON(http.StatusOK, obj)
}