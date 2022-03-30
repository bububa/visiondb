package handpose

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

func TrainHandler(c *gin.Context) {
	err := service.HandPoseService().Train()
	if handler.CheckErr(err, c) {
		return
	}
	handler.Success(c, nil)
}
