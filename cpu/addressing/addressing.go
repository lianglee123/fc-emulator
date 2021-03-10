package addressing

type Mode int

// 寻址模式
//go:generate stringer -type=Mode -output addressing.string.go
const (
	IMP Mode = iota
	IMM
	ZP
	ZPX
	ZPY
	IZX
	IZY
	ABS
	ABX
	ABY
	IND
	REL
)
