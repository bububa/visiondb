package handguesture

import (
	"bytes"
	"image/gif"
	"net/http"

	"github.com/bububa/openvision/go/common"
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/estimator"
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
	srv := service.HandGuestureService()
	item, err := srv.DB.GetLabelItem(*req.LabelID, *req.ItemID)
	if handler.CheckErr(err, c) {
		return
	}
	embedding := item.GetEmbedding()

	frames := len(embedding) / estimator.HandGuestureFrameSize
	objs := make([]common.PalmObject, 0, frames)
	for frame := 0; frame < frames; frame++ {
		pts := make([]common.Point3d, 0, 21)
		for i := 0; i < 21; i++ {
			pt := common.Pt3d(float64(embedding[i*frame]), float64(embedding[i*frame+21]), float64(embedding[i*frame+21*2]))
			pts = append(pts, pt)
		}
		o := common.PalmObject{
			Skeleton3d: pts,
		}
		objs = append(objs, o)
	}
	out := srv.Estimator.DrawGuesture(objs, 128, "", 8)
	buf := new(bytes.Buffer)
	gif.EncodeAll(buf, out)
	c.Data(http.StatusOK, "image/gif", buf.Bytes())
}
