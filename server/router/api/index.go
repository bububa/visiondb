package api

import (
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/router/api/faceid"
	"github.com/bububa/visiondb/server/router/api/handpose"
	"github.com/bububa/visiondb/server/router/api/handguesture"
	"github.com/bububa/visiondb/server/router/api/camera"
)

func Router(r *gin.Engine) {
	g := r.Group("/api")
	faceid.Router(g)
	handpose.Router(g)
    handguesture.Router(g)
	camera.Router(g)
}
