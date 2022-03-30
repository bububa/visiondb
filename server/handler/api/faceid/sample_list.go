package faceid

import (
	"math"

	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/server/service"
)

type SampleListRequest struct {
	ID       *int `json:"id" form:"id" binding:"required"`
	Page     int  `json:"page" form:"page" binding:"required"`
	PageSize int  `json:"page_size" form:"page_size" binding:"required"`
}

func ListSampleHandler(c *gin.Context) {
	var req SampleListRequest
	if handler.CheckErr(c.ShouldBind(&req), c) {
		return
	}

	items, err := service.FaceIDService().DB.GetLabelItems(*req.ID)
	if handler.CheckErr(err, c) {
		return
	}
	total := len(items)
	var pageCount = 1
	if req.PageSize > 0 {
		pageCount = int(math.Ceil(float64(total) / float64(req.PageSize)))
		startPtr := (req.Page - 1) * req.PageSize
		endPtr := startPtr + req.PageSize
		if endPtr > total {
			endPtr = total
		}
		items = items[startPtr:endPtr]
	}
	records := make([]model.Record, 0, len(items))
	for idx, v := range items {
		rec := model.Record{
			ID:   idx + (req.Page-1)*req.PageSize,
			Name: v.GetHash(),
		}
		records = append(records, rec)
	}
	handler.Success(c, gin.H{"page_count": pageCount, "list": records})
}
