package camera

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

func ResetHandler(c *gin.Context) {
	service.Camera().SetConverter(nil)
	handler.Success(c, nil)
}
