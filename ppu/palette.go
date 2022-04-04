package ppu

import (
	"image"
	"image/color"
	"image/draw"
)

// 这个是NES所有的颜色映射，调色板只有16中颜色，使用的是这里面的颜色索引。
var AllColor = [64]uint32{
	0x808080, 0x0000BB, 0x3700BF, 0x8400A6,
	0xBB006A, 0xB7001E, 0xB30000, 0x912600,
	0x7B2B00, 0x003E00, 0x00480D, 0x003C22,
	0x002F66, 0x000000, 0x050505, 0x050505,

	0xC8C8C8, 0x0059FF, 0x443CFF, 0xB733CC,
	0xFF33AA, 0xFF375E, 0xFF371A, 0xD54B00,
	0xC46200, 0x3C7B00, 0x1E8415, 0x009566,
	0x0084C4, 0x111111, 0x090909, 0x090909,

	0xFFFFFF, 0x0095FF, 0x6F84FF, 0xD56FFF,
	0xFF77CC, 0xFF6F99, 0xFF7B59, 0xFF915F,
	0xFFA233, 0xA6BF00, 0x51D96A, 0x4DD5AE,
	0x00D9FF, 0x666666, 0x0D0D0D, 0x0D0D0D,

	0xFFFFFF, 0x84BFFF, 0xBBBBFF, 0xD0BBFF,
	0xFFBFEA, 0xFFBFCC, 0xFFC4B7, 0xFFCCAE,
	0xFFD9A2, 0xCCE199, 0xAEEEB7, 0xAAF7EE,
	0xB3EEFF, 0xDDDDDD, 0x111111, 0x111111,
}

type BgPaletteImp struct {
	data           [16]byte
	UniversalColor byte
}

func NewPalette(data [16]byte) Palette {
	data[0x4] = data[0x0]
	data[0x8] = data[0x0]
	data[0xC] = data[0x0]
	return &BgPaletteImp{data: data, UniversalColor: data[0]}
}

// tilePixelColor 的值为 nameTable和attributeTable组合后的四位
func (p *BgPaletteImp) Color(colorIndex byte) color.RGBA {
	if colorIndex == 0x0 || colorIndex == 0x4 || colorIndex == 0x8 || colorIndex == 0xC {
		colorIndex = 0
	}
	u := AllColor[p.data[colorIndex]]
	return uint32ToRgb(u)
}

func (p *BgPaletteImp) Draw() image.Image {
	blockSize := 32
	img := image.NewRGBA(image.Rect(0, 0, len(p.data)*blockSize, blockSize))
	for i := 0; i < len(p.data); i++ {
		x0, y0 := i*blockSize, 0
		x1, y1 := x0+blockSize, y0+blockSize
		draw.Draw(img, image.Rect(x0, y0, x1, y1), &image.Uniform{C: p.Color(byte(i))}, image.Point{}, draw.Src)
	}
	return img
}

type Palette interface {
	Color(colorIndex byte) color.RGBA
	Draw() image.Image
}
