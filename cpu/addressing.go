package cpu

import (
	"fc-emulator/cpu/addressing"
	"fmt"
)

func (c *CPU) readOp() byte {
	return 0
}

// 注意，寻址返回的是一个地址，而不是值
// gones 的实现，IMM寻址返回的是一个地址
// emulator.py的实现，IMM寻址返回的是一个值
func (c *CPU) AddressIMM() uint16 {
	addr := c.register.PC
	c.increasePC()
	return addr
}

// Zero Page
func (c *CPU) AddressZP() uint16 {
	addr := c.register.PC
	c.increasePC()
	return uint16(c.memo.Read(addr))
}

// Zero Page X
func (c *CPU) AddressZPX() uint16 {
	addr := c.register.PC
	c.increasePC()
	return uint16(c.memo.Read(addr) + c.register.X)
}

// Zero Page Y
func (c *CPU) AddressZPY() uint16 {
	addr := c.register.PC
	c.increasePC()
	return uint16(c.memo.Read(addr) + c.register.Y)
}

func (c *CPU) AddressAbs() uint16 {
	val := c.memo.ReadWord(c.register.PC)
	c.increasePC()
	c.increasePC()
	return val
}

func (c *CPU) AddressAbsX() uint16 {
	addr := c.AddressAbs()
	return addr + uint16(c.register.X)
}

func (c *CPU) AddressAbsY() uint16 {
	addr := c.AddressAbs()
	return addr + uint16(c.register.Y)
}

func (c *CPU) AddressIndirect() uint16 {
	addr := c.AddressAbs()
	// HardWare Bug: 无法跨 Page
	// 例如JMP ($10FF), 理论上讲要读取$10FF和$1100这两个字节的数据,
	// 因为$10FF和$1100不在同一个Page上，所以
	// 实际上是读取的$10FF和$1000这两个字节的数据.
	if addr&0xFF == 0xFF {
		byte1 := c.memo.Read(addr)
		byte2 := c.memo.Read(addr & 0xFF00)
		return littleEndian(byte1, byte2)
	} else {
		return c.memo.ReadWord(addr)
	}
}

func (c *CPU) AddressIndexedDirectX() uint16 {
	if c.register.PC == 0xCFF2 || c.register.PC == 0xCFF0 || c.register.PC == 0xCFF1 || c.register.PC == 0xCFF3 {
		m := c.memo.Read(c.register.PC)
		fmt.Println("memeo vale: ", m)
	}
	v := c.memo.Read(c.register.PC)
	c.register.IncreasePC()
	addr := c.register.X + v
	byte1 := c.memo.Read(uint16(addr)) // 注意，这里不能用c.memo.ReadWord(), 这两者在逻辑上是有差距的
	byte2 := c.memo.Read(uint16(addr + 1))
	return littleEndian(byte1, byte2)
}

func (c *CPU) AddressIndirectIndexedY() uint16 {
	v := c.memo.Read(c.register.PC)
	c.register.IncreasePC()
	byte1 := c.memo.Read(uint16(v))
	byte2 := c.memo.Read(uint16(v + 1))
	baseAddr := littleEndian(byte1, byte2)
	addr := baseAddr + uint16(c.register.Y)
	// Warning: different with py
	//if isCrossPage(baseAddr, addr) {
	//	return
	//}
	return addr
}

func isCrossPage(v1, v2 uint16) bool {
	return v1&0xFF00 != v2&0xFF00
}
func (c *CPU) AddressRel() uint16 {
	v := c.memo.Read(c.register.PC)
	c.register.IncreasePC()
	return uint16(int64(int8(v)) + int64(c.register.PC))
}

func (c *CPU) Addressing(mode addressing.Mode) uint16 {
	switch mode {
	case addressing.IMP:
		return 0
	case addressing.IMM:
		return c.AddressIMM()
	case addressing.ZPG:
		return c.AddressZP()
	case addressing.ZPX:
		return c.AddressZPX()
	case addressing.ZPY:
		return c.AddressZPY()
	case addressing.ABS:
		return c.AddressAbs()
	case addressing.ABX:
		return c.AddressAbsX()
	case addressing.ABY:
		return c.AddressAbsY()
	case addressing.REL:
		return c.AddressRel()
	case addressing.IND:
		return c.AddressIndirect()
	case addressing.INX:
		return c.AddressIndexedDirectX()
	case addressing.INY:
		return c.AddressIndirectIndexedY()
	default:
		panic(fmt.Errorf("unsupported addressing mode: %s", mode))
	}
}
