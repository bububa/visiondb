package handpose

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

func AddSampleHandler(c *gin.Context) {
	labelID, err := strconv.Atoi(c.GetHeader("X-Label-Id"))
	if handler.CheckErr(err, c) {
		return
	}
	file, err := c.FormFile("file")
	if handler.CheckErr(err, c) {
		return
	}
	fn, err := file.Open()
	if handler.CheckErr(err, c) {
		return
	}
	defer fn.Close()
	srv := service.HandPoseService()
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	fn.Read(head)

	var itemIDs []int
	if filetype.IsImage(head) {
		fn.Seek(0, io.SeekStart)
		img, _, err := image.Decode(fn)
		if handler.CheckErr(err, c) {
			return
		}
		items, _, err := srv.Estimator.Features(img)
		if handler.CheckErr(err, c) {
			return
		}
		startID, totalItems, err := srv.DB.AddLabelItems(labelID, items...)
		if handler.CheckErr(err, c) {
			return
		}
		itemIDs = make([]int, 0, totalItems)
		for i := 0; i < totalItems; i++ {
			itemIDs = append(itemIDs, startID+i)
		}
	}

	handler.Success(c, gin.H{"label_id": labelID, "item_ids": itemIDs})
}
