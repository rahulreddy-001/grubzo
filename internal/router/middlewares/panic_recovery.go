package middlewares

import (
	"fmt"
	"grubzo/internal/router/ext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type zapWriter struct {
	logger *zap.Logger
}

func (w zapWriter) Write(p []byte) (n int, err error) {
	w.logger.Error("panic", zap.String("info", string(p)))
	return len(p), nil
}

func RecoverPanic(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(zapWriter{logger: logger}, func(c *gin.Context, err any) {
		ext.Ctx(c).RespondWithError(ext.Panic(fmt.Errorf("%v", err)))
		c.Abort()
	})
}
