package middlewares

import (
	"errors"
	"grubzo/internal/utils/ce"

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
		internalError := ce.Panic(errors.New("Panic recovered"))
		ce.RespondWithError(c, internalError)
	})
}
