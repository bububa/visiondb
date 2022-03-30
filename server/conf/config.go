package conf

import "github.com/bububa/camera"

type Config struct {
	AppName      string         `required:"true"`
	Port         int            `required:"true"`
	UI           string         `required:"true"`
	FaceID       FaceID         `required:"true"`
	HandPose     HandPose       `required:"true"`
	HandGuesture HandPose       `required:"true"`
	Camera       camera.Options `required:"true"`
	FontPath     string
	LogPath      string
	Debug        bool
}

type FaceID struct {
	DetecterModelPath   string
	RecognizerModelPath string
	DatabasePath        string
}

type HandPose struct {
	DetecterModelPath  string
	EstimatorModelPath string
	DatabasePath       string
}
