package ppu

import (
	"fc-emulator/rom"
)

type PPUMemo interface {
	Read(addr uint16) byte
	ReadWord(addr uint16) uint16
	Write(addr uint16, val byte)
}

type DefaultPPUMemo struct {
	Data []byte
}

func NewPPUMemo(rom *rom.NesRom) *DefaultPPUMemo {
	data := make([]byte, 0x10000)
	for i, b := range rom.ChrRom {
		data[i] = b
	}
	return &DefaultPPUMemo{
		Data: data,
	}
}

func (m *DefaultPPUMemo) Read(addr uint16) byte {
	if addr >= 0x4000 {
		addr = addr % 0x4000
	}
	return m.Data[addr]
}

func (m *DefaultPPUMemo) Write(addr uint16, val byte) {
	if addr >= 0x4000 {
		addr = addr % 0x4000
	}
	m.Data[addr] = val
}

func (m *DefaultPPUMemo) ReadWord(addr uint16) uint16 {
	panic("implement me")
}

//$3F00	Universal background color
//$3F01-$3F03	Background palette 0
//$3F05-$3F07	Background palette 1
//$3F09-$3F0B	Background palette 2
//$3F0D-$3F0F	Background palette 3
