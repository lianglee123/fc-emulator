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
func (c *CPU) AddressIMM() (uint16, bool) {
	addr := c.register.PC
	c.increasePC()
	return addr, false
}

// 只有absX, absY， idx这三种寻址下的指令才有额外cycle的问题，所以只用实现这三种寻址模式即可。
// Zero Page
func (c *CPU) AddressZP() (uint16, bool) {
	addr := c.register.PC
	c.increasePC()
	return uint16(c.memo.Read(addr)), false
}

// Zero Page X
func (c *CPU) AddressZPX() (uint16, bool) {
	addr := c.register.PC
	c.increasePC()
	return uint16(c.memo.Read(addr) + c.register.X), false
}

// Zero Page Y
func (c *CPU) AddressZPY() (uint16, bool) {
	addr := c.register.PC
	c.increasePC()
	return uint16(c.memo.Read(addr) + c.register.Y), false
}

func (c *CPU) AddressAbs() (uint16, bool) {

	val := c.memo.ReadWord(c.register.PC)
	c.increasePC()
	c.increasePC()
	return val, false
}

func (c *CPU) AddressAbsX() (uint16, bool) {
	baseAddr, _ := c.AddressAbs()
	addr := baseAddr + uint16(c.register.X)
	return addr, c.IsCrossPage(baseAddr, addr)
}

func (c *CPU) AddressAbsY() (uint16, bool) {
	baseAddr, _ := c.AddressAbs()
	addr := baseAddr + uint16(c.register.Y)
	return addr, c.IsCrossPage(baseAddr, addr)
}

func (c *CPU) AddressIndirect() (uint16, bool) {
	addr, _ := c.AddressAbs()
	// HardWare Bug: 无法跨 Page
	// 例如JMP ($10FF), 理论上讲要读取$10FF和$1100这两个字节的数据,
	// 因为$10FF和$1100不在同一个Page上，所以
	// 实际上是读取的$10FF和$1000这两个字节的数据.
	if addr&0xFF == 0xFF {
		byte1 := c.memo.Read(addr)
		byte2 := c.memo.Read(addr & 0xFF00)
		return littleEndian(byte1, byte2), false
	} else {
		return c.memo.ReadWord(addr), false
	}
}

func (c *CPU) AddressIndexedDirectX() (uint16, bool) {
	if c.register.PC == 0xCFF2 || c.register.PC == 0xCFF0 || c.register.PC == 0xCFF1 || c.register.PC == 0xCFF3 {
		m := c.memo.Read(c.register.PC)
		fmt.Println("memeo vale: ", m)
	}
	v := c.memo.Read(c.register.PC)
	c.register.IncreasePC()
	addr := c.register.X + v
	byte1 := c.memo.Read(uint16(addr)) // 注意，这里不能用c.memo.ReadWord(), 这两者在逻辑上是有差距的
	byte2 := c.memo.Read(uint16(addr + 1))
	return littleEndian(byte1, byte2), false
}

func (c *CPU) AddressIndirectIndexedY() (uint16, bool) {
	v := c.memo.Read(c.register.PC)
	c.register.IncreasePC()
	byte1 := c.memo.Read(uint16(v))
	byte2 := c.memo.Read(uint16(v + 1))
	baseAddr := littleEndian(byte1, byte2)
	addr := baseAddr + uint16(c.register.Y)
	return addr, c.IsCrossPage(baseAddr, addr)
}

func isCrossPage(v1, v2 uint16) bool {
	return v1&0xFF00 != v2&0xFF00
}

func (c *CPU) AddressRel() (uint16, bool) {
	v := c.memo.Read(c.register.PC)
	c.register.IncreasePC()
	return uint16(int64(int8(v)) + int64(c.register.PC)), false
}

func (c *CPU) IsCrossPage(addr1, addr2 uint16) bool {
	return (addr1 & 0xFF) != (addr2 & 0xFF)
}

func (c *CPU) Addressing(mode addressing.Mode) (uint16, bool) {
	switch mode {
	case addressing.IMP:
		return 0, false
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
