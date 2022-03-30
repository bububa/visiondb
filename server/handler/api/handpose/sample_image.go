package handpose

import (
	"bytes"
	"net/http"

	imgEncoder "github.com/bububa/camera/image"
	"github.com/bububa/openvision/go/common"
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
	srv := service.HandPoseService()
	item, err := srv.DB.GetLabelItem(*req.LabelID, *req.ItemID)
	if handler.CheckErr(err, c) {
		return
	}
	embedding := item.GetEmbedding()

	l := len(embedding) / 3
	pts := make([]common.Point3d, 0, l)
	for i := 0; i < l; i++ {
		pt := common.Pt3d(float64(embedding[i]), float64(embedding[i+l]), float64(embedding[i+l*2]))
		pts = append(pts, pt)
	}
	o := common.PalmObject{
		Skeleton3d: pts,
	}
	out := srv.Estimator.Draw3D(o, 128, "")
	buf := new(bytes.Buffer)
	imgEncoder.NewEncoder(buf).Encode(out)
	c.Data(http.StatusOK, "image/jpeg", buf.Bytes())
}
