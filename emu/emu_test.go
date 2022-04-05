package emu

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	a := []byte{}
	fmt.Println(a[100:1000])
}
