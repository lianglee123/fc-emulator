package emu

import (
	"fc-emulator/cpu"
	"fc-emulator/rom"
)

type Emu struct {
	cpu *cpu.CPU
}

func NewEmu() *Emu {
	return &Emu{}
}

func (e *Emu) Load(fileName string) error {
	nesRom, err := rom.LoadNesRom(fileName)
	if err != nil {
		return err
	}
	cpuMemo := cpu.NewMemo(nesRom)
	c := cpu.NewCPU(cpuMemo, true)
	c.Reset()
	e.cpu = c
	return nil
}

func (e *Emu) Start() {
	for {
		e.cpu.ExecuteOneInstruction()
	}
}
