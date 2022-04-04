package rom

import (
	"fc-emulator/utils"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadNesRom(t *testing.T) {
	rom, err := LoadNesRom("../static/nestest.nes")
	require.NoError(t, err)
	fmt.Println(rom.String())
	require.Equal(t, len(rom.PrgRom), 16*utils.Kb)
	require.Equal(t, len(rom.ChrRom), 8*utils.Kb)
	rom, err = LoadNesRom("../static/mario.nes")
	require.NoError(t, err)
	require.Equal(t, len(rom.PrgRom), 32*utils.Kb)
	require.Equal(t, len(rom.ChrRom), 8*utils.Kb)
	fmt.Println(rom.String())

	rom, err = LoadNesRom("../static/balloon.nes")
	require.NoError(t, err)
	require.Equal(t, len(rom.PrgRom), 16*utils.Kb)
	require.Equal(t, len(rom.ChrRom), 8*utils.Kb)
	fmt.Println(rom.String())
}
