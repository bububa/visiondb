package handguesture

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/server/service"
)

type LabelInfoRequest struct {
	ID *int `uri:"id" binding:"required"`
}

func LabelInfoHandler(c *gin.Context) {
	var req LabelInfoRequest
	if handler.CheckErr(c.ShouldBindUri(&req), c) {
		return
	}

	labelName, itemsCount, err := service.HandGuestureService().DB.GetLabelByID(*req.ID)
	if handler.CheckErr(err, c) {
		return
	}
	record := model.Record{
		ID:         *req.ID,
		Name:       labelName,
		ItemsCount: itemsCount,
	}
	handler.Success(c, record)
}
