package faceid

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type SampleImageRequest struct {
	LabelID *int `uri:"label" binding:"required"`
	ItemID  *int `uri:"item" binding:"required"`
}

func SampleImageHandler(c *gin.Context) {
	var req SampleImageRequest
	if handler.CheckErr(c.ShouldBindUri(&req), c) {
		return
	}

	item, err := service.FaceIDService().DB.GetLabelItem(*req.LabelID, *req.ItemID)
	if handler.CheckErr(err, c) {
		return
	}

	c.Data(http.StatusOK, "image/jpeg", item.GetRaw())
}
