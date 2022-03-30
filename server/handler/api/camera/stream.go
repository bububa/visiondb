package camera

import (
	"context"
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	imageEncoder "github.com/bububa/camera/image"
	"github.com/gin-gonic/gin"

	"github.com/bububa/visiondb/logger"
	"github.com/bububa/visiondb/server/handler"
	"github.com/bububa/visiondb/server/service"
)

type StreamRequest struct {
	Delay int `form:"delay"`
}

func StreamHandler(c *gin.Context) {
	cam := service.Camera()
	var req StreamRequest
	c.ShouldBind(&req)
	if req.Delay > 0 {
		cam.SetDelay(req.Delay)
	}
	mimeWriter := multipart.NewWriter(c.Writer)
	mimeWriter.SetBoundary("--boundary")

	c.Header("Connection", "close")
	c.Header("Cache-Control", "no-store, no-cache")
	c.Header("Content-Type", fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary()))

	cn := c.Writer.(http.CloseNotifier).CloseNotify()

	logger := handler.HandlerLogger(c)
	err := cam.Start()
	if handler.CheckErr(err, c) {
		return
	}
	defer cam.Stop()
loop:
	for {
		select {
		case <-cn:
			break loop
		case img := <-cam.Read():
			if img == nil {
				break loop
			}
			if err := writeStream(c.Request.Context(), mimeWriter, img); err != nil {
				logger.Error().Err(err).Send()
				if errors.Is(err, context.DeadlineExceeded) {
					break loop
				}
				continue
			}
		}
	}
	mimeWriter.Close()
}

func writeStream(ctx context.Context, mimeWriter *multipart.Writer, img image.Image) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	ch := make(chan error, 1)
	go func() {
		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Type", "image/jpeg")

		if partWriter, err := mimeWriter.CreatePart(partHeader); err != nil {
			logger.Error().Err(err).Send()
			ch <- err
			return
		} else if err = imageEncoder.NewEncoder(partWriter).Encode(img); err != nil {
			logger.Error().Err(err).Send()
			ch <- err
			return
		}
		ch <- nil
	}()
	select {
	case <-timeoutCtx.Done():
		logger.Error().Err(timeoutCtx.Err()).Msg("ctx timeout")
		return timeoutCtx.Err()
	case ret := <-ch:
		return ret
	}
}
