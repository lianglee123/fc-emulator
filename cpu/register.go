package cpu

type Flag uint8

const (
	FLAG_C Flag = 1 << 0 // Carry
	FLAG_Z Flag = 1 << 1 // Zero
	FLAG_I Flag = 1 << 2 // Disable interrupt
	FLAG_D Flag = 1 << 3 // Decimal Mode ( unused in nes )
	FLAG_B Flag = 1 << 4 // Break
	FLAG_U Flag = 1 << 5 // Unused ( always 1 )
	FLAG_V Flag = 1 << 6 // Overflow
	FLAG_N Flag = 1 << 7 // Negative
)

func (r *Register) getFlag(f Flag) bool {
	return (r.P & uint8(f)) != 0
}

func (r *Register) getCarry() uint8 {
	if r.getFlag(FLAG_C) {
		return 1
	} else {
		return 0
	}
}

func (r *Register) setFlag(f Flag, val bool) {
	if val {
		r.P = r.P | uint8(f)
	} else {
		r.P = r.P & (^uint8(f))
	}
}

func (r *Register) setNZFlag(value uint8) {
	r.setFlag(FLAG_Z, value == 0)
	r.setFlag(FLAG_N, value >= 128)
}

type Register struct {
	PC uint16
	S  uint8
	P  uint8
	A  uint8
	X  uint8
	Y  uint8
}

func (r *Register) IncreasePC() uint16 {
	r.PC++
	return r.PC
}
