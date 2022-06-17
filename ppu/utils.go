package ppu

import (
	"fc-emulator/utils"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strings"
)

//像素:256x240
//一个tile: 8x8
//所以整个图像的Tile分布: 32x30
//
//画背景图流程：
//从NameTable开始解析
//NameTabel的一个字节代表一个Tile
// 使用该字节+PPUCTRL的指示，去PatternTable取出颜色Tile数据(16byte)
// 该16byte决定了Tile中每个像素点颜色数据的后两位bit (8*8*2 bit = 16byte)
//
//从AttributeTable中获取每个像素点的前两个bit.
//
//然后根据Tile的颜色索引，去拿到实际的颜色。
//至此，一个那么至此一个Tile的图像数据还原完毕。
//
//NameTable组织形式：一个字节代表一个Tile，该字节为PatternTable中的索引
//一个32*30个Tile。
//
//PatternTable组织形式: 16byte为一个单位，每16byte表示一个Tile。包含Tile中
//每个像素点的后两bit。
//
//AttributeTable: 64字节。
//按Tile把屏幕分成4*4的方块。
//屏幕中共有64个方块(最后一行的方块是不完整的)
//AttributeTable的表示64个区域的颜色数据。
//每个区域(16个tile), 1byte,
//每个区域再分为4份，每一份2*2个tile
//这1byte的每2bit平分给这四份。控制着这些Tile像素点颜色的前两位。
func Save2jpeg(filename string, im image.Image) {
	if !strings.HasSuffix(filename, ".jpeg") {
		filename += ".jpeg"
	}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f, im, nil)
	if err != nil {
		panic(err)
	}
}

func DrawImage(nameTable []byte, patternTable []byte, attributeTable []byte, palette []byte) []byte {
	m := image.NewRGBA(image.Rect(0, 0, 256, 240))
	for tileIndex, patternIndex := range nameTable {
		tileColorData := GetColorDataFromPatternTable(int(patternIndex), patternTable)
		colorMaskByte := getColorMaskDataFromAttributeTable(tileIndex, attributeTable)
		//fmt.Printf("colorMaskByte: %b", colorMaskByte)
		//fmt.Println(colorMaskByte)
		mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
		//printInGrid(tileColorData)
		tileRGBData := tileColorData2RGBData(&tileColorData, palette)
		tileImage := convert2tileImage(&tileRGBData)

		appendTile2Screen(m, tileIndex, tileImage)
	}
	Save2jpeg("bg", m)
	return nil
}
func GenTileImage(nameTable []byte, patternTable []byte, attributeTable []byte, palette []byte) []*image.RGBA {
	res := make([]*image.RGBA, 0, len(nameTable))
	for tileIndex, patternIndex := range nameTable {
		tileColorData := GetColorDataFromPatternTable(int(patternIndex), patternTable)
		colorMaskByte := getColorMaskDataFromAttributeTable(tileIndex, attributeTable)
		//fmt.Printf("colorMaskByte: %b", colorMaskByte)
		//fmt.Println(colorMaskByte)
		mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
		//printInGrid(tileColorData)
		tileRGBData := tileColorData2RGBData(&tileColorData, palette)
		tileImage := convert2tileImage(&tileRGBData)
		res = append(res, tileImage)
	}
	return res
}

func ScreenRec() image.Rectangle {
	return image.Rect(0, 0, 8*32, 8*30)
}
func NewScreenImage() *image.RGBA {
	res := image.NewRGBA(ScreenRec())
	draw.Draw(res, res.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)
	return res
}
func DrawPatternTable(patternTable []byte, attributeColorMaskFn AttributeColorMask, palette Palette) *image.RGBA {
	if len(patternTable)/16 != 256 {
		panic(fmt.Sprintf("tile count is not 256, is %v", len(patternTable)/16))
	}
	nameTable := make([]byte, 0, 256)
	for i := 0; i < 256; i++ {
		nameTable = append(nameTable, byte(i))
	}
	imgs := make([]*image.RGBA, 0, 256)
	for tileIndex, patternIndex := range nameTable {
		tileColorData := GetColorDataFromPatternTable(int(patternIndex), patternTable)
		colorMaskByte := attributeColorMaskFn(tileIndex)
		mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
		tileRGBData := tileColorData2RGBData2(&tileColorData, palette)
		tileImage := convert2tileImage(&tileRGBData)
		imgs = append(imgs, tileImage)
	}
	return ComposeImage(imgs, 8, 16)
}

func ScaleImage(img image.Image, scale int) image.Image {
	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	return resize.Resize(uint(w*scale), uint(h*scale), img, resize.Lanczos3)
	//if scale <= 1 {
	//	return img
	//}
	//minP := img.Bounds().Min
	//maxP := img.Bounds().Max
	//xLen := (maxP.X - minP.X) * scale
	//yLen := (maxP.Y - minP.Y) * scale
	//res := image.NewRGBA(image.Rect(0, 0, xLen, yLen))
	//for y := 0; y < maxP.Y; y++ {
	//	for x := 0; x < maxP.X; x++ {
	//		c := img.At(x, y)
	//		xS := x * scale
	//		yS := y * scale
	//		draw.Draw(res, image.Rect(xS, yS, xS+scale, yS+scale), &image.Uniform{C: c}, image.Point{}, draw.Src)
	//	}
	//}
	//return res
}

func ComposeImage(imags []*image.RGBA, blockSize int, rowCount int) *image.RGBA {
	w := rowCount * blockSize
	colCount := len(imags) / rowCount
	if len(imags)%rowCount != 0 {
		colCount += 1
	}
	h := colCount * blockSize
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	for i, img := range imags {
		scaleFactor := blockSize / imags[0].Bounds().Max.X
		img = ScaleImage(img, scaleFactor).(*image.RGBA)
		x := (i % rowCount) * blockSize
		y := (i / rowCount) * blockSize
		draw.Draw(res, image.Rect(x, y, x+blockSize, y+blockSize), img, image.Point{}, draw.Src)
	}
	//withGrid := true
	// 画竖线
	for i := 1; i < rowCount; i++ {
		x := i * blockSize
		draw.Draw(res, image.Rect(x, 0, x+1, h), image.NewUniform(color.White), image.Point{}, draw.Src)
	}
	// 画横线
	for i := 1; i < colCount; i++ {
		y := i * blockSize
		draw.Draw(res, image.Rect(0, y, w, y+1), image.NewUniform(color.White), image.Point{}, draw.Src)
	}
	return res
}

func DrawBg(nameTable []byte, patternTable []byte, attributeColorMaskFn AttributeColorMask, palette Palette) image.Image {
	imgs := make([]*image.RGBA, 0, len(nameTable))
	for tileIndex, patternIndex := range nameTable {
		tileColorData := GetColorDataFromPatternTable(int(patternIndex), patternTable)
		colorMaskByte := attributeColorMaskFn(tileIndex)
		mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
		tileRGBData := tileColorData2RGBData2(&tileColorData, palette)
		tileImage := convert2tileImage(&tileRGBData)
		imgs = append(imgs, tileImage)
	}
	sc := image.NewRGBA(ScreenRec())
	for i, img := range imgs {
		x := (i % 32) * 8
		y := (i / 32) * 8
		draw.Draw(sc, image.Rect(x, y, x+8, y+8), img, image.Point{}, draw.Src)
	}
	return sc
}

type Sprite struct {
	X               int
	Y               int
	FlipVertical    bool
	FlipHorizontal  bool
	BehindBg        bool
	OriginColorData *[8][8]byte
	ColorData       *[8][8]byte
	OriginImage     *image.RGBA
	Image           *image.RGBA
}

func DrawSpritesTable(oam [256]byte, patternTable []byte, palette Palette) image.Image {
	sprites := make([]*Sprite, 0, 64)
	for i := 0; i < 256; i = i + 4 {
		spriteData := [4]byte{oam[i], oam[i+1], oam[i+2], oam[i+3]}
		sprites = append(sprites, drawSprite(spriteData, patternTable, palette))
	}
	return ComposeSprites(sprites)
}
func DrawSpritesInSC(screenImage *image.RGBA, oam [256]byte, patternTable []byte, palette Palette) image.Image {
	sprites := make([]*Sprite, 0, 64)
	for i := 0; i < 256; i = i + 4 {
		spriteData := [4]byte{oam[i], oam[i+1], oam[i+2], oam[i+3]}
		sprites = append(sprites, drawSprite(spriteData, patternTable, palette))
	}
	for _, sprite := range sprites {
		draw.Draw(screenImage, image.Rect(sprite.X, sprite.Y, sprite.X+8, sprite.Y+8),
			sprite.Image, image.Point{}, draw.Src)
	}
	return screenImage
}

func ComposeSprites(sprites []*Sprite) image.Image {
	images := make([]*image.RGBA, 0, len(sprites))
	for _, sprite := range sprites {
		images = append(images, sprite.OriginImage)
	}
	return ComposeImage(images, 32, 8)
}
func drawSprite(spriteData [4]byte, patternTable []byte, palette Palette) *Sprite {
	y := spriteData[0]
	tileIndex := spriteData[1]
	tileColorData := GetColorDataFromPatternTable(int(tileIndex), patternTable)
	colorMaskByte := (spriteData[2] & 0b11) << 2

	mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
	tileRGBData := tileColorData2RGBData2(&tileColorData, palette)

	tileImage := convert2tileImage(&tileRGBData)
	sprite := &Sprite{
		X:               int(spriteData[3]),
		Y:               int(y),
		FlipVertical:    utils.IsSet(spriteData[2], 7),
		FlipHorizontal:  utils.IsSet(spriteData[2], 6),
		BehindBg:        true,
		OriginColorData: &tileColorData,
		OriginImage:     tileImage,
	}
	sprite.ColorData = Flip(*sprite.OriginColorData, sprite.FlipVertical, sprite.FlipHorizontal)
	tileRGBData = tileColorData2RGBData2(sprite.ColorData, palette)
	sprite.Image = convert2tileImage(&tileRGBData)
	return sprite
}

func Flip(rgb [8][8]byte, flipVertical, flipHorizontal bool) *[8][8]byte {
	if flipVertical {
		for y := 0; y < 8; y++ {
			for x := 0; x < 4; x++ {
				rgb[y][x], rgb[y][7-x] = rgb[y][7-x], rgb[y][x]
			}
		}
	}
	if flipHorizontal {
		for x := 0; x < 8; x++ {
			for y := 0; y < 4; y++ {
				rgb[y][x], rgb[7-y][x] = rgb[7-y][x], rgb[y][x]
			}
		}
	}
	return &rgb

}

func SystemPaletteColor() *image.RGBA {
	c := make([]uint32, 64, 64)
	for i := 0; i < len(AllColor); i++ {
		c[i] = AllColor[i]
	}
	return GenPalette(c)
}
func GenPalette(rgbData []uint32) *image.RGBA {
	blockSize := 32
	m := image.NewRGBA(image.Rect(0, 0, 16*blockSize, 4*blockSize))
	count := 0
	for y := 0; y < 4; y++ {
		for x := 0; x < 16; x++ {
			xStart := x * blockSize
			yStart := y * blockSize
			block := image.Rect(xStart, yStart, xStart+blockSize, yStart+blockSize)
			rgb := rgbData[y*16+x]
			//fmt.Println(y*16+x)
			c := color.RGBA{
				R: uint8((rgb >> 16) & 0xFF),
				G: uint8((rgb >> 8) & 0xFF),
				B: uint8(rgb & 0xFF),
				A: 255,
			}
			//fmt.Println(fmt.Sprintf("%X", c.R),
			//	fmt.Sprintf("%X", c.G),
			//	fmt.Sprintf("%X", c.B))

			draw.Draw(m, block, &image.Uniform{C: c}, image.Point{0, 0}, draw.Src)
			count += 1
			//if count == 16 {
			//	Save2jpeg("palette", m)
			//	return
			//}
		}
	}
	return m
	//Save2jpeg("palette", m)

}
