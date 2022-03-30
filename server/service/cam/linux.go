//go:build linux
// +build linux

package cam

import (
	"github.com/bububa/camera"
	"github.com/bububa/camera/linux"
)

func getDevice(opts camera.Options) (camera.Device, error) {
	return linux.New(opts)
}
