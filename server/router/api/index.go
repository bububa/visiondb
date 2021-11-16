package api

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/router/api/faceid"
)

func Router(r *gin.Engine) {
	g := r.Group("/api")
	faceid.Router(g)
}
