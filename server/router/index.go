package router

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"

	"github.com/bububa/visiondb/server/middlewares/logger"
	"github.com/bububa/visiondb/server/router/api"
	"github.com/bububa/visiondb/server/router/ws"
)

func NewRouter() (*gin.Engine, *melody.Melody) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.SetLogger())

	api.Router(r)
	m := ws.Router(r)
	return r, m

}
