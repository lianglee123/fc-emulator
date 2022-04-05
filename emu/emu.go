package emu

import (
	"fc-emulator/cpu"
	"fc-emulator/pad"
	"fc-emulator/ppu"
	"fc-emulator/rom"
	"time"
)

type Emu struct {
	CPU           *cpu.CPU
	PPU           ppu.PPU
	Opt           *EmuOpt
	Rom           *rom.NesRom
	Pad1          pad.Pad
	Pad2          pad.Pad
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
	pad1 := pad.NewPad()
	pad2 := pad.NewPad()
	cpuMemo := cpu.NewMemo(nesRom, _ppu, pad1, pad2)
	c := cpu.NewCPU(cpuMemo, e.Opt.Debug)
	c.Reset()
	e.CPU = c
	e.Pad1 = pad1
	e.Pad2 = pad2
	return nil
}

func (e *Emu) Start() {
	cnt := 0
	for {
		e.PPU.EnterVblank()
		if e.PPU.CanInterrupt() {
			e.CPU.ExecNMI()
		}
		for i := 0; i < 1000; i++ {
			_, err := e.CPU.ExecuteOneInstruction()
			if err != nil {
				panic(err)
			}
			cnt += 1
			//if cnt%1000 == 0 {
			//	//fmt.Println("Count ", cnt)
			//}
		}
		if e.FrameCallback != nil {
			e.FrameCallback()
		}
		time.Sleep(30 * time.Millisecond)
	}
}
