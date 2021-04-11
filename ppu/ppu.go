package ppu

type PPU interface {
	ReadForCPU(addr uint16) byte
	WriteForCPU(addr uint16, val byte)
	SetOAM(values []byte)
}
