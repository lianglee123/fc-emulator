package emu

import (
	"fc-emulator/cpu"
	"fc-emulator/ppu"
	"fc-emulator/rom"
	"fmt"
)

type Emu struct {
	CPU           *cpu.CPU
	PPU           ppu.PPU
	Opt           *EmuOpt
	Rom           *rom.NesRom
	FrameCallback func()
}

type EmuOpt struct {
	Debug bool
}

func NewEmu(opt *EmuOpt) *Emu {
	if opt == nil {
		opt = &EmuOpt{Debug: false}
	}
	return &Emu{Opt: opt}
}

func (e *Emu) Load(fileName string) error {
	nesRom, err := rom.LoadNesRom(fileName)
	e.Rom = nesRom
	_ppu := ppu.NewPPU(nesRom)
	if err != nil {
		return err
	}
	e.PPU = _ppu
	cpuMemo := cpu.NewMemo(nesRom, _ppu)
	c := cpu.NewCPU(cpuMemo, e.Opt.Debug)
	c.Reset()
	e.CPU = c
	return nil
}

func (e *Emu) Start() {
	cnt := 0
	for {
		e.PPU.EnterVblank()
		if e.PPU.CanInterrupt() {
			e.CPU.ExecNMI()
		}
		for i := 0; i < 5000; i++ {
			_, err := e.CPU.ExecuteOneInstruction()
			if err != nil {
				panic(err)
			}
			cnt += 1
			if cnt%10000 == 0 {
				fmt.Println("Count ", cnt)
			}
		}
		if e.FrameCallback != nil {
			e.FrameCallback()
		}
	}
}
