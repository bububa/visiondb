package faceid

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type ClassifyCameraRequest struct {
	Delay int `form:"delay"`
}

func ClassifyCameraHandler(c *gin.Context) {
	cam := service.Camera()
	srv := service.FaceIDService()
	cam.SetConverter(srv.ClassifyImage)
	handler.Success(c, nil)
}
