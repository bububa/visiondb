package estimator

import (
	"fmt"
	"image"
	"sync"

	"github.com/bububa/openvision/go/common"
	"github.com/bububa/openvision/go/hand/drawer"
	"github.com/bububa/openvision/go/hand/pose3d"

	"github.com/bububa/visiondb/pb"
	"github.com/bububa/visiondb/server/conf"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/utils"
)

// HandPose represents HnadPose estimator
type HandPose struct {
	drawer    *drawer.Drawer
	estimator *pose3d.Mediapipe
	mutex     *sync.Mutex
}

// NewHandPose returns a new HandPose
func NewHandPose(e *pose3d.Mediapipe) *HandPose {
	return &HandPose{
		estimator: e,
		drawer:    drawer.New(drawer.WithFont(conf.DefaultFont)),
		mutex:     new(sync.Mutex),
	}
}

// Close .
func (h *HandPose) Close() {
	h.Lock()
	defer h.Unlock()
	h.estimator.Destroy()
}

// Lock .
func (h *HandPose) Lock() {
	h.mutex.Lock()
}

// Unlock .
func (h *HandPose) Unlock() {
	h.mutex.Unlock()
}

// SetFont set label font
func (h *HandPose) SetFont(font *common.Font) {
	h.drawer.Font = font
}

// Features extract hand features from image
func (h *HandPose) Features(img image.Image) ([]*pb.Item, []common.PalmObject, error) {
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

// Draw draw classifyResult
func (h *HandPose) Draw(img image.Image, results []model.ClassifyResult, rois []common.PalmObject) image.Image {
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

// Draw3D draw s3 palm skeleton
func (h *HandPose) Draw3D(roi common.PalmObject, size float64, bg string) image.Image {
	return h.drawer.DrawPalm3D(roi, size, bg)
}
