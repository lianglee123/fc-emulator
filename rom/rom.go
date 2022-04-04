package rom

import (
	"errors"
	"fc-emulator/utils"
	"fmt"
	"io/ioutil"
)

const nesPrefix = "NES\x1A"

// https://www.nesdev.org/wiki/INES
type NesRom struct {
	MapperNumber int
	IsTrainer    bool
	HasBattery   bool
	PrgRom       []byte
	ChrRom       []byte
	Trainer      []byte
	Header       *Header
}

type NesFlag1 struct {
	data                   byte
	Mirroring              bool
	HasBatterySRAM         bool
	HasTrainer             bool
	ScreenMode             bool
	MapperVersionLower4Bit byte
}

type NesFlag2 struct {
	data                  byte
	VSUnisystem           bool // not used, no need to understand
	PlayChoice10          bool // not used, no need to understand
	NesVersion            byte // if is 2, meaning NES2.0
	MapperVersionHigh4Bit byte
}

type Header struct {
	Data         []byte
	PrgCount     int // unit is 16k
	ChrCount     int //  unit is 8k
	Flag6        *Flag6
	Flag7        *Flag7
	MapperNumber byte
}

func (h *Header) String() string {
	return fmt.Sprintf("MapperNumber: %d, prgCount: %d, chrCount: %d, Flag6: %s",
		h.MapperNumber, h.PrgCount, h.ChrCount, h.Flag6.String())
}

func (r *NesRom) String() string {
	return fmt.Sprintf("Header: %s, PrgRomSize(kb): %d, ChrRomSize(kb): %d ", r.Header.String(),
		len(r.PrgRom)/utils.Kb, len(r.ChrRom)/utils.Kb)
}

func NewHeader(data []byte) *Header {
	if len(data) != 16 {
		panic("header must has 16 byte")
	}
	flag6 := parseFlag6(data[6])
	flag7 := parseFlag7(data[7])
	mapperNumber := (flag7.MapperHighVersion << 4) | flag6.MapperLowerVersion
	return &Header{
		Data:         data,
		PrgCount:     int(data[4]),
		ChrCount:     int(data[5]),
		Flag6:        parseFlag6(data[6]),
		Flag7:        parseFlag7(data[7]),
		MapperNumber: mapperNumber,
	}
}

type NameTableMirrorMode int

const (
	VerticalMirror   NameTableMirrorMode = 1
	HorizontalMirror NameTableMirrorMode = 2
)

var MirrorModeNameMap map[NameTableMirrorMode]string = map[NameTableMirrorMode]string{
	VerticalMirror:   "VerticalMirror",
	HorizontalMirror: "HorizontalMirror",
}

func LoadNesRom(filename string) (*NesRom, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if string(data[:4]) != nesPrefix {
		return nil, errors.New("nes format error")
	}
	if data[7]&0x0c == 0x08 {
		return nil, errors.New("NES2.0 is not currently supported")
	}
	rom := &NesRom{
		Header:     NewHeader(data[:16]),
		IsTrainer:  (data[6] & 0x04) != 0x00, // 第4bit
		HasBattery: (data[6] & 0x02) != 0x00, // 第2bit
	}
	if rom.Header.MapperNumber != 0 {
		return nil, utils.NewError("Not support Mapper ", rom.MapperNumber)
	}
	prgStartIndex := 16
	if rom.IsTrainer {
		rom.Trainer = data[16:512]
		prgStartIndex += 512
	} else {
		rom.Trainer = make([]byte, 512)
	}
	prgSize := rom.Header.PrgCount * 16 * utils.Kb
	rom.PrgRom = data[prgStartIndex : prgStartIndex+prgSize]
	rom.ChrRom = data[prgStartIndex+prgSize:]
	return rom, nil
}

//76543210
//||||||||
//|||||||+- Mirroring: 0: 水平镜像（PPU 章节再介绍）
//|||||||              1: 垂直镜像（PPU 章节再介绍）
//||||||+-- 1: 卡带上有没有带电池的 SRAM
//|||||+--- 1: Trainer 标志
//||||+---- 1: 4-Screen 模式（PPU 章节再介绍）
//++++----- Mapper 号的低 4 bit

type Flag6 struct {
	HasBattery         bool
	MirrorMode         NameTableMirrorMode
	Trainer            bool
	FourScreenMode     bool
	MapperLowerVersion byte
}

func (f *Flag6) String() string {
	tmp := "HasBattery: %v, MirrorMode: %s, Trainer: %v, FourScreenMode: %v, MapperLowerVersion:  %v"
	return fmt.Sprintf(tmp, f.HasBattery, MirrorModeNameMap[f.MirrorMode], f.Trainer, f.FourScreenMode, f.MapperLowerVersion)
}

func parseFlag6(v byte) *Flag6 {
	return &Flag6{
		MirrorMode:         getMirrorMode(v),
		HasBattery:         utils.GetBitFromRight(v, 1) == 1,
		Trainer:            utils.GetBitFromRight(v, 2) == 1,
		FourScreenMode:     utils.GetBitFromRight(v, 3) == 1,
		MapperLowerVersion: v >> 4,
	}
}

func getMirrorMode(v byte) NameTableMirrorMode {
	if utils.GetBitFromRight(v, 0) == 0 {
		return HorizontalMirror
	} else {
		return VerticalMirror
	}

}

type Flag7 struct {
	VSUnisystem       bool
	PlayChoice10      bool
	NesFormat         int
	MapperHighVersion byte
}

//76543210
//||||||||
//|||||||+- VS Unisystem，不需要了解
//||||||+-- PlayChoice-10，不需要了解
//||||++--- 如果为 2，代表 NES 2.0 格式，不需要了解
//++++----- Mapper 号的高 4 bit
func parseFlag7(v byte) *Flag7 {
	return &Flag7{
		MapperHighVersion: v >> 4,
	}
}
