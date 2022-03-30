package faceid

import (
	"github.com/gin-gonic/gin"
	"github.com/things-go/nocache"

	"github.com/bububa/visiondb/server/handler/api/faceid"
)

func Router(r *gin.RouterGroup) {
	g := r.Group("/faceid")
	g.GET("/label/list", faceid.ListLabelHandler)
	g.GET("/label/:id", faceid.LabelInfoHandler)
	g.POST("/label/add", faceid.AddLabelHandler)
	g.POST("/label/del", faceid.DelLabelHandler)
	g.POST("/label/update", faceid.UpdateLabelHandler)
	g.POST("/label/reset", faceid.ResetLabelHandler)
	g.GET("/sample/list", faceid.ListSampleHandler)
	g.POST("/sample/add", faceid.AddSampleHandler)
	g.POST("/sample/camera", faceid.SampleCameraHandler)
	g.POST("/sample/del", faceid.DelSampleHandler)
	g.GET("/sample/image/:label/:item", nocache.NoCache(), faceid.SampleImageHandler)
	g.POST("/train", faceid.TrainHandler)
	g.POST("/classify/image", faceid.ClassifyImageHandler)
	g.POST("/classify/camera", faceid.ClassifyCameraHandler)
}
