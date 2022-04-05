package pad

import "fmt"

type Pad interface {
	ReadForCPU() byte
	UpdateButton(buttonType ButtonType, pressDown bool)
	WriteForCPU(value byte)
}

type DefaultPad struct {
	strobe      bool
	data        byte
	buttonIndex int
}

func NewPad() Pad {
	return &DefaultPad{
		strobe:      false,
		data:        0,
		buttonIndex: 0,
	}
}

type ButtonType byte

const (
	BUTTON_RIGHT  ButtonType = 0x1 << 7
	BUTTON_LEFT   ButtonType = 0x1 << 6
	BUTTON_DOWN   ButtonType = 0x1 << 5
	BUTTON_UP     ButtonType = 0x1 << 4
	BUTTON_START  ButtonType = 0x1 << 3
	BUTTON_SELECT ButtonType = 0x1 << 2
	BUTTON_B      ButtonType = 0x1 << 1
	BUTTON_A      ButtonType = 0x1 << 0
)

func (p *DefaultPad) UpdateButton(buttonType ButtonType, pressDown bool) {
	fmt.Printf("update button %v %v", buttonType, pressDown)
	if pressDown {
		p.data |= byte(buttonType)
	} else {
		p.data = p.data & (^byte(buttonType))
	}
}
func (p *DefaultPad) WriteForCPU(value byte) {
	//fmt.Printf("write pad 0x%b\n", value)
	if value&0x01 == 0x01 {
		p.strobe = true
	} else {
		p.strobe = false
		p.buttonIndex = 0
	}
}

func (p *DefaultPad) ReadForCPU() byte {
	if p.strobe {
		return p.data & byte(BUTTON_A)
	} else {
		res := (p.data >> p.buttonIndex) & 0x01
		p.buttonIndex += 1
		fmt.Printf("read pad 0x%b\n", res)
		return res
	}
}
