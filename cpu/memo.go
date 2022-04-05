package cpu

import (
	"fc-emulator/pad"
	"fc-emulator/ppu"
	"fc-emulator/rom"
	"fc-emulator/utils"
	"fmt"
)

type Memo interface {
	Read(addr uint16) byte
	ReadWord(addr uint16) uint16
	Write(addr uint16, val byte)
}

type DefaultMemo struct {
	Ram     [2 * utils.Kb]byte
	Trainer []byte
	PrgRom  []byte
	ppu     ppu.PPU
	pad1    pad.Pad
	pad2    pad.Pad
}

func NewMemo(rom *rom.NesRom, _ppu ppu.PPU, pad1, pad2 pad.Pad) Memo {
	prgRom := rom.PrgRom
	if len(rom.PrgRom) == 16*utils.Kb { // Mapper0 Prg Mirror
		prgRom = append(rom.PrgRom, rom.PrgRom...)
	}
	memo := &DefaultMemo{
		Ram:     [2 * utils.Kb]byte{},
		Trainer: rom.Trainer,
		PrgRom:  prgRom,
		ppu:     _ppu,
		pad1:    pad1,
		pad2:    pad2,
	}
	return memo
}

func (m *DefaultMemo) Read(addr uint16) byte {
	addr = m.handleMirror(addr)
	if between(addr, 0, 0x07FF) { // 2k RAM
		return m.Ram[addr]
	} else if between(addr, 0x2000, 0x3FFF) { // ppu register
		return m.ppu.ReadForCPU(addr)
	} else if between(addr, 0x4000, 0x4013) { // some io register
		return 0
	} else if addr == 0x4014 {
		return m.ppu.ReadForCPU(addr)
	} else if addr == 0x4016 {
		return m.pad1.ReadForCPU()
	} else if addr == 0x4017 {
		return m.pad2.ReadForCPU()
	} else if between(addr, 0x4015, 0x5fff) {
		// some io register and expansion Rom
		return 0
	} else if between(addr, 0x7000, 0x71FF) {
		// Trainer, SRAM, 带电池的RAM
		return m.Trainer[addr]
	} else if between(addr, 0x8000, 0xffff) {
		addr -= 0x8000
		return m.PrgRom[addr]
	} else {
		panic(fmt.Sprintf("Read Wrong Data Addr: %X", addr))
	}
}

func (m *DefaultMemo) ReadWord(addr uint16) uint16 {
	byte1 := m.Read(addr)
	byte2 := m.Read(addr + 1)
	return uint16(byte2)<<8 | uint16(byte1) // 小端序
}

func (m *DefaultMemo) Write(addr uint16, val byte) {
	addr = m.handleMirror(addr)
	if between(addr, 0, 0x07FF) { // 2k RAM
		m.Ram[addr] = val
	} else if between(addr, 0x2000, 0x3FFF) { // ppu register
		m.ppu.WriteForCPU(addr, val)
	} else if between(addr, 0x4000, 0x4013) {
		// io register
	} else if addr == 0x4014 {
		// DMA直写，把整个PAGE的地址写进OAM
		pageNo := uint16(val)
		left := pageNo * 256
		m.ppu.SetOAM(m.copyRam(left, left+256))
	} else if addr == 0x4016 {
		m.pad1.WriteForCPU(val)
	} else if addr == 0x4017 {
		m.pad2.WriteForCPU(val)
	} else if between(addr, 0x4015, 0x5fff) {
		// some io register and expansion Rom
	} else {
		panic(Str("Write Wrong Data Addr", addr, val))
	}
}

func (m *DefaultMemo) copyRam(begin, end uint16) []byte {
	res := make([]byte, 0, end-begin)
	for begin < end {
		res = append(res, m.Ram[begin])
		begin += 1
	}
	return res
}

func (m *DefaultMemo) handleMirror(addr uint16) uint16 {
	// 0x0800: 2k, 0x2000: 8k
	if addr >= 0x0800 && addr <= 0x1FFF {
		addr = addr % 0x0800 // Mirrors of $0000-$07FF， 这里的2k的RAM被Mirror了三次
	} else if addr >= 0x2008 && addr <= 0x3FFF {
		addr = 0x2000 + (addr-0x2008)%8
	}
	return addr
}

func between(val, lower, upper uint16) bool {
	return val >= lower && val <= upper
}

func Str(values ...interface{}) string {
	return fmt.Sprint(values...)
}
