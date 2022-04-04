package ppu

import (
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
		tileColorData := getColorDataFromPatternTable(int(patternIndex), patternTable)
		colorMaskByte := getColorMaskDataFromAttributeTable(tileIndex, attributeTable)
		//fmt.Printf("colorMaskByte: %b", colorMaskByte)
		//fmt.Println(colorMaskByte)
		mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
		//printInGrid(tileColorData)
		tileRGBData := tileColorData2RGBData(&tileColorData, palette)
		tileImage := tileImage(&tileRGBData)

		appendTile2Screen(m, tileIndex, tileImage)
	}
	Save2jpeg("bg", m)
	return nil
}
func GenTileImage(nameTable []byte, patternTable []byte, attributeTable []byte, palette []byte) []*image.RGBA {
	res := make([]*image.RGBA, 0, len(nameTable))
	for tileIndex, patternIndex := range nameTable {
		tileColorData := getColorDataFromPatternTable(int(patternIndex), patternTable)
		colorMaskByte := getColorMaskDataFromAttributeTable(tileIndex, attributeTable)
		//fmt.Printf("colorMaskByte: %b", colorMaskByte)
		//fmt.Println(colorMaskByte)
		mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
		//printInGrid(tileColorData)
		tileRGBData := tileColorData2RGBData(&tileColorData, palette)
		tileImage := tileImage(&tileRGBData)
		res = append(res, tileImage)
	}
	return res
}

func ScreenRec() image.Rectangle {
	return image.Rect(0, 0, 8*32, 8*30)
}

func GenBg(nameTable []byte, patternTable []byte, attributeTable []byte, palette Palette) image.Image {
	imgs := make([]*image.RGBA, 0, len(nameTable))
	for tileIndex, patternIndex := range nameTable {
		tileColorData := getColorDataFromPatternTable(int(patternIndex), patternTable)
		colorMaskByte := getColorMaskDataFromAttributeTable(tileIndex, attributeTable)
		mergeRgbDataWithColorMaskByte(&tileColorData, colorMaskByte)
		tileRGBData := tileColorData2RGBData2(&tileColorData, palette)
		tileImage := tileImage(&tileRGBData)
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
func GenPalette(rgbData []uint32) {
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
	Save2jpeg("palette", m)

}
func DrawPatternTable(patternTable []byte) {
	for patternIndex := 0; patternIndex < 16; patternIndex++ {
		//tileColorData := getColorDataFromPatternTable(patternIndex, patternTable)

	}
}
