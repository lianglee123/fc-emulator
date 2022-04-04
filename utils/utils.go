package utils

import (
	"errors"
	"fmt"
)

func NewError(errInfos ...interface{}) error {
	return errors.New(fmt.Sprintln(errInfos...))
}

// zero index base
func GetBitFromLeft(v byte, index int) byte {
	if index >= 8 || index < 0 {
		panic("index should between [1, 8)")
	}
	return (v >> (7 - index)) & 0x1
}

// zero index base
func GetBitFromRight(v byte, index int) byte {
	if index >= 8 || index < 0 {
		panic("index should between [0, 8)")
	}
	return (v >> index) & 0x1
}
