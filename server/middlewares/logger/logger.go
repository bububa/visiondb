package logger

import (
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config defines the config for logger middleware
type Config struct {
	Logger *zerolog.Logger
	// UTC a boolean stating whether to use UTC time zone or local.
	UTC            bool
	SkipPath       []string
	SkipPathRegexp *regexp.Regexp
}

// SetLogger initializes the logging middleware.
func SetLogger(config ...Config) gin.HandlerFunc {
	var newConfig Config
	if len(config) > 0 {
		newConfig = config[0]
	}
	var skip map[string]struct{}
	if length := len(newConfig.SkipPath); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range newConfig.SkipPath {
			skip[path] = struct{}{}
		}
	}

	var sublog zerolog.Logger
	if newConfig.Logger == nil {
		sublog = log.Logger
	} else {
		sublog = *newConfig.Logger
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		logger := sublog.With().Logger()
		dumpLogger := &logger
		if requestID := requestid.Get(c); requestID != "" {
			dumpLogger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
				return ctx.Str("X-Request-ID", requestID)
			})
		}
		dumpLogger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
			return ctx.Str("method", c.Request.Method).
				Str("path", path).
				Str("ip", c.ClientIP()).
				Str("remote", remoteIP(c))
		})
		ua, err := url.QueryUnescape(c.Request.UserAgent())
		if err != nil {
			ua = c.Request.UserAgent()
			if ua != "" {
				dumpLogger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
					return ctx.Str("ua", ua)
				})
			}
		}
		c.Set("LOGGER", dumpLogger)

		c.Next()
		track := true

		if _, ok := skip[c.Request.URL.Path]; ok {
			track = false
		}

		if track &&
			newConfig.SkipPathRegexp != nil &&
			newConfig.SkipPathRegexp.MatchString(path) {
			track = false
		}

		if track {
			end := time.Now()
			latency := end.Sub(start)
			if newConfig.UTC {
				end = end.UTC()
			}
			msg := "Request"
			if len(c.Errors) > 0 {
				msg = c.Errors.String()
			}

			dumpLoggerCtx, exists := c.Get("LOGGER")
			if !exists {
				return
			}
			dumpLogger := dumpLoggerCtx.(*zerolog.Logger)
			dumpLogger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
				return ctx.Int("status", c.Writer.Status()).
					Int("size", c.Writer.Size()).
					Dur("latency", latency).
					Str("referer", c.Request.Referer())
			})
			switch {
			case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
				dumpLogger.Warn().
					Msg(msg)
			case c.Writer.Status() >= http.StatusInternalServerError:
				dumpLogger.Error().
					Msg(msg)
			case c.Writer.Status() == http.StatusFound:
				location := c.Writer.Header().Get("Location")
				dumpLogger.Info().Str("redirect", location).Msg(msg)
			default:
				dumpLogger.Info().
					Msg(msg)
			}
		}

	}

}

func remoteIP(c *gin.Context) string {
	if values, _ := c.Request.Header["X-Forwarded-For"]; len(values) > 0 {
		clientIP := values[0]
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}
	if values, _ := c.Request.Header["X-Real-Ip"]; len(values) > 0 {
		clientIP := strings.TrimSpace(values[0])
		if len(clientIP) > 0 {
			return clientIP
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
