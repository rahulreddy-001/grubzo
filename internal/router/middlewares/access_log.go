package middlewares

import (
	"grubzo/internal/utils/random"
	"strconv"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func getRequestID(c *gin.Context) string {
	rid := c.Request.Header.Get(echo.HeaderXRequestID)
	if len(rid) == 0 {
		rid = random.SecureAlphaNumeric(16)
	}
	return rid
}

func AccessLogging(logger *zap.Logger, dev bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stop := time.Now()

		req := c.Request
		res := c.Writer
		if dev {
			logger.Sugar().Infof("%3d | %s | %s %s %d", res.Status, stop.Sub(start), req.Method, req.URL, res.Size)
		} else {
			logger.Info("",
				zap.String("requestId", getRequestID(c)),
				zapdriver.HTTP(&zapdriver.HTTPPayload{
					RequestMethod: req.Method,
					Status:        res.Status(),
					UserAgent:     req.UserAgent(),
					RemoteIP:      c.ClientIP(),
					Referer:       req.Referer(),
					Protocol:      req.Proto,
					RequestURL:    req.URL.String(),
					RequestSize:   req.Header.Get(echo.HeaderContentLength),
					ResponseSize:  strconv.Itoa(res.Size()),
					Latency:       strconv.FormatFloat(stop.Sub(start).Seconds(), 'f', 9, 64) + "s",
				}))
		}
	}
}
