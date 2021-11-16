package faceid

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler/api/faceid"
)

func Router(r *gin.RouterGroup) {
	g := r.Group("/faceid")
	g.GET("/label/list", faceid.ListLabelHandler)
}
