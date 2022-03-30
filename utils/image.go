package utils

import (
	"image"

	"github.com/bububa/visiondb/server/model"
)

type ImageConverter = func(in image.Image) (out image.Image, err error)
type ImagesClassifier = func(in []image.Image) (out model.ClassifyResult, err error)
