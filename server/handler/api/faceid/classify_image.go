package faceid

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	camImg "github.com/bububa/camera/image"
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/server/service"
)

func ClassifyImageHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if handler.CheckErr(err, c) {
		return
	}
	fn, err := file.Open()
	if handler.CheckErr(err, c) {
		return
	}
	defer fn.Close()
	srv := service.FaceIDService()
	img, _, err := image.Decode(fn)
	if handler.CheckErr(err, c) {
		return
	}
	items, rects, err := srv.Estimator.Features(img)
	if handler.CheckErr(err, c) {
		return
	}
	results := make([]model.ClassifyResult, 0, len(items))
	var matched model.ClassifyResult
	for _, itm := range items {
		embedding := itm.GetEmbedding()
		err := srv.Classify(embedding, &matched)
		if handler.CheckErr(err, c) {
			return
		}
		results = append(results, matched)
	}
	buf := new(bytes.Buffer)
	out := srv.Estimator.Draw(img, results, rects)
	err = jpeg.Encode(buf, out, nil)
	if handler.CheckErr(err, c) {
		return
	}
	b64 := camImg.EncodeToString(buf.Bytes())
	link := fmt.Sprintf("data:image/jpeg;base64,%s", b64)
	handler.Success(c, gin.H{"image": link, "items": results})
}
