package handpose

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type UpdateLabelRequest struct {
	ID   *int   `json:"id" form:"id" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

func UpdateLabelHandler(c *gin.Context) {
	var req UpdateLabelRequest
	if handler.CheckErr(c.ShouldBind(&req), c) {
		return
	}

	name := strings.TrimSpace(strings.ToLower(req.Name))
	err := service.HandPoseService().DB.ChangeLabelName(*req.ID, name)
	if handler.CheckErr(err, c) {
		return
	}
	handler.Success(c, gin.H{"id": *req.ID, "name": name})
}
