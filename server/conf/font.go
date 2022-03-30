package conf

import (
	"github.com/bububa/openvision/go/common"
	"github.com/llgcode/draw2d"
)

var DefaultFont = NewFont(
	&draw2d.FontData{
		Name: "NotoSansCJKsc",
		//Name:   "Roboto",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	},
	9,
	nil,
)

// NewFont returns a new common.Font
func NewFont(data *draw2d.FontData, size float64, cache draw2d.FontCache) *common.Font {
	fnt := &common.Font{
		Data: data,
		Size: size,
	}
	if cache != nil {
		fnt.Load(cache)
	}
	return fnt
}
