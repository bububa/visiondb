package estimator

import (
	"image"

	"github.com/bububa/openvision/go/common"
	"github.com/bububa/openvision/go/face/aligner"
	"github.com/bububa/openvision/go/face/detecter"
	"github.com/bububa/openvision/go/face/recognizer"
)

// FaceID represents FaceID estimator
type FaceID struct {
	aligner    *aligner.Aligner
	detecter   detecter.Detecter
	recognizer recognizer.Recognizer
}

// NewFaceID returns a new FaceID
func NewFaceID(d detecter.Detecter, r recognizer.Recognizer) *FaceID {
	return &FaceID{
		aligner:    aligner.NewAligner(),
		detecter:   d,
		recognizer: r,
	}
}

// Close .
func (f *FaceID) Close() {
	f.aligner.Destroy()
	f.detecter.Destroy()
	f.recognizer.Destroy()
}

// Features extract face face features from image
func (f *FaceID) Features(img image.Image) ([][]float64, error) {
	rgbImg := common.NewImage(img)
	faces, err := f.detecter.Detect(rgbImg)
	if err != nil {
		return nil, err
	}
	ret := make([][]float64, 0, len(faces))
	for _, face := range faces {
		rect := face.Rect
		features, err := f.recognizer.ExtractFeatures(rgbImg, rect)
		if err != nil {
			continue
		}
		ret = append(ret, features)
	}
	return ret, nil
}
