package ppu

import (
	"fc-emulator/rom"
	"fc-emulator/utils"
	"fmt"
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func genImage() {
	m := image.NewRGBA(image.Rect(0, 0, 256, 240))

	draw.Draw(m, m.Bounds(), &image.Uniform{C: color.Black},
		image.Point{}, draw.Src)
	draw.Draw(m, image.Rect(20, 20, 100, 100), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	f, err := os.Create("demo.jpeg")
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		panic(err)
	}
}

//https://wangbjun.site/2020/coding/golang/image.html
func TestGenImage(t *testing.T) {
	m := image.NewRGBA(image.Rect(0, 0, 256, 240))
	b := m.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			m.Set(x, y, color.NRGBA{
				R: uint8(y),
				G: 0,
				B: uint8(x),
				A: 255,
			})
		}
	}
	Save2jpeg("demo", m)
}
func TestGetBit(t *testing.T) {
	for i := 0; i < 8; i++ {
		fmt.Print(utils.GetBitFromLeft(24, i))
	}
}
func TestShift(t *testing.T) {

}

func TestDrawImage(t *testing.T) {
	nameTables := []byte{}
	for i := 0; i < 0xFF; i++ {
		nameTables = append(nameTables, byte(i))
	}
	palette := []byte{
		0x22, 0x29, 0x1A, 0x0F, // 00
		0x22, 0x36, 0x17, 0x0F, // 01
		0x22, 0x30, 0x21, 0x0F, // 10
		0x22, 0x17, 0x17, 0x0F, // 11
	}
	nesRom, err := rom.LoadNesRom("../static/mario.nes")
	require.NoError(t, err)
	//patternTable := nesRom.ChrRom[4*utils.Kb:8*utils.Kb]
	patternTable := nesRom.ChrRom[0 : 4*utils.Kb]
	attributeTable := []byte{}
	for i := 0; i < 64; i++ {
		attributeTable = append(attributeTable, 0)
	}
	tileImages := GenTileImage(nameTables, patternTable, attributeTable, palette)
	blockSize := 32
	scaleFactor := blockSize / 8
	m := image.NewRGBA(image.Rect(0, 0, blockSize*16, blockSize*16))
	for i, tileImage := range tileImages {
		x := (i % 16) * blockSize
		y := (i / 16) * blockSize
		draw.Draw(m, image.Rect(x, y, x+blockSize, y+blockSize), ScaleImage(tileImage, scaleFactor), image.Point{}, draw.Src)
	}
	Save2jpeg("patternTable00", m)
	OpenImageFile("patternTable00")
}

func ScaleImage(img *image.RGBA, scale int) *image.RGBA {
	if scale <= 1 {
		return img
	}
	minP := img.Bounds().Min
	maxP := img.Bounds().Max
	fmt.Println(img.Bounds())
	xLen := (maxP.X - minP.X) * scale
	yLen := (maxP.Y - minP.Y) * scale
	res := image.NewRGBA(image.Rect(0, 0, xLen, yLen))
	for y := 0; y < maxP.Y; y++ {
		for x := 0; x < maxP.X; x++ {
			c := img.At(x, y)
			xS := x * scale
			yS := y * scale
			draw.Draw(res, image.Rect(xS, yS, xS+scale, yS+scale), &image.Uniform{C: c}, image.Point{}, draw.Src)
		}
	}
	return res
}
func TestTemp(t *testing.T) {
	//m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	//c := uint32ToRgb(PaletteRGB[0x30])
	//draw.Draw(m, m.Bounds(), &image.Uniform{C: c}, image.Point{0, 0}, draw.Src)
	//Save2jpeg("temp", m)
}
func TestDrawPatternTable(t *testing.T) {
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			m.Set(i, j, uint32ToRgb(uint32(i&j*10)))
		}
	}
	Save2jpeg("temp", m)
	OpenImageFile("temp")
}

func OpenImageFile(filename string) {
	if !strings.HasSuffix(filename, ".jpeg") {
		filename += ".jpeg"
	}
	fmt.Println(exec.Command("open", filename).Run())
}

func TestDrawPalette(t *testing.T) {
	//GenPalette(AllColor)
	//OpenImageFile("palette.jpeg")
}

func TestPPU(t *testing.T) {

}
