package faceid

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type DelSampleRequest struct {
	LabelID  *int `json:"label_id" form:"label_id" binding:"required"`
	SampleID *int `json:"sample_id" form:"sample_id" binding:"required"`
}

func DelSampleHandler(c *gin.Context) {
	var req DelSampleRequest
	if handler.CheckErr(c.ShouldBind(&req), c) {
		return
	}

	count, err := service.FaceIDService().DB.DeleteLabelItems(*req.LabelID, *req.SampleID)
	if handler.CheckErr(err, c) {
		return
	}
	handler.Success(c, gin.H{"id": *req.LabelID, "items_count": count})
}
