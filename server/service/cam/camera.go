package cam

import (
	"github.com/bububa/camera"
)

func New(opts camera.Options) (*camera.Camera, error) {
	device, err := getDevice(opts)
	if err != nil {
		return nil, err
	}
	return camera.NewCamera(device), nil
}
