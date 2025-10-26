package ce

import (
	"net/http"
	"regexp"
	"strings"

	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/blendle/zapdriver"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

func sanitizeStack(stack string) []string {
	stack = strings.ReplaceAll(stack, "\t", " ")
	re := regexp.MustCompile(`[^\x20-\x7E\n]`)
	stack = re.ReplaceAllString(stack, "")
	lines := strings.Split(stack, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return cleaned
}
func (i *InternalError) JSON() map[string]any {
	return map[string]any{
		"panic":  i.Panic,
		"stack":  sanitizeStack(i.Stack),
		"fields": i.Fields,
		"error":  i.Err.Error(),
	}
}

func InternalServerError(err error) *InternalError {
	return &InternalError{
		Err:   err,
		Stack: fmt.Sprintf("%s", debug.Stack()),
		Fields: []zap.Field{
			zapdriver.ErrorReport(runtime.Caller(2)),
			zap.String("errorString", err.Error()),
			zap.Error(err),
		},
		Panic: false,
	}
}

func Panic(err error) *InternalError {
	return &InternalError{
		Err:   err,
		Stack: fmt.Sprintf("%s", debug.Stack()),
		Fields: []zap.Field{
			zapdriver.ErrorReport(runtime.Caller(2)),
			zap.String("errorString", err.Error()),
			zap.Error(err),
		},
		Panic: true,
	}
}

func RespondWithError(c *gin.Context, err error) {
	if cerr, ok := err.(*Error); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": cerr.Msg,
		})
		return
	}
	if cerr, ok := err.(*InternalError); ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
			"debug": cerr.JSON(),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
		"debug": InternalServerError(err).JSON(),
	})
}

func BadRequestParams(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "invalid request parameters",
	})
}

func BadRequestBody(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "invalid request body",
	})
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}
