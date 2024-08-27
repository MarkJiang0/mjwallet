package middleware

import (
	"fmt"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/markjiang0/mjwallet/pkg/errors"
	"github.com/markjiang0/mjwallet/pkg/logging"
	"github.com/markjiang0/mjwallet/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RecoveryConfig struct {
	Skip int // default: 3
}

var DefaultRecoveryConfig = RecoveryConfig{
	Skip: 3,
}

// Recovery from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return RecoveryWithConfig(DefaultRecoveryConfig)
}

func RecoveryWithConfig(config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rv := recover(); rv != nil {
				ctx := c.Request.Context()
				ctx = logging.NewTag(ctx, logging.TagKeyRecovery)

				var fields []zap.Field
				fields = append(fields, zap.StackSkip("stack", config.Skip))

				if gin.IsDebugging() {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					fields = append(fields, zap.Strings("headers", headers))
				}

				logging.Context(ctx).Error(fmt.Sprintf("[Recovery] %s panic recovered", time.Now().Format("2006/01/02 - 15:04:05")), fields...)
				util.ResError(c, errors.InternalServerError("", "Internal server error, please try again later"))
			}
		}()

		c.Next()
	}
}
