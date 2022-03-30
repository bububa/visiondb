package faceid

import (
	"image"

	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/server/service"
)

type SampleCameraRequest struct {
	ID    *int `json:"id" binding:"required"`
	Delay int  `json:"delay"`
}

func SampleCameraHandler(c *gin.Context) {
	cam := service.Camera()
	var req SampleCameraRequest
	c.ShouldBind(&req)
	if req.Delay > 0 {
		cam.SetDelay(req.Delay)
	}

	converter := func(in image.Image) (image.Image, error) {
		srv := service.FaceIDService()
		items, rects, err := srv.Estimator.Features(in)
		if err != nil {
			return nil, err
		}
		if len(items) == 0 {
			return in, nil
		}
		results := make([]model.ClassifyResult, 1)
		out := srv.Estimator.Draw(in, results, rects[0:1])
		if _, _, err := srv.DB.AddLabelItems(*req.ID, items[0]); err != nil {
			return out, err
		}
		return out, nil
	}

	cam.SetConverter(converter)
	handler.Success(c, nil)
}
