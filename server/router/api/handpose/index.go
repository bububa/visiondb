package handpose

import (
	"github.com/gin-gonic/gin"
	"github.com/things-go/nocache"

	"github.com/bububa/visiondb/server/handler/api/handpose"
)

func Router(r *gin.RouterGroup) {
	g := r.Group("/handpose")
	g.GET("/label/list", handpose.ListLabelHandler)
	g.GET("/label/:id", handpose.LabelInfoHandler)
	g.POST("/label/add", handpose.AddLabelHandler)
	g.POST("/label/del", handpose.DelLabelHandler)
	g.POST("/label/update", handpose.UpdateLabelHandler)
	g.POST("/label/reset", handpose.ResetLabelHandler)
	g.GET("/sample/list", handpose.ListSampleHandler)
	g.POST("/sample/add", handpose.AddSampleHandler)
	g.POST("/sample/camera", handpose.SampleCameraHandler)
	g.POST("/sample/del", handpose.DelSampleHandler)
	g.GET("/sample/image/:label/:item", nocache.NoCache(), handpose.SampleImageHandler)
	g.POST("/train", handpose.TrainHandler)
	g.POST("/classify/image", handpose.ClassifyImageHandler)
	g.POST("/classify/camera", handpose.ClassifyCameraHandler)
}
