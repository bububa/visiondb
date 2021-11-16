package service

import (
	"os"

	"github.com/bububa/openvision/go/face/detecter"
	"github.com/bububa/openvision/go/face/recognizer"

	"github.com/bububa/visiondb/estimator"
	"github.com/bububa/visiondb/server/conf"
	"github.com/bububa/visiondb/storage"
	"github.com/bububa/visiondb/utils"
)

type faceIDService struct {
	Estimator *estimator.FaceID
	DB        storage.Storage
}

func (s *faceIDService) Close() error {
	s.Estimator.Close()
	return s.DB.Flush()
}

var faceIDInstance *faceIDService

// FaceIDService represents faceid service
func FaceIDService() *faceIDService {
	return faceIDInstance
}

func closeFaceID() error {
	if faceIDInstance == nil {
		return nil
	}
	return faceIDInstance.Close()
}

func initFaceID(config conf.FaceID) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	d := detecter.NewRetinaFace()
	if err := d.LoadModel(utils.PathJoin(wd, config.DetectorModelPath)); err != nil {
		return err
	}
	r := recognizer.NewMobilefacenet()
	if err := r.LoadModel(utils.PathJoin(wd, config.RecognizerModelPath)); err != nil {
		return err
	}
	faceIDInstance = &faceIDService{
		Estimator: estimator.NewFaceID(d, r),
		DB:        storage.NewProtoBufStorage(config.DatabasePath),
	}

	return nil
}
