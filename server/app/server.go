package app

import (
	"fmt"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	"github.com/bububa/visiondb/logger"
	IRouter "github.com/bububa/visiondb/server/router"
)

func StartServer(c *cli.Context) error {
	if c.IsSet("port") {
		config.Port = c.Int("port")
	}
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r, m := IRouter.NewRouter()
	logger.Info().Str("service", serverName).Int("port", config.Port).Str("GitTag", GitSummary).Str("GitRevision", GitRevision).Str("STATUS", "START").Send()
	defer logger.Info().Str("service", serverName).Int("port", config.Port).Str("GitTag", GitSummary).Str("GitRevision", GitRevision).Str("STATUS", "EXIT").Send()
	srv := endless.NewServer(fmt.Sprintf(":%d", config.Port), r)
	srv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT] = append(
		srv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT],
		func() {
			if m == nil {
				return
			}
			m.CloseWithMsg([]byte("service shuting down"))
		},
	)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Error().Err(err)
		return err
	}
	logger.Info().Str("service", config.AppName).Int("port", config.Port).Str("GitTag", GitSummary).Str("GitRevision", GitRevision).Str("STATUS", "SHUTTING DOWN").Send()
	return nil
}
