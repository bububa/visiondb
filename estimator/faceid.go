package estimator

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"sync"

	"github.com/bububa/openvision/go/common"
	"github.com/bububa/openvision/go/face"
	"github.com/bububa/openvision/go/face/aligner"
	"github.com/bububa/openvision/go/face/detecter"
	"github.com/bububa/openvision/go/face/drawer"
	"github.com/bububa/openvision/go/face/recognizer"

	"github.com/bububa/visiondb/pb"
	"github.com/bububa/visiondb/server/conf"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/utils"
)

// FaceID represents FaceID estimator
type FaceID struct {
	drawer         *drawer.Drawer
	aligner        *aligner.Aligner
	detecter       detecter.Detecter
	recognizer     recognizer.Recognizer
	alignedImgPool *sync.Pool
	bufPool        *sync.Pool
	mutex          *sync.Mutex
}

// NewFaceID returns a new FaceID
func NewFaceID(d detecter.Detecter, r recognizer.Recognizer) *FaceID {
	return &FaceID{
		aligner:    aligner.NewAligner(),
		detecter:   d,
		recognizer: r,
		drawer:     drawer.New(drawer.WithFont(conf.DefaultFont)),
		alignedImgPool: &sync.Pool{
			New: func() interface{} {
				return common.NewImage(nil)
			},
		},
		bufPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		mutex: new(sync.Mutex),
	}
}

// Close .
func (f *FaceID) Close() {
	f.Lock()
	defer f.Unlock()
	f.aligner.Destroy()
	f.detecter.Destroy()
	f.recognizer.Destroy()
}

func (f *FaceID) Lock() {
	f.mutex.Lock()
}

func (f *FaceID) Unlock() {
	f.mutex.Unlock()
}

func (f *FaceID) SetFont(font *common.Font) {
	f.drawer.Font = font
}

// Features extract face face features from image
func (f *FaceID) Features(img image.Image) ([]*pb.Item, []common.Rectangle, error) {
	rgbImg := common.NewImage(img)
	f.Lock()
	defer f.Unlock()
	faces, err := f.detecter.Detect(rgbImg)
	if err != nil {
		return nil, nil, err
	}
	ret := make([]*pb.Item, 0, len(faces))
	alignedImg := f.alignedImgPool.Get().(*common.Image)
	defer f.alignedImgPool.Put(alignedImg)
	rects := make([]common.Rectangle, 0, len(faces))
	for _, face := range faces {
		alignedImg.Reset()
		if err := f.aligner.Align(rgbImg, face, alignedImg); err != nil {
			continue
		}
		features, err := f.recognizer.ExtractFeatures(alignedImg, common.FullRect)
		if err != nil {
			continue
		}
		buf := f.bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		if err := jpeg.Encode(buf, alignedImg.Image, nil); err != nil {
			continue
		}
		itm := new(pb.Item)
		itm.Embedding = utils.FloatSlice64To32(features)
		itm.Raw = make([]byte, buf.Len())
		copy(itm.Raw, buf.Bytes())
		itm.GenHash()
		ret = append(ret, itm)
		f.bufPool.Put(buf)
		rects = append(rects, face.Rect)
	}
	return ret, rects, nil
}

func (f *FaceID) Draw(img image.Image, results []model.ClassifyResult, rois []common.Rectangle) image.Image {
	faces := make([]face.FaceInfo, 0, len(results))
	for idx, ret := range results {
		f := face.FaceInfo{
			Rect: rois[idx],
		}
		if ret.ID >= 0 && ret.Name != "" {
			f.Label = fmt.Sprintf("%s:%.4f", ret.Name, ret.Score)
		}
		faces = append(faces, f)
	}
	return f.drawer.Draw(img, faces)
}
