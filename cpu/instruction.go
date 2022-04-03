package cpu

import (
	"fc-emulator/cpu/addressing"
	"fc-emulator/cpu/opcode"
)

type Instruction struct {
	Code   opcode.Code
	Mode   addressing.Mode
	Handle func(cpu *CPU, addr uint16)
}

func codeCount() int {
	var l = 0
	for _, v := range instructionTable {
		if v != nil {
			l++
		}

	}
	return l
}

// http://obelisk.me.uk/6502/reference.html#JMP
// https://www.masswerk.at/6502/6502_instruction_set.html
// https://www.pagetable.com/c64ref/6502/
var instructionTable = [256]*Instruction{

	// Sets the program counter to the address specified by the operand.
	0x4C: {opcode.JMP, addressing.ABS, (*CPU).JMP},
	0x6C: {opcode.JMP, addressing.IND, (*CPU).JMP},

	// Loads a byte of memory into the X register setting the zero and negative flags as appropriate.
	0xA2: {opcode.LDX, addressing.IMM, (*CPU).LDX},
	0xA6: {opcode.LDX, addressing.ZPG, (*CPU).LDX},
	0xB6: {opcode.LDX, addressing.ZPY, (*CPU).LDX},
	0xAE: {opcode.LDX, addressing.ABS, (*CPU).LDX},
	0xBE: {opcode.LDX, addressing.ABY, (*CPU).LDX},

	// Stores the contents of the X register into memory.
	0x86: {opcode.STX, addressing.ZPG, (*CPU).STX},
	0x96: {opcode.STX, addressing.ZPY, (*CPU).STX},
	0x8E: {opcode.STX, addressing.ABS, (*CPU).STX},

	0x84: {opcode.STY, addressing.ZPG, (*CPU).STY},
	0x94: {opcode.STY, addressing.ZPX, (*CPU).STY},
	0x8C: {opcode.STY, addressing.ABS, (*CPU).STY},

	// Stores the contents of the accumulator into memory.
	0x85: {opcode.STA, addressing.ZPG, (*CPU).STA},
	0x95: {opcode.STA, addressing.ZPX, (*CPU).STA},
	0x8D: {opcode.STA, addressing.ABS, (*CPU).STA},
	0x9D: {opcode.STA, addressing.ABX, (*CPU).STA},
	0x99: {opcode.STA, addressing.ABY, (*CPU).STA},
	0x81: {opcode.STA, addressing.INX, (*CPU).STA},
	0x91: {opcode.STA, addressing.INY, (*CPU).STA},

	// The JSR instruction pushes the address (minus one) of the return point on
	// to the stack and then sets the program counter to the target memory address.
	0x20: {opcode.JSR, addressing.ABS, (*CPU).JSR},

	// Set the carry flag to one.
	0x38: {opcode.SEC, addressing.IMP, (*CPU).SEC},

	// If the carry flag is set then add the relative displacement to
	// the program counter to cause a branch to a new location.
	0xB0: {opcode.BCS, addressing.REL, (*CPU).BCS},

	// Set the carry flag to zero.
	0x18: {opcode.CLC, addressing.IMP, (*CPU).CLC},

	// If the carry flag is clear then add the relative displacement to the program
	// counter to cause a branch to a new location.
	0x90: {opcode.BCC, addressing.REL, (*CPU).BCC},

	// Loads a byte of memory into the accumulator setting the zero and
	// negative flags as appropriate.
	0xA9: {opcode.LDA, addressing.IMM, (*CPU).LDA},
	0xA5: {opcode.LDA, addressing.ZPG, (*CPU).LDA},
	0xB5: {opcode.LDA, addressing.ZPX, (*CPU).LDA},
	0xAD: {opcode.LDA, addressing.ABS, (*CPU).LDA},
	0xBD: {opcode.LDA, addressing.ABX, (*CPU).LDA},
	0xB9: {opcode.LDA, addressing.ABY, (*CPU).LDA},
	0xA1: {opcode.LDA, addressing.INX, (*CPU).LDA},
	0xB1: {opcode.LDA, addressing.INY, (*CPU).LDA},

	// If the zero flag is set then add the relative displacement
	// to the program counter to cause a branch to a new location.
	0xF0: {opcode.BEQ, addressing.REL, (*CPU).BEQ},

	// If the zero flag is clear then add the relative displacement to
	// the program counter to cause a branch to a new location.
	0xD0: {opcode.BNE, addressing.REL, (*CPU).BNE},

	// This instructions is used to test if one or more bits are set in a target memory location.
	// The mask pattern in A is ANDed with the value in memory to set or clear the zero flag,
	// but the result is not kept.
	// Bits 7 and 6 of the value from memory are copied into the N and V flags.
	0x24: {opcode.BIT, addressing.ZPG, (*CPU).BIT},
	0x2C: {opcode.BIT, addressing.ABS, (*CPU).BIT},

	// If the overflow flag is set then add the relative displacement to
	// the program counter to cause a branch to a new location.
	0x70: {opcode.BVS, addressing.REL, (*CPU).BVS},

	// If the overflow flag is clear then add the relative displacement
	// to the program counter to cause a branch to a new location.
	0x50: {opcode.BVC, addressing.REL, (*CPU).BVC},

	// If the negative flag is clear then add the relative displacement
	// to the program counter to cause a branch to a new location.
	0x10: {opcode.BPL, addressing.REL, (*CPU).BPL},

	// The RTS instruction is used at the end of a subroutine to return
	// to the calling routine. It pulls the program counter (minus one) from the stack.
	0x60: {opcode.RTS, addressing.IMP, (*CPU).RTS},

	// Set the interrupt disable flag to one.
	0x78: {opcode.SEI, addressing.IMP, (*CPU).SEI},

	// Set the decimal mode flag to one.
	0xF8: {opcode.SED, addressing.IMP, (*CPU).SED},

	// Pushes a copy of the status flags on to the stack.
	0x08: {opcode.PHP, addressing.IMP, (*CPU).PHP},

	// Pulls an 8 bit value from the stack and into the accumulator.
	// The zero and negative flags are set as appropriate.
	0x68: {opcode.PLA, addressing.IMP, (*CPU).PLA},

	// A logical AND is performed, bit by bit, on the accumulator contents
	// using the contents of a byte of memory.
	0x29: {opcode.AND, addressing.IMM, (*CPU).AND},
	0x25: {opcode.AND, addressing.ZPG, (*CPU).AND},
	0x35: {opcode.AND, addressing.ZPX, (*CPU).AND},
	0x2D: {opcode.AND, addressing.ABS, (*CPU).AND},
	0x3D: {opcode.AND, addressing.ABX, (*CPU).AND},
	0x39: {opcode.AND, addressing.ABY, (*CPU).AND},
	0x21: {opcode.AND, addressing.INX, (*CPU).AND},
	0x31: {opcode.AND, addressing.INY, (*CPU).AND},

	// This instruction compares the contents of the accumulator with another memory
	// held value and sets the zero and carry flags as appropriate.
	0xC9: {opcode.CMP, addressing.IMM, (*CPU).CMP},
	0xC5: {opcode.CMP, addressing.ZPG, (*CPU).CMP},
	0xD5: {opcode.CMP, addressing.ZPX, (*CPU).CMP},
	0xCD: {opcode.CMP, addressing.ABS, (*CPU).CMP},
	0xDD: {opcode.CMP, addressing.ABX, (*CPU).CMP},
	0xD9: {opcode.CMP, addressing.ABY, (*CPU).CMP},
	0xC1: {opcode.CMP, addressing.INX, (*CPU).CMP},
	0xD1: {opcode.CMP, addressing.INY, (*CPU).CMP},

	// Sets the decimal mode flag to zero.
	0xD8: {opcode.CLD, addressing.IMP, (*CPU).CLD},

	// Pushes a copy of the accumulator on to the stack.
	0x48: {opcode.PHA, addressing.IMP, (*CPU).PHA},

	// Pulls an 8 bit value from the stack and into the processor flags.
	// The flags will take on new states as determined by the value pulled.
	0x28: {opcode.PLP, addressing.IMP, (*CPU).PLP},

	// If the negative flag is set then add the relative displacement
	// to the program counter to cause a branch to a new location.
	0x30: {opcode.BMI, addressing.REL, (*CPU).BMI},

	// A,Z,N = A|M
	// An inclusive OR is performed, bit by bit,
	// on the accumulator contents using the contents of a byte of memory.
	0x09: {opcode.ORA, addressing.IMM, (*CPU).ORA},
	0x05: {opcode.ORA, addressing.ZPG, (*CPU).ORA},
	0x15: {opcode.ORA, addressing.ZPX, (*CPU).ORA},
	0x0D: {opcode.ORA, addressing.ABS, (*CPU).ORA},
	0x1D: {opcode.ORA, addressing.ABX, (*CPU).ORA},
	0x19: {opcode.ORA, addressing.ABY, (*CPU).ORA},
	0x01: {opcode.ORA, addressing.INX, (*CPU).ORA},
	0x11: {opcode.ORA, addressing.INY, (*CPU).ORA},

	// Clears the overflow flag.
	0xB8: {opcode.CLV, addressing.IMP, (*CPU).CLV},

	// An exclusive OR is performed, bit by bit,
	// on the accumulator contents using the contents of a byte of memory.
	0x49: {opcode.EOR, addressing.IMM, (*CPU).EOR},
	0x45: {opcode.EOR, addressing.ZPG, (*CPU).EOR},
	0x55: {opcode.EOR, addressing.ZPX, (*CPU).EOR},
	0x4D: {opcode.EOR, addressing.ABS, (*CPU).EOR},
	0x5D: {opcode.EOR, addressing.ABX, (*CPU).EOR},
	0x59: {opcode.EOR, addressing.ABY, (*CPU).EOR},
	0x41: {opcode.EOR, addressing.INX, (*CPU).EOR},
	0x51: {opcode.EOR, addressing.INY, (*CPU).EOR},

	// This instruction adds the contents of a memory location to
	// the accumulator together with the carry bit.
	// If overflow occurs the carry bit is set, this enables multiple byte addition to be performed.
	0x69: {opcode.ADC, addressing.IMM, (*CPU).ADC},
	0x65: {opcode.ADC, addressing.ZPG, (*CPU).ADC},
	0x75: {opcode.ADC, addressing.ZPX, (*CPU).ADC},
	0x6D: {opcode.ADC, addressing.ABS, (*CPU).ADC},
	0x7D: {opcode.ADC, addressing.ABX, (*CPU).ADC},
	0x79: {opcode.ADC, addressing.ABY, (*CPU).ADC},
	0x61: {opcode.ADC, addressing.INX, (*CPU).ADC},
	0x71: {opcode.ADC, addressing.INY, (*CPU).ADC},

	0xA0: {opcode.LDY, addressing.IMM, (*CPU).LDY},
	0xA4: {opcode.LDY, addressing.ZPG, (*CPU).LDY},
	0xB4: {opcode.LDY, addressing.ZPX, (*CPU).LDY},
	0xAC: {opcode.LDY, addressing.ABS, (*CPU).LDY},
	0xBC: {opcode.LDY, addressing.ABX, (*CPU).LDY},

	// This instruction compares the contents of the Y register with another memory
	// held value and sets the zero and carry flags as appropriate.
	0xC0: {opcode.CPY, addressing.IMM, (*CPU).CPY},
	0xC4: {opcode.CPY, addressing.ZPG, (*CPU).CPY},
	0xCC: {opcode.CPY, addressing.ABS, (*CPU).CPY},

	// This instruction compares the contents of the X register with another
	// memory held value and sets the zero and carry flags as appropriate.
	0xE0: {opcode.CPX, addressing.IMM, (*CPU).CPX},
	0xE4: {opcode.CPX, addressing.ZPG, (*CPU).CPX},
	0xEC: {opcode.CPX, addressing.ABS, (*CPU).CPX},

	// This instruction subtracts the contents of a memory location to
	// the accumulator together with the not of the carry
	// bit. If overflow occurs the carry bit is clear, this enables
	// multiple byte subtraction to be performed.
	0xE9: {opcode.SBC, addressing.IMM, (*CPU).SBC},
	0xEB: {opcode.SBC, addressing.IMM, (*CPU).SBC},
	0xE5: {opcode.SBC, addressing.ZPG, (*CPU).SBC},
	0xF5: {opcode.SBC, addressing.ZPX, (*CPU).SBC},
	0xED: {opcode.SBC, addressing.ABS, (*CPU).SBC},
	0xFD: {opcode.SBC, addressing.ABX, (*CPU).SBC},
	0xF9: {opcode.SBC, addressing.ABY, (*CPU).SBC},
	0xE1: {opcode.SBC, addressing.INX, (*CPU).SBC},
	0xF1: {opcode.SBC, addressing.INY, (*CPU).SBC},

	// Adds one to the Y register setting the zero and negative flags as appropriate
	0xC8: {opcode.INY, addressing.IMP, (*CPU).INY},

	// Adds one to the X register setting the zero and negative flags as appropriate.
	0xE8: {opcode.INX, addressing.IMP, (*CPU).INX},

	0x88: {opcode.DEY, addressing.IMP, (*CPU).DEY},
	0xCA: {opcode.DEX, addressing.IMP, (*CPU).DEX},

	0xC6: {opcode.DEC, addressing.ZPG, (*CPU).DEC},
	0xD6: {opcode.DEC, addressing.ZPX, (*CPU).DEC},
	0xCE: {opcode.DEC, addressing.ABS, (*CPU).DEC},
	0xDE: {opcode.DEC, addressing.ABX, (*CPU).DEC},

	0xA8: {opcode.TAY, addressing.IMP, (*CPU).TAY},
	0x98: {opcode.TYA, addressing.IMP, (*CPU).TYA},

	0xAA: {opcode.TAX, addressing.IMP, (*CPU).TAX},
	0x8A: {opcode.TXA, addressing.IMP, (*CPU).TXA},

	0xBA: {opcode.TSX, addressing.IMP, (*CPU).TSX},
	0x9A: {opcode.TXS, addressing.IMP, (*CPU).TXS},

	0x40: {opcode.RTI, addressing.IMP, (*CPU).RTI},

	0x4A: {opcode.LSR, addressing.IMP, (*CPU).LSRImp},
	0x46: {opcode.LSR, addressing.ZPG, (*CPU).LSR},
	0x56: {opcode.LSR, addressing.ZPX, (*CPU).LSR},
	0x4E: {opcode.LSR, addressing.ABS, (*CPU).LSR},
	0x5E: {opcode.LSR, addressing.ABX, (*CPU).LSR},

	0x0A: {opcode.ASL, addressing.IMP, (*CPU).ASLImp},
	0x06: {opcode.ASL, addressing.ZPG, (*CPU).ASL},
	0x16: {opcode.ASL, addressing.ZPX, (*CPU).ASL},
	0x0E: {opcode.ASL, addressing.ABS, (*CPU).ASL},
	0x1E: {opcode.ASL, addressing.ABX, (*CPU).ASL},

	0x6A: {opcode.ROR, addressing.IMP, (*CPU).RORImp},
	0x66: {opcode.ROR, addressing.ZPG, (*CPU).ROR},
	0x76: {opcode.ROR, addressing.ZPX, (*CPU).ROR},
	0x6E: {opcode.ROR, addressing.ABS, (*CPU).ROR},
	0x7E: {opcode.ROR, addressing.ABX, (*CPU).ROR},

	0x2A: {opcode.ROL, addressing.IMP, (*CPU).ROLImp},
	0x26: {opcode.ROL, addressing.ZPG, (*CPU).ROL},
	0x36: {opcode.ROL, addressing.ZPX, (*CPU).ROL},
	0x2E: {opcode.ROL, addressing.ABS, (*CPU).ROL},
	0x3E: {opcode.ROL, addressing.ABX, (*CPU).ROL},

	0xE6: {opcode.INC, addressing.ZPG, (*CPU).INC},
	0xF6: {opcode.INC, addressing.ZPX, (*CPU).INC},
	0xEE: {opcode.INC, addressing.ABS, (*CPU).INC},
	0xFE: {opcode.INC, addressing.ABX, (*CPU).INC},

	0x04: {opcode.NOP, addressing.ZPG, (*CPU).NOP},
	0x0C: {opcode.NOP, addressing.ABS, (*CPU).NOP},
	0x14: {opcode.NOP, addressing.ZPX, (*CPU).NOP},
	0x1A: {opcode.NOP, addressing.IMP, (*CPU).NOP},
	0x1C: {opcode.NOP, addressing.ABX, (*CPU).NOP},
	0x34: {opcode.NOP, addressing.ZPX, (*CPU).NOP},
	0x3A: {opcode.NOP, addressing.IMP, (*CPU).NOP},
	0x3C: {opcode.NOP, addressing.ABX, (*CPU).NOP},
	0x44: {opcode.NOP, addressing.ZPG, (*CPU).NOP},
	0x54: {opcode.NOP, addressing.ZPX, (*CPU).NOP},
	0x5A: {opcode.NOP, addressing.IMP, (*CPU).NOP},
	0x5C: {opcode.NOP, addressing.ABX, (*CPU).NOP},
	0x64: {opcode.NOP, addressing.ZPG, (*CPU).NOP},
	0x74: {opcode.NOP, addressing.ZPX, (*CPU).NOP},
	0x7A: {opcode.NOP, addressing.IMP, (*CPU).NOP},
	0x7C: {opcode.NOP, addressing.ABX, (*CPU).NOP},
	0x80: {opcode.NOP, addressing.IMM, (*CPU).NOP},
	0x82: {opcode.NOP, addressing.IMM, (*CPU).NOP},
	0x89: {opcode.NOP, addressing.IMM, (*CPU).NOP},
	0xC2: {opcode.NOP, addressing.IMM, (*CPU).NOP},
	0xD4: {opcode.NOP, addressing.ZPX, (*CPU).NOP},
	0xDA: {opcode.NOP, addressing.IMP, (*CPU).NOP},
	0xDC: {opcode.NOP, addressing.ABX, (*CPU).NOP},
	0xE2: {opcode.NOP, addressing.IMM, (*CPU).NOP},
	0xEA: {opcode.NOP, addressing.IMP, (*CPU).NOP},
	0xF4: {opcode.NOP, addressing.ZPX, (*CPU).NOP},
	0xFA: {opcode.NOP, addressing.IMP, (*CPU).NOP},
	0xFC: {opcode.NOP, addressing.ABX, (*CPU).NOP},

	0xA3: {opcode.LAX, addressing.INX, (*CPU).LAX},
	0xA7: {opcode.LAX, addressing.ZPG, (*CPU).LAX},
	0xAB: {opcode.LAX, addressing.IMM, (*CPU).LAX},
	0xAF: {opcode.LAX, addressing.ABS, (*CPU).LAX},
	0xB3: {opcode.LAX, addressing.INY, (*CPU).LAX},
	0xB7: {opcode.LAX, addressing.ZPY, (*CPU).LAX},
	0xBF: {opcode.LAX, addressing.ABY, (*CPU).LAX},

	0x83: {opcode.SAX, addressing.INX, (*CPU).SAX},
	0x87: {opcode.SAX, addressing.ZPG, (*CPU).SAX},
	0x8F: {opcode.SAX, addressing.ABS, (*CPU).SAX},
	0x97: {opcode.SAX, addressing.ZPY, (*CPU).SAX},

	0xC3: {opcode.DCP, addressing.INX, (*CPU).DCP},
	0xC7: {opcode.DCP, addressing.ZPG, (*CPU).DCP},
	0xCF: {opcode.DCP, addressing.ABS, (*CPU).DCP},
	0xD3: {opcode.DCP, addressing.INY, (*CPU).DCP},
	0xD7: {opcode.DCP, addressing.ZPX, (*CPU).DCP},
	0xDB: {opcode.DCP, addressing.ABY, (*CPU).DCP},
	0xDF: {opcode.DCP, addressing.ABX, (*CPU).DCP},

	0xE3: {opcode.ISB, addressing.INX, (*CPU).ISB},
	0xE7: {opcode.ISB, addressing.ZPG, (*CPU).ISB},
	0xEF: {opcode.ISB, addressing.ABS, (*CPU).ISB},
	0xF3: {opcode.ISB, addressing.INY, (*CPU).ISB},
	0xF7: {opcode.ISB, addressing.ZPX, (*CPU).ISB},
	0xFB: {opcode.ISB, addressing.ABY, (*CPU).ISB},
	0xFF: {opcode.ISB, addressing.ABX, (*CPU).ISB},

	0x03: {opcode.SLO, addressing.INX, (*CPU).SLO},
	0x07: {opcode.SLO, addressing.ZPG, (*CPU).SLO},
	0x0F: {opcode.SLO, addressing.ABS, (*CPU).SLO},
	0x13: {opcode.SLO, addressing.INY, (*CPU).SLO},
	0x17: {opcode.SLO, addressing.ZPX, (*CPU).SLO},
	0x1B: {opcode.SLO, addressing.ABY, (*CPU).SLO},
	0x1F: {opcode.SLO, addressing.ABX, (*CPU).SLO},

	0x23: {opcode.RLA, addressing.INX, (*CPU).RLA},
	0x27: {opcode.RLA, addressing.ZPG, (*CPU).RLA},
	0x2F: {opcode.RLA, addressing.ABS, (*CPU).RLA},
	0x33: {opcode.RLA, addressing.INY, (*CPU).RLA},
	0x37: {opcode.RLA, addressing.ZPX, (*CPU).RLA},
	0x3B: {opcode.RLA, addressing.ABY, (*CPU).RLA},
	0x3F: {opcode.RLA, addressing.ABX, (*CPU).RLA},

	0x43: {opcode.SRE, addressing.INX, (*CPU).SRE},
	0x47: {opcode.SRE, addressing.ZPG, (*CPU).SRE},
	0x4F: {opcode.SRE, addressing.ABS, (*CPU).SRE},
	0x53: {opcode.SRE, addressing.INY, (*CPU).SRE},
	0x57: {opcode.SRE, addressing.ZPX, (*CPU).SRE},
	0x5B: {opcode.SRE, addressing.ABY, (*CPU).SRE},
	0x5F: {opcode.SRE, addressing.ABX, (*CPU).SRE},

	0x63: {opcode.RRA, addressing.INX, (*CPU).RRA},
	0x67: {opcode.RRA, addressing.ZPG, (*CPU).RRA},
	0x6F: {opcode.RRA, addressing.ABS, (*CPU).RRA},
	0x73: {opcode.RRA, addressing.INY, (*CPU).RRA},
	0x77: {opcode.RRA, addressing.ZPX, (*CPU).RRA},
	0x7B: {opcode.RRA, addressing.ABY, (*CPU).RRA},
	0x7F: {opcode.RRA, addressing.ABX, (*CPU).RRA},
}

func (c *CPU) RRA(addr uint16) {
	c.ROR(addr)
	c.ADC(addr)
}

func (c *CPU) SRE(addr uint16) {
	c.LSR(addr)
	c.EOR(addr)
}
func (c *CPU) RLA(addr uint16) {
	c.ROL(addr)
	c.AND(addr)
}

func (c *CPU) SLO(addr uint16) {
	c.ASL(addr)
	c.ORA(addr)
}

func (c *CPU) ISB(addr uint16) {
	c.INC(addr)
	c.SBC(addr)
}

func (c *CPU) DCP(addr uint16) {
	c.DEC(addr)
	c.CMP(addr)
}

func (c *CPU) SAX(addr uint16) {
	v := c.register.A & c.register.X
	c.memo.Write(addr, v)
}

func (c *CPU) LAX(addr uint16) {
	m := c.memo.Read(addr)
	c.register.A = m
	c.register.X = m
	c.register.setNZFlag(m)
}

func (c *CPU) NOP(addr uint16) {}

func (c *CPU) INC(addr uint16) {
	v := c.memo.Read(addr) + 1
	c.memo.Write(addr, v)
	c.register.setNZFlag(v)
}

func (c *CPU) ROL(addr uint16) {
	m := c.memo.Read(addr)
	r := m << 1
	if c.register.getFlag(FLAG_C) {
		r |= 1
	}
	c.memo.Write(addr, r)
	c.register.setFlag(FLAG_C, m > 0x7f)
	c.register.setNZFlag(r)
}

func (c *CPU) ROLImp(addr uint16) {
	temp := c.register.A
	c.register.A <<= 1
	if c.register.getFlag(FLAG_C) {
		c.register.A |= 1
	}
	c.register.setFlag(FLAG_C, temp > 0x7f)
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) ROR(addr uint16) {
	m := c.memo.Read(addr)
	r := m >> 1
	if c.register.getFlag(FLAG_C) {
		r |= 0x80
	}
	c.memo.Write(addr, r)
	c.register.setFlag(FLAG_C, m&1 != 0)
	c.register.setNZFlag(r)
}

func (c *CPU) RORImp(addr uint16) {
	temp := c.register.A
	c.register.A >>= 1
	if c.register.getFlag(FLAG_C) {
		c.register.A |= 0x80
	}
	c.register.setFlag(FLAG_C, temp&1 != 0)
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) ASL(addr uint16) {
	m := c.memo.Read(addr)
	r := m << 1
	c.memo.Write(addr, r)
	c.register.setNZFlag(r)
	c.register.setFlag(FLAG_C, m >= 128)
}

func (c *CPU) ASLImp(addr uint16) {
	c.register.setFlag(FLAG_C, c.register.A >= 128)
	c.register.A <<= 1
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) LSR(addr uint16) {
	m := c.memo.Read(addr)
	c.register.setFlag(FLAG_C, m&1 != 0)
	r := m >> 1
	c.memo.Write(addr, r)
	c.register.setNZFlag(r)
}

func (c *CPU) LSRImp(addr uint16) {
	c.register.setFlag(FLAG_C, c.register.A&1 != 0)
	c.register.A >>= 1
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) TXS(addr uint16) {
	c.register.S = c.register.X
}

func (c *CPU) TSX(addr uint16) {
	c.register.X = c.register.S
	c.register.setNZFlag(c.register.X)
}

func (c *CPU) TXA(addr uint16) {
	c.register.A = c.register.X
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) TAX(addr uint16) {
	c.register.X = c.register.A
	c.register.setNZFlag(c.register.X)
}

func (c *CPU) TYA(addr uint16) {
	c.register.A = c.register.Y
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) TAY(addr uint16) {
	c.register.Y = c.register.A
	c.register.setNZFlag(c.register.Y)
}

func (c *CPU) DEY(addr uint16) {
	c.register.Y -= 1
	c.register.setNZFlag(c.register.Y)
}

func (c *CPU) DEC(addr uint16) {
	v := c.memo.Read(addr) - 1
	c.memo.Write(addr, v)
	c.register.setNZFlag(v)
}

func (c *CPU) DEX(addr uint16) {
	c.register.X -= 1
	c.register.setNZFlag(c.register.X)
}
func (c *CPU) INX(addr uint16) {
	c.register.X += 1
	c.register.setNZFlag(c.register.X)
}

func (c *CPU) INY(addr uint16) {
	c.register.Y += 1
	c.register.setNZFlag(c.register.Y)
}

func (c *CPU) CPY(addr uint16) {
	m := c.memo.Read(addr)
	c.register.setFlag(FLAG_C, c.register.Y >= m)
	c.register.setNZFlag(c.register.Y - m)
}

func (c *CPU) CPX(addr uint16) {
	m := c.memo.Read(addr)
	c.register.setFlag(FLAG_C, c.register.X >= m)
	c.register.setNZFlag(c.register.X - m)
}

func (c *CPU) EOR(operandAddr uint16) {
	c.register.A ^= c.memo.Read(operandAddr)
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) CLV(operandAddr uint16) {
	c.register.setFlag(FLAG_V, false)
}

func (c *CPU) ORA(operandAddr uint16) {
	c.register.A |= c.memo.Read(operandAddr)
	c.register.setNZFlag(c.register.A)
}
func (c *CPU) BMI(operandAddr uint16) {
	if c.register.getFlag(FLAG_N) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) PLP(operandAddr uint16) {
	oldFb := c.register.getFlag(FLAG_B)
	c.register.P = c.StackPop() | byte(FLAG_U) // unused永远是1
	c.register.setFlag(FLAG_B, oldFb)
}

func (c *CPU) PHP(operandAddr uint16) {
	v := c.register.P | byte(FLAG_B)
	c.StackPush(v)
}

func (c *CPU) PHA(operandAddr uint16) {
	c.StackPush(c.register.A)
}

func (c *CPU) CLD(operandAddr uint16) {
	c.register.setFlag(FLAG_D, false)
}

func (c *CPU) CMP(operandAddr uint16) {
	m := c.memo.Read(operandAddr)
	c.register.setFlag(FLAG_C, c.register.A >= m)
	c.register.setNZFlag(c.register.A - m)
}

func (c *CPU) AND(operandAddr uint16) {
	c.register.A = c.register.A & c.memo.Read(operandAddr)
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) PLA(operandAddr uint16) {
	c.register.A = c.StackPop()
	c.register.setNZFlag(c.register.A)
}

func (c *CPU) RTS(operandAddr uint16) {
	byte1 := c.StackPop()
	byte2 := c.StackPop()
	c.register.PC = littleEndian(byte1, byte2) + 1
}

func (c *CPU) RTI(operandAddr uint16) {
	oldFb := c.register.getFlag(FLAG_B)
	c.register.P = c.StackPop() | byte(FLAG_U) // unused永远是1
	c.register.setFlag(FLAG_B, oldFb)
	c.register.PC = c.StackPopWord()
}

func (c *CPU) BPL(operandAddr uint16) {
	if !c.register.getFlag(FLAG_N) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) BVC(operandAddr uint16) {
	if !c.register.getFlag(FLAG_V) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) BVS(operandAddr uint16) {
	if c.register.getFlag(FLAG_V) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) BIT(operandAddr uint16) {
	m := c.memo.Read(operandAddr)
	c.register.setFlag(FLAG_Z, c.register.A&m == 0)
	c.register.setFlag(FLAG_V, m&byte(FLAG_V) > 0)
	c.register.setFlag(FLAG_N, m >= 128)
}

func (c *CPU) BEQ(operandAddr uint16) {
	if c.register.getFlag(FLAG_Z) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) BNE(operandAddr uint16) {
	if !c.register.getFlag(FLAG_Z) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) LDA(operandAddr uint16) {
	c.register.A = c.memo.Read(operandAddr)
	c.register.setNZFlag(c.register.A)
	if operandAddr == 0x40 {
		println("------------LDA------------")
		println(c.memo.Read(operandAddr))
		println(c.memo.Read(operandAddr + 1))
		println(c.memo.Read(operandAddr - 1))
	}
}

func (c *CPU) BCC(operandAddr uint16) {
	if !c.register.getFlag(FLAG_C) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) CLC(operandAddr uint16) {
	c.register.setFlag(FLAG_C, false)
}

func (c *CPU) BCS(operandAddr uint16) {
	if c.register.getFlag(FLAG_C) {
		c.register.PC = operandAddr
	}
}

func (c *CPU) SEC(operandAddr uint16) {
	c.register.setFlag(FLAG_C, true)
}

func (c *CPU) JMP(operandAddr uint16) {
	c.register.PC = operandAddr
}

func (c *CPU) STX(operandAddr uint16) {
	c.memo.Write(operandAddr, c.register.X)
}

func (c *CPU) STY(operandAddr uint16) {
	c.memo.Write(operandAddr, c.register.Y)
}

func (c *CPU) STA(operandAddr uint16) {
	c.memo.Write(operandAddr, c.register.A)
}

func (c *CPU) JSR(operandAddr uint16) {
	val := c.register.PC - 1
	c.StackPush(uint8(val >> 8))
	c.StackPush(uint8(val & 0xFF))
	c.register.PC = operandAddr
}

func (c *CPU) StackPush(v uint8) {
	addr := uint16(c.register.S) + 0x0100 // sp指针从0x1FF处向下增长
	c.memo.Write(addr, v)
	c.register.S -= 1
}

func (c *CPU) StackPushWord(w uint16) {
	// 要和出站顺序对应
	c.StackPush(byte(w >> 8))
	c.StackPush(byte(w & 0xFF))
}

func (c *CPU) StackPop() uint8 {
	c.register.S += 1
	addr := uint16(c.register.S) + 0x0100
	return c.memo.Read(addr)
}

func (c *CPU) StackPopWord() uint16 {
	byte1 := c.StackPop()
	byte2 := c.StackPop()
	return littleEndian(byte1, byte2)
}

func (c *CPU) SEI(operandAddr uint16) { // set interrupt disable
	c.register.setFlag(FLAG_I, true)
}

func (c *CPU) SED(operandAddr uint16) {
	c.register.setFlag(FLAG_D, true)
}

func (c *CPU) LDX(addr uint16) {
	c.register.X = c.memo.Read(addr)
	c.register.setFlag(FLAG_Z, c.register.X == 0)
	c.register.setFlag(FLAG_N, c.register.X >= 128)
}

func (c *CPU) LDY(addr uint16) {
	c.register.Y = c.memo.Read(addr)
	c.register.setNZFlag(c.register.Y)
}

func (c *CPU) ADC(operandAddr uint16) {
	// http://www.6502.org/tutorials/vflag.html#2.4
	// FLAG_V表示的是在进行有符号数运算时的溢出位
	// FLAG_C表示的是在进行无符号运算时的溢出位
	// 如果进行的是有符号数运算，那么操作数和结果都是补码
	// 因为CPU并不知道进行的时什么类型的运算，所以CPU会同时设置FLAG_V和FLAG_C
	m := c.memo.Read(operandAddr)
	var unsignedVal = uint16(c.register.A) + uint16(m) + uint16(c.register.getCarry())
	c.register.setFlag(FLAG_C, unsignedVal > 255)
	c.register.setNZFlag(uint8(unsignedVal))

	// 注意： int16(val) 和 int16(int8(val)并不相等, 这里必须使用int16(int8(val))
	var signedVal = int16(int8(m)) + int16(int8(c.register.A)) + int16(int8(c.register.getCarry()))
	c.register.setFlag(FLAG_V, signedVal > 127 || signedVal < -128)
	c.register.A = uint8(unsignedVal)
}

// Warning: Need valid
func (c *CPU) SBC(addr uint16) {
	operand := c.memo.Read(addr)
	operand2 := ^operand
	r := uint16(c.register.A) + uint16(operand2)
	if c.register.P&byte(FLAG_C) != 0 {
		r++
	}
	r2 := uint8(r)
	c.register.setFlag(FLAG_C, r > 0xFF)
	c.register.setFlag(FLAG_V, (c.register.A^operand2)&0x80 == 0 && (c.register.A^r2)&0x80 != 0)
	c.register.setNZFlag(r2)
	c.register.A = r2
}
