package cpu

import (
	"errors"
	"fc-emulator/cpu/addressing"
	"fc-emulator/cpu/opcode"
	"fc-emulator/memo"
	"fmt"
)

type Bus interface {
	Tick(int)
}

type DefaultBus struct {
}

func (bus *DefaultBus) Tick(n int) {

}

type CPU struct {
	memo     memo.Memo
	register *Register
	debug    bool
	bus      Bus
}

func NewCPU(memo memo.Memo, debug bool) *CPU {
	return &CPU{
		memo: memo,
		register: &Register{
			PC: 0,
			S:  0,
			P:  0,
			A:  0,
			X:  0,
			Y:  0,
		},
		debug: debug,
		bus:   &DefaultBus{},
	}
}

// http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-reset-stack-push-3
func (c *CPU) Reset() {
	c.register.A = 0
	c.register.X = 0
	c.register.Y = 0
	c.register.P = 0x34 // IRQ disabled
	c.register.S = 0xFD
	c.register.PC = c.memo.ReadWord(IV_RESET)
	c.memo.Write(0x4017, 0x00) // frame irq enabled
	c.memo.Write(0x4015, 0x00) // all channels enabled
}

func (c *CPU) increasePC() {
	c.register.IncreasePC()
}

type TraceLog struct {
	OldReg       Register
	NewReg       Register
	opcodeNumber byte
	Code         opcode.Code
	Mode         addressing.Mode
	Addr         uint16
}

func (c *CPU) ExecuteOneInstruction() (*TraceLog, error) {
	if c.debug {
		return c.ExecuteOneInstructionInDebug()
	} else {
		err := c.ExecuteOneInstructionNoDebug()
		return nil, err
	}
}

func (c *CPU) ExecuteOneInstructionInDebug() (*TraceLog, error) {
	opcodeNumber := c.memo.Read(c.register.PC)
	traceLog := &TraceLog{OldReg: *c.register, opcodeNumber: opcodeNumber}
	c.increasePC()

	instruction := instructionTable[opcodeNumber]
	if instruction == nil {
		return nil, errors.New(fmt.Sprintf("opcode 0x%02X is not support", opcodeNumber))
	}
	traceLog.Mode = instruction.Mode
	traceLog.Code = instruction.Code
	addr, crossPage := c.Addressing(instruction.Mode)
	traceLog.Addr = addr
	instruction.Handle(c, addr)
	c.bus.Tick(instruction.Cycle)
	if instruction.CheckPageCross && crossPage {
		c.bus.Tick(1)
	}
	traceLog.NewReg = *c.register
	return traceLog, nil
}

func (c *CPU) ExecuteOneInstructionNoDebug() error {
	opcodeNumber := c.memo.Read(c.register.PC)
	c.increasePC()
	instruction := instructionTable[opcodeNumber]
	if instruction == nil {
		panic(fmt.Sprintf("opcode 0x%02X is not support", opcodeNumber))
		return nil
	}
	addr, crossPage := c.Addressing(instruction.Mode) // 寻址，读参数
	instruction.Handle(c, addr)
	c.bus.Tick(instruction.Cycle)
	if instruction.CheckPageCross && crossPage {
		c.bus.Tick(1)
	}
	return nil
}

func (c *CPU) reset() {
	addr := c.memo.ReadWord(IV_RESET)
	c.register.PC = addr
}
