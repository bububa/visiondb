package faceid

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type AddLabelRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func AddLabelHandler(c *gin.Context) {
	var req AddLabelRequest
	if handler.CheckErr(c.ShouldBind(&req), c) {
		return
	}

	err := service.FaceIDService().DB.AddLabel(req.Name)
	if handler.CheckErr(err, c) {
		return
	}
	handler.Success(c, nil)
}
