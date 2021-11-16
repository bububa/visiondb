package service

import (
	"github.com/bububa/visiondb/server/conf"
)

// Init initialize services
func Init(config *conf.Config) error {
	if err := initFaceID(config.FaceID); err != nil {
		return err
	}
	return nil
}

// Close resource recycle
func Close() error {
	if err := closeFaceID(); err != nil {
		return err
	}
	return nil
}
