package utils

func LittleEndian(byte1, byte2 byte) uint16 {
	return uint16(byte2)<<8 | uint16(byte1)
}

func Uint16(lowByte, highByte byte) *uint16 {
	v := uint16(highByte)<<8 | uint16(lowByte)
	return &v
}
