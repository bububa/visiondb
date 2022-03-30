package camera

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

func StartHandler(c *gin.Context) {
	service.Camera().Start()
	handler.Success(c, nil)
}
