package ws

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"

	"github.com/bububa/visiondb/server/handler/ws"
)

func Router(r *gin.Engine) *melody.Melody {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 100
	m.HandleConnect(ws.ConnectHandler)

	m.HandleDisconnect(ws.DisconnectHandler)

	m.HandleMessage(ws.MsgHandler)

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	return m
}
