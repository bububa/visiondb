package handguesture

import (
	"github.com/gin-gonic/gin"
	"github.com/things-go/nocache"

	"github.com/bububa/visiondb/server/handler/api/handguesture"
)

func Router(r *gin.RouterGroup) {
	g := r.Group("/handguesture")
	g.GET("/label/list", handguesture.ListLabelHandler)
	g.GET("/label/:id", handguesture.LabelInfoHandler)
	g.POST("/label/add", handguesture.AddLabelHandler)
	g.POST("/label/del", handguesture.DelLabelHandler)
	g.POST("/label/update", handguesture.UpdateLabelHandler)
	g.POST("/label/reset", handguesture.ResetLabelHandler)
	g.GET("/sample/list", handguesture.ListSampleHandler)
	g.POST("/sample/camera", handguesture.SampleCameraHandler)
	g.POST("/sample/del", handguesture.DelSampleHandler)
	g.GET("/sample/image/:label/:item", nocache.NoCache(), handguesture.SampleImageHandler)
	g.POST("/train", handguesture.TrainHandler)
	g.POST("/classify/camera", handguesture.ClassifyCameraHandler)
}
