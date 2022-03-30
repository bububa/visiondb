package handpose

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type ResetLabelRequest struct {
	ID *int `json:"id" form:"id" binding:"required"`
}

func ResetLabelHandler(c *gin.Context) {
	var req ResetLabelRequest
	if handler.CheckErr(c.ShouldBind(&req), c) {
		return
	}

	err := service.HandPoseService().DB.ResetLabelItems(*req.ID)
	if handler.CheckErr(err, c) {
		return
	}
	handler.Success(c, gin.H{"id": *req.ID})
}
