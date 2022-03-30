package camera

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

func CloseHandler(c *gin.Context) {
	service.Camera().Stop()
	handler.Success(c, nil)
}
