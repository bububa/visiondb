package faceid

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type DelLabelRequest struct {
	ID int `json:"id" form:"id" binding:"required"`
}

func DelLabelHandler(c *gin.Context) {
	var req DelLabelRequest
	if handler.CheckErr(c.ShouldBind(&req), c) {
		return
	}

	err := service.FaceIDService().DB.DeleteLabel(req.ID)
	if handler.CheckErr(err, c) {
		return
	}
	handler.Success(c, nil)
}
