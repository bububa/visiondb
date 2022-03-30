//go:build android
// +build android

package cam

import (
	"github.com/bububa/camera"
	"github.com/bububa/camera/android"
)

func getDevice(opts camera.Options) (camera.Device, error) {
	return android.New(opts)
}
