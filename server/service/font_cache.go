package service

import (
	"github.com/bububa/visiondb/server/conf"
	"github.com/llgcode/draw2d"
)

var globalFontCache draw2d.FontCache

// NewFontCache load font cache
func initGlobalFontCache(fontFolder string) error {
	if fontFolder == "" {
		return nil
	}
	globalFontCache = draw2d.NewSyncFolderFontCache(fontFolder)
	return conf.DefaultFont.Load(globalFontCache)
}
