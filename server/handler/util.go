package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/bububa/visiondb/logger"
)

// Logger
func Logger(c *gin.Context, callerSkip int) zerolog.Logger {
	if loggerCtx, exists := c.Get("LOGGER"); exists {
		return loggerCtx.(*zerolog.Logger).With().CallerWithSkipFrameCount(callerSkip).Logger()
	}
	return logger.Logger.With().CallerWithSkipFrameCount(callerSkip).Logger()
}

// HandlerLogger returns handler logger
func HandlerLogger(c *gin.Context) zerolog.Logger {
	return Logger(c, 2)
}

// CheckErr check error response
func CheckErr(err error, c *gin.Context) (ret bool) {
	ret = err != nil
	if ret {
		hookLogger := Logger(c, 3)
		hookLogger.Error().Err(err).Send()
		c.JSON(http.StatusOK, ErrorResponse(ErrorCode, err.Error()))
	}
	return
}

// CheckWithCode check error with error code
func CheckWithCode(flag bool, code int, err string, c *gin.Context) (ret bool) {
	ret = flag
	if ret {
		hookLogger := Logger(c, 3)
		hookLogger.Error().Msg(err)
		c.JSON(http.StatusOK, ErrorResponse(code, err))
	}
	return
}

// Success returns http succes response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse(data))
}
