package test

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	// 这里打印的都是数值的原码
	var a byte = 128
	fmt.Println(int16(a), int8(a), int16(int8(a)))
	fmt.Println(uint16(a), uint8(a), uint16(uint8(a)))

}

func TestSlice(t *testing.T) {
	var a = []int{1, 2, 3, 4, 5, 6}
	var b = a[1:4]
	fmt.Println(a)
	fmt.Println(b)
	b[0] = 222
	fmt.Println(a)
	fmt.Println(b)
}

func TestArray(t *testing.T) {
	opcode := [8]int{
		1: 2,
		5: 6,
	}
	fmt.Println(opcode)
}

func TestArrary2(t *testing.T) {
	opcodes := [8]int{}
	fmt.Println(opcodes)
}

func TestInt(t *testing.T) {
	var a uint8 = 1
	var b uint8 = 2
	fmt.Println(a - b)
}
