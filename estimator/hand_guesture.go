package estimator

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"sync"

	"github.com/bububa/openvision/go/common"
	"github.com/bububa/openvision/go/hand/drawer"
	"github.com/bububa/openvision/go/hand/pose3d"

	"github.com/bububa/visiondb/pb"
	"github.com/bububa/visiondb/server/conf"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/utils"
)

const HandGuestureFrameSize = 10

// HandGuesture represents HnadPose estimator
type HandGuesture struct {
	drawer    *drawer.Drawer
	estimator *pose3d.Mediapipe
	mutex     *sync.Mutex
}

// NewHandGuesture returns a new HandGuesture
func NewHandGuesture(e *pose3d.Mediapipe) *HandGuesture {
	return &HandGuesture{
		estimator: e,
		drawer:    drawer.New(drawer.WithFont(conf.DefaultFont)),
		mutex:     new(sync.Mutex),
	}
}

// Close .
func (h *HandGuesture) Close() {
	h.Lock()
	defer h.Unlock()
	h.estimator.Destroy()
}

// Lock .
func (h *HandGuesture) Lock() {
	h.mutex.Lock()
}

// Unlock .
func (h *HandGuesture) Unlock() {
	h.mutex.Unlock()
}

// SetFont set label font
func (h *HandGuesture) SetFont(font *common.Font) {
	h.drawer.Font = font
}

// Detect extract hand features from image
func (h *HandGuesture) Detect(img image.Image) ([]*pb.Item, []common.PalmObject, error) {
	rgbImg := common.NewImage(img)
	h.Lock()
	defer h.Unlock()
	rois, err := h.estimator.Detect(rgbImg)
	if err != nil {
		return nil, nil, err
	}
	ret := make([]*pb.Item, 0, len(rois))
	objs := make([]common.PalmObject, 0, len(rois))
	for _, roi := range rois {
		kl := len(roi.Skeleton3d)
		if err != nil || kl == 0 {
			continue
		}
		embedding := make([]float64, kl*3)
		for idx, pt := range roi.Skeleton3d {
			embedding[idx] = pt.X
			embedding[idx+kl] = pt.Y
			embedding[idx+kl*2] = pt.Z
		}
		itm := new(pb.Item)
		itm.Embedding = utils.FloatSlice64To32(embedding)
		itm.Raw = nil
		itm.GenHash()
		ret = append(ret, itm)
		objs = append(objs, roi)
	}
	return ret, objs, nil
}

// Features extract hand features from images
func (h *HandGuesture) Features(imgs []image.Image) (*pb.Item, []common.PalmObject, error) {
	var (
		embeddings = make([]float64, 0, HandGuestureFrameSize*21*3)
		objs       = make([]common.PalmObject, 0, HandGuestureFrameSize)
	)
	for _, img := range imgs {
		rgbImg := common.NewImage(img)
		h.Lock()
		defer h.Unlock()
		rois, err := h.estimator.Detect(rgbImg)
		if err != nil {
			continue
		}
		if len(rois) == 0 {
			continue
		}
		obj := rois[0]
		kl := len(obj.Skeleton3d)
		if err != nil || kl == 0 {
			continue
		}
		objs = append(objs, obj)
		embedding := make([]float64, kl*3)
		for idx, pt := range obj.Skeleton3d {
			embedding[idx] = pt.X
			embedding[idx+kl] = pt.Y
			embedding[idx+kl*2] = pt.Z
		}
		embeddings = append(embeddings, embedding...)
	}
	itm := new(pb.Item)
	itm.Embedding = utils.FloatSlice64To32(embeddings)
	itm.Raw = nil
	itm.GenHash()
	return itm, objs, nil
}

// Draw draw classifyResult
func (h *HandGuesture) Draw(img image.Image, results []model.ClassifyResult, rois []common.PalmObject) image.Image {
	hands := make([]common.PalmObject, 0, len(results))
	for idx, ret := range results {
		o := rois[idx]
		if ret.ID >= 0 && ret.Name != "" {
			o.Name = fmt.Sprintf("%s:%.4f", ret.Name, ret.Score)
		}
		hands = append(hands, o)
	}
	return h.drawer.DrawPalm(img, hands)
}

// DrawGuesture draw  palm skeleton guesture
func (h *HandGuesture) DrawGuesture(rois []common.PalmObject, size float64, bg string, delay int) *gif.GIF {
	images := make([]*image.Paletted, 0, len(rois))
	delays := make([]int, 0, len(rois))
	opt := gif.Options{NumColors: 256, Quantizer: nil, Drawer: nil}
	for _, roi := range rois {
		img := h.drawer.DrawPalm3D(roi, size, bg)
		buf := new(bytes.Buffer)
		if err := gif.Encode(buf, img, &opt); err != nil {
			continue
		}
		im, err := gif.Decode(buf)
		if err != nil {
			continue
		}
		if i, ok := im.(*image.Paletted); ok {
			images = append(images, i)
			delays = append(delays, delay)
		}
	}
	return &gif.GIF{
		Image: images,
		Delay: delays,
	}
}
