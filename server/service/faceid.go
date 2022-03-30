package service

import (
	"image"
	"os"

	"github.com/bububa/openvision/go/classifier/svm"
	"github.com/bububa/openvision/go/face/detecter"
	"github.com/bububa/openvision/go/face/recognizer"

	"github.com/bububa/visiondb/estimator"
	"github.com/bububa/visiondb/logger"
	"github.com/bububa/visiondb/pb"
	"github.com/bububa/visiondb/server/conf"
	"github.com/bububa/visiondb/server/model"
	"github.com/bububa/visiondb/storage"
	"github.com/bububa/visiondb/utils"
)

type faceIDService struct {
	dbPath     string
	Classifier svm.Classifier
	Trainer    svm.Trainer
	Estimator  *estimator.FaceID
	DB         storage.Storage
}

func (s *faceIDService) Close() error {
	s.Estimator.Close()
	s.Trainer.Destroy()
	s.Classifier.Destroy()
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
	if err := d.LoadModel(utils.PathJoin(wd, config.DetecterModelPath)); err != nil {
		return err
	}
	r := recognizer.NewMobilefacenet()
	if err := r.LoadModel(utils.PathJoin(wd, config.RecognizerModelPath)); err != nil {
		return err
	}
	trainer := svm.NewMultiClassTrainer()
	dbPath := utils.PathJoin(wd, config.DatabasePath)
	dbFile := utils.PathJoin(dbPath, "db")
	db := storage.NewProtoBufStorage(dbFile)
	if err := db.Reload(); err != nil {
		return err
	}
	if shape, _ := db.Shape(); shape == nil {
		db.SetShape(&pb.Shape{
			Width:  128,
			Height: 1,
		})
	}
	classifier := svm.NewMultiClassClassifier()
	modelFile := utils.PathJoin(dbPath, "model")
	if _, err := os.Stat(modelFile); err == nil {
		classifier.LoadModel(modelFile)
	}
	faceIDInstance = &faceIDService{
		dbPath:     dbPath,
		Trainer:    trainer,
		Classifier: classifier,
		Estimator:  estimator.NewFaceID(d, r),
		DB:         db,
	}
	return nil
}

func (s *faceIDService) Classify(vec []float32, result *model.ClassifyResult) error {
	scores, err := s.Classifier.Classify(vec)
	if err != nil {
		return err
	}
	var (
		labelID = -1
		score   float64
	)
	for idx, v := range scores {
		if v > 0 && v >= score {
			labelID = idx
			score = v
		}
	}
	result.ID = labelID
	result.Score = score
	result.Scores = scores
	if result.ID < 0 {
		return nil
	}
	if result.Name, _, err = s.DB.GetLabelByID(result.ID); err != nil {
		return err
	}
	return nil
}

func (s *faceIDService) Train() error {
	modelFile := utils.PathJoin(s.dbPath, "model")
	s.Trainer.Reset()
	labels, _, err := s.DB.Labels()
	if err != nil {
		return err
	}
	shape, err := s.DB.Shape()
	if err != nil {
		return err
	}
	s.Trainer.SetLabels(len(labels))
	s.Trainer.SetFeatures(int(shape.GetWidth()))
	for labelID := range labels {
		items, err := s.DB.GetLabelItems(labelID)
		if err != nil {
			return err
		}
		for _, itm := range items {
			s.Trainer.AddData(labelID+1, itm.GetEmbedding())
		}
	}
	s.Trainer.Train(modelFile)
	s.Classifier.LoadModel(modelFile)
	return nil
}

func (s *faceIDService) ClassifyImage(in image.Image) (image.Image, error) {
	items, rects, err := s.Estimator.Features(in)
	if err != nil {
		return nil, err
	}
	results := make([]model.ClassifyResult, 0, len(items))
	var matched model.ClassifyResult
	for _, itm := range items {
		embedding := itm.GetEmbedding()
		if err := s.Classify(embedding, &matched); err != nil {
			logger.Error().Err(err).Send()
		}
		results = append(results, matched)
	}
	return s.Estimator.Draw(in, results, rects), nil

}
