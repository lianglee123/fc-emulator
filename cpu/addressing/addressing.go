package addressing

type Mode int

// 寻址模式
//go:generate stringer -type=Mode -output addressing.string.go
const (
	IMP Mode = iota
	IMM
	ZPG
	ZPX
	ZPY
	INX
	INY
	ABS
	ABX
	ABY
	IND
	REL
)
