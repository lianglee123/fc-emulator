package cpu

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegister1(t *testing.T) {
	rg := &Register{}
	rg.setFlag(FLAG_I, true)
	rg.setFlag(FLAG_U, true)
	require.Equal(t, rg.P, 0x24)
}
