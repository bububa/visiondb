package handguesture

import (
	"math"

	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/server/service"
)

type LabelListRequest struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"page_size" form:"page_size" binding:"required"`
}

func ListLabelHandler(c *gin.Context) {
	var req LabelListRequest
	if handler.CheckErr(c.ShouldBind(&req), c) {
		return
	}

	labels, counts, err := service.HandGuestureService().DB.Labels()
	if handler.CheckErr(err, c) {
		return
	}
	total := len(labels)
	var pageCount = 1
	if req.PageSize > 0 {
		pageCount = int(math.Ceil(float64(total) / float64(req.PageSize)))
		// end := req.Page * req.PageSize
		// if end > total {
		// 	end = total
		// }
		// labels = labels[(req.Page-1)*req.PageSize : total]
	}
	records := make([]model.Record, total)
	for idx, v := range labels {
		records[idx] = model.Record{
			ID:         idx,
			Name:       v,
			ItemsCount: counts[idx],
		}
	}
	handler.Success(c, gin.H{"page_count": pageCount, "list": records})
}
