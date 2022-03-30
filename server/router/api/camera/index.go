package camera

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler/api/camera"
)

func Router(r *gin.RouterGroup) {
	g := r.Group("/camera")
	g.GET("/start", camera.StartHandler)
	g.GET("/close", camera.CloseHandler)
	g.GET("/stream", camera.StreamHandler)
	g.POST("/reset", camera.ResetHandler)
}
