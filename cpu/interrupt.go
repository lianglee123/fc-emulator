package cpu

import "fmt"

const (
	IV_NMI   uint16 = 0xFFFA
	IV_RESET uint16 = 0xFFFC
	IV_IRQ   uint16 = 0xFFFE
	IV_BRK   uint16 = 0xFFFE
)

func (c *CPU) ExecIRQ() {
	fmt.Println("EXECUTE IRQ")
	if c.register.getFlag(FLAG_I) {
		fmt.Println("IRQ Skip")
		return
	}

	c.StackPushWord(c.register.PC)
	c.StackPush((c.register.P | uint8(FLAG_U)) | (^uint8(FLAG_B)))
	c.register.setFlag(FLAG_I, true)
	c.register.PC = c.memo.ReadWord(IV_IRQ)
}

// 进入中断后，由中断handler负责pop 原来的pc, 返回到原来的执行链路上。
func (c *CPU) ExecNMI() {
	//fmt.Println("EXECUTE NMI")
	c.StackPushWord(c.register.PC)
	c.StackPush((c.register.P | uint8(FLAG_U)) | (^uint8(FLAG_B)))
	c.register.setFlag(FLAG_I, true)
	c.register.PC = c.memo.ReadWord(IV_NMI)
}

func (c *CPU) ExecBRK() {
	fmt.Println("Exec BRK")
	c.StackPushWord(c.register.PC)
	c.StackPush(c.register.P | uint8(FLAG_U) | uint8(FLAG_B))
	c.register.setFlag(FLAG_I, true)
	c.register.PC = c.memo.ReadWord(IV_BRK)
}
