package service

import (
	"github.com/bububa/visiondb/server/conf"
)

// Init initialize services
func Init(config *conf.Config) error {
	if err := initGlobalFontCache(config.FontPath); err != nil {
		return err
	}
	if err := initFaceID(config.FaceID); err != nil {
		return err
	}
	if err := initHandPose(config.HandPose); err != nil {
		return err
	}
	if err := initHandGuesture(config.HandGuesture); err != nil {
		return err
	}
	if err := initCamera(config.Camera); err != nil {
		return err
	}
	return nil
}

// Close resource recycle
func Close() error {
	closeCamera()
	if err := closeFaceID(); err != nil {
		return err
	}
	if err := closeHandPose(); err != nil {
		return err
	}
	if err := closeHandGuesture(); err != nil {
		return err
	}
	return nil
}
