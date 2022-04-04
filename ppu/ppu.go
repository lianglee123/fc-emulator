package ppu

import (
	"fc-emulator/rom"
	"fc-emulator/utils"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"image"
	"image/color"
	"image/draw"
	"strconv"
)

type PPU interface {
	ReadForCPU(addr uint16) byte
	WriteForCPU(addr uint16, val byte)
	SetOAM(values []byte)
	Render() image.Image
	CanInterrupt() bool
	EnterVblank()
}

type PPUImpl struct {
	Rom            *rom.NesRom
	Memo           *DefaultPPUMemo
	Register       *RegisterManager
	readDataBuffer byte
	OAM            [256]byte
}

func NewPPU(rom *rom.NesRom) PPU {
	memo := NewPPUMemo(rom)
	return &PPUImpl{
		Memo:     memo,
		Register: NewRegisterManager(),
	}
}

func (p *PPUImpl) DrawBGPatternTable() {
	p.Register.PPUCTRL.BackgroundPatternTableAddress()
}

func (p *PPUImpl) DrawSpritePatternTable() {

}
func (p *PPUImpl) CanInterrupt() bool {
	return p.Register.PPUCTRL.CanGenerateNMIBreakAtStartOfVerticalBlankingInterval()
}

func (p *PPUImpl) ReadForCPU(addr uint16) byte {
	if addr == 0x4014 {
		return p.Register.OAMDMA
	}
	addr = 0x2000 + addr&0b111
	switch addr {
	case 0x2000:
		return byte(p.Register.PPUCTRL)
	case 0x2001:
		return p.Register.PPUMASK
	case 0x2002:
		v := p.Register.PPUSTATUS
		p.Register.PPUSTATUS = v & 0b01111111
		return v
	case 0x2003:
		return p.Register.OAMADDR
	case 0x2004:
		return p.OAM[p.Register.OAMADDR]
	case 0x2005:
		panic("only for write")
	case 0x2006:
		panic("only for write")
	case 0x2007:
		return p.readData()
	default:
		panic("wrong ppu register for cpu")
	}
}
func (p *PPUImpl) EnterVblank() {
	v := p.Register.PPUSTATUS
	v |= 0b10000000 // set 「v」 flag
	v ^= 0b01000000 // toggle 「s」 flag
	p.Register.PPUSTATUS = v
}

func (p *PPUImpl) incrementPPUADDR() {
	p.Register.PPUADDR.Add(p.Register.PPUCTRL.VRAMAddressIncrementPerCPUReadOrWriteOfPPUDATA())
}

func (p *PPUImpl) readData() byte {
	addr := p.Register.PPUADDR.Value()
	p.incrementPPUADDR()
	if addr <= 0x2FFF {
		res := p.readDataBuffer
		p.readDataBuffer = p.Memo.Read(addr)
		return res
	} else if addr >= 0x3000 && addr <= 0x3EFF {
		panic("should not be used in reality")
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		// https://www.nesdev.org/wiki/PPU_palettes
		if addr >= 0x3F20 {
			addr = addr%0x20 + 0x3F00
		}
		if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f1c {
			addr = addr - 0x10
		}
		return p.Memo.Read(addr)
	} else {
		panic("unexpected access to mirrored space")
	}

}
func (p *PPUImpl) writeData(value byte) {
	addr := p.Register.PPUADDR.Value()
	p.incrementPPUADDR()
	if addr <= 0x1FFF {
		panic("attempt to write chr rom space")
	} else if addr >= 0x2000 && addr <= 0x2FFF {
		p.Memo.Write(addr, value)
	} else if addr >= 0x3000 && addr <= 0x3EFF {
		panic("not used in reality")
	} else if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f1c {
		addr = addr - 0x10
		p.Memo.Write(addr, value)
	} else if addr >= 0x3F00 && addr <= 0x3F1F {
		p.Memo.Write(addr, value)
	} else {
		panic(fmt.Sprintf("unexpected ppu addr 0x%X", addr))
	}
}

func (p *PPUImpl) WriteForCPU(addr uint16, val byte) {
	if addr == 0x4014 {
		p.Register.OAMDMA = val // DMA 4014直写，暂时不实现
		panic("DMA used")
		return
	}
	addr = 0x2000 + addr&0b111
	switch addr {
	case 0x2000:
		fmt.Printf("write ppu  contrl %s \n", strconv.FormatUint(uint64(val), 2))
		p.Register.PPUCTRL = PPUCTRL(val)
	case 0x2001:
		p.Register.PPUMASK = val
	case 0x2002:
		p.Register.PPUSTATUS = val
	case 0x2003:
		p.Register.OAMADDR = val
	case 0x2004:
		p.OAM[p.Register.OAMADDR] = val
		p.Register.OAMADDR += 1
	case 0x2005:
		p.Register.PPUSCROLL.Write(val)
	case 0x2006:
		p.Register.PPUADDR.Write(val)
	case 0x2007:
		p.writeData(val)
	default:
		panic("wrong ppu register for cpu")
	}
}

func (p *PPUImpl) SetOAM(values []byte) {
	for i := 0; i < 256; i++ {
		p.OAM[i] = values[i]
	}
}

func (p *PPUImpl) Render() image.Image {
	return p.renderBg()
}

func (p *PPUImpl) DrawBGPalette() image.Image {
	return p.BgPalette().Draw()
}

func (p *PPUImpl) DrawSpritePalette() image.Image {
	return p.SpritePalette().Draw()
}

func (p *PPUImpl) bgPatternTable() []byte {
	startAddr := p.Register.PPUCTRL.BackgroundPatternTableAddress()
	endAddr := startAddr + uint16(4*utils.Kb)
	return p.Memo.Data[startAddr:endAddr]
}
func (p *PPUImpl) renderBg() image.Image {
	patternTable := p.bgPatternTable()
	bgPalette := p.BgPalette()
	nameTable := p.nameTable()
	attributeTable := p.attributeTable()
	return GenBg(nameTable, attributeTable, patternTable, bgPalette)
}

func (p *PPUImpl) BgPalette() Palette {
	data := p.Memo.Data[0x3F00 : 0x3F00+0x10]
	_data := [16]byte{}
	for i := 0; i < 16; i++ {
		_data[i] = data[i]
	}
	return NewPalette(_data)
}

func (p *PPUImpl) SpritePalette() Palette {
	data := p.Memo.Data[0x3F00+0x10 : 0x3F00+0x20]
	_data := [16]byte{}
	for i := 0; i < 16; i++ {
		_data[i] = data[i]
	}
	_data[0] = p.Memo.Data[0x3F00]
	return NewPalette(_data)
}

func (p *PPUImpl) nameTable() []byte {
	baseAddr := p.Register.PPUCTRL.NameTableBaseAddress()
	return p.Memo.Data[baseAddr : baseAddr+960]
}
func (p *PPUImpl) attributeTable() []byte {
	baseAddr := p.Register.PPUCTRL.NameTableBaseAddress()
	return p.Memo.Data[baseAddr+960 : baseAddr+1024]
}

func tileImage(tileRgb *[8][8]color.RGBA) *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			m.Set(x, y, tileRgb[y][x])
		}
	}
	return m
}

type Color uint32

func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32((c >> 4) & 0xFF)
	g = uint32((c >> 2) & 0xFF)
	b = uint32(c & 0xFF)
	a = 0xFFFF
	return
}

func ShowPicture(imageFileName string) {
	myApp := app.New()
	w := myApp.NewWindow("Image")

	//image := canvas.NewImageFromResource(theme.FyneLogo())
	// image := canvas.NewImageFromURI(uri)
	// image := canvas.NewImageFromImage(src)
	// image := canvas.NewImageFromReader(reader, name)
	_image := canvas.NewImageFromFile(imageFileName)
	_image.FillMode = canvas.ImageFillOriginal
	w.SetContent(_image)
	w.ShowAndRun()
}

func uint32ToRgb(v uint32) color.RGBA {
	return color.RGBA{
		R: uint8((v >> 16) & 0xFF),
		G: uint8((v >> 8) & 0xFF),
		B: uint8(v & 0xFF),
		A: 255,
	}
}

func tileIndex2PixelPoint(tileIndex int) Point {
	return Point{X: (tileIndex % 32) * 8, Y: (tileIndex / 32) * 8}
}

func tileColorData2RGBData(tileColorData *[8][8]byte, palette []byte) [8][8]color.RGBA {
	res := [8][8]color.RGBA{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			c := tileColorData[i][j]
			u := AllColor[palette[c]]
			res[i][j] = uint32ToRgb(u)
		}
	}
	return res
}

func tileColorData2RGBData2(tileColorData *[8][8]byte, palette Palette) [8][8]color.RGBA {
	res := [8][8]color.RGBA{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			c := tileColorData[i][j]
			res[i][j] = palette.Color(c)
		}
	}
	return res
}

func appendTile2Screen(s draw.Image, tileIndex int, tileImage *image.RGBA) {

	//for i := 0; i < 8; i++ {
	//	for j := 0; j < 8; j++ {
	//		s.Set( j, i, tileRGBData[i][j])
	//	}
	//}
	pixelPoint := tileIndex2PixelPoint(tileIndex)
	draw.Draw(s, image.Rect(pixelPoint.X, pixelPoint.Y, pixelPoint.X+8, pixelPoint.Y+8),
		tileImage, image.Point{}, draw.Src)
	//for y := pixelPoint.Y; y < pixelPoint.Y; y++ {
	//	for x := pixelPoint.X; x < pixelPoint.X+8; x++ {
	//		s.Set(x, y, tileRGBData[y][x])
	//	}
	//}
}

func mergeRgbDataWithColorMaskByte(rgbData *[8][8]byte, colorMaskByte byte) {
	for i := 0; i < len(rgbData); i++ {
		for j := 0; j < len(rgbData[i]); j++ {
			//fmt.Println(rgbData[i][j], colorMaskByte, rgbData[i][j]|colorMaskByte)
			rgbData[i][j] |= colorMaskByte
		}
	}
}
func getColorMaskDataFromAttributeTable(tileIndex int, attributeTable []byte) byte {
	//return 0b00  << 2
	//return 0b01  << 2
	return 0b10 << 2
	//return 0b11  << 2
	//attributeIndex := tileIndex2attributeIndex(tileIndex)
	//part := tileIndex2attributePart(tileIndex)
	//return getPixelDataFromAttributeTable(attributeIndex, part, attributeTable)
}

func getPixelDataFromAttributeTable(attributeIndex int, blockPart AttributeBlockPart, attributeTable []byte) byte {
	value := attributeTable[attributeIndex]
	switch blockPart {
	case TOP_LEFT_PART:
		return (value & 0b00000011) << 2
	case TOP_RIGHT_PART:
		return value & 0b00001100
	case BOTTOM_LEFT_PART:
		return (value & 0b00110000) >> 2
	case BOTTOM_RIGHT_PART:
		return (value & 0b11000000) >> 4
	default:
		panic("wrong attribute block part")
	}
}

func getColorDataFromPatternTable(patternIndex int, patternTable []byte) [8][8]byte {
	//fmt.Println("****************************")
	patternData := patternTable[patternIndex*16 : patternIndex*16+16]
	bit0Data := patternData[:8]
	//fmt.Println("bit0: ", bit0Data[0])
	//fmt.Println("###bit0Data####")
	//printInBit(bit0Data)
	bit1Data := patternData[8:16]
	res := [8][8]byte{}
	for i := 0; i < 8; i++ {
		bit0Byte := bit0Data[i]
		bit1Byte := bit1Data[i]
		for j := 0; j < 8; j++ {
			bit0 := utils.GetBitFromLeft(bit0Byte, j)
			bit1 := utils.GetBitFromLeft(bit1Byte, j)
			res[i][j] = bit0 | (bit1 << 1)
		}
	}
	//printInGrid(res)
	return res
}

func printInGrid(res [8][8]byte) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Printf("%X  ", res[i][j])
		}
		fmt.Println("")
	}
}
func printInBit(data []byte) {
	for i := 0; i < len(data); i++ {
		s := fmt.Sprintf("%08b\n", data[i])
		for j := 0; j < 8; j++ {
			fmt.Printf("%s  ", s[j:j+1])
		}
		fmt.Println("")
	}

}

func tileIndex2Point(tileIndex int) Point {
	y := tileIndex / 32
	x := tileIndex % 32
	return Point{x, y}
}

func tilePoint2Index(p Point) int {
	return p.X + 32*p.Y
}

func attributeIndex2Point(attributeIndex int) Point {
	return Point{
		X: attributeIndex % 8,
		Y: attributeIndex / 8,
	}
}

func attributePoint2Index(p Point) int {
	return p.X + p.Y*8
}

func tileIndex2attributeIndex(tileIndex int) int {
	tilePoint := tileIndex2Point(tileIndex)
	attributePoint := tilePoint2attributePoint(tilePoint)
	return attributePoint2Index(attributePoint)
}

func tileIndex2attributePart(tileIndex int) AttributeBlockPart {
	tilePoint := tileIndex2Point(tileIndex)
	return getAttributeBlockPart(tilePoint)
}

func tilePoint2attributePoint(tilePoint Point) Point {
	return Point{tilePoint.X / 4, tilePoint.Y / 4}
}

type Point struct {
	X int
	Y int
}

func getAttributeBlockPart(tilePoint Point) AttributeBlockPart {
	isLeft := false
	xM := tilePoint.X % 4
	if xM > 0 && xM <= 2 {
		isLeft = true
	}
	isTop := false
	yM := tilePoint.Y % 4
	if yM > 0 && yM <= 2 {
		isTop = true
	}
	if isTop {
		if isLeft {
			return TOP_LEFT_PART
		} else {
			return TOP_RIGHT_PART
		}
	} else {
		if isLeft {
			return BOTTOM_LEFT_PART
		} else {
			return BOTTOM_RIGHT_PART
		}
	}
}

type AttributeBlockPart int

const (
	TOP_LEFT_PART     AttributeBlockPart = 1
	TOP_RIGHT_PART    AttributeBlockPart = 2
	BOTTOM_LEFT_PART  AttributeBlockPart = 3
	BOTTOM_RIGHT_PART AttributeBlockPart = 4
)

// value = (bottomright << 6) | (bottomleft << 4) | (topright << 2) | (topleft << 0)
// https://www.nesdev.org/wiki/PPU_attribute_tables
