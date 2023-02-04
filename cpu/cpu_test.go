package cpu

import (
	"errors"
	"fc-emulator/cpu/addressing"
	"fc-emulator/memo"
	"fc-emulator/pad"
	"fc-emulator/ppu"
	"fc-emulator/rom"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

type CPUStatus struct {
	Register Register
}

type LogDiffer struct {
	logs []string
	Ptr  int
}

func NewLogDiffer(logName string) *LogDiffer {
	d := &LogDiffer{}
	data, err := ioutil.ReadFile(logName)
	if err != nil {
		panic(err)
	}
	d.logs = strings.Split(string(data), "\n")
	return d
}
func (d *LogDiffer) HasNext() bool {
	return d.Ptr < len(d.logs)
}

func (d *LogDiffer) Diff(traceLog *TraceLog, cpu *CPU) error {
	var oprandAddr string
	if traceLog.Mode == addressing.IMM {
		oprandAddr = fmt.Sprintf("%04X", cpu.memo.Read(traceLog.Addr))
	} else if traceLog.Mode == addressing.IMP {
		oprandAddr = "----"
	} else {
		oprandAddr = fmt.Sprintf("%04X", traceLog.Addr)
	}
	msgTpl := "%04d PC: %04X %02X OP: %s %s(%s) A:%02X X:%02X Y:%02X P:%02X S:%02X"
	actualLog := fmt.Sprintf(msgTpl,
		d.Ptr+1, traceLog.OldReg.PC, traceLog.opcodeNumber,
		traceLog.Code, oprandAddr, traceLog.Mode,
		traceLog.OldReg.A, traceLog.OldReg.X, traceLog.OldReg.Y, traceLog.OldReg.P, traceLog.OldReg.S)
	expectLog := d.logs[d.Ptr]
	d.Ptr += 1
	fmt.Println(actualLog)
	if actualLog != expectLog {
		return errors.New(fmt.Sprintf("Expect: %s\nActual: %s", expectLog, actualLog))
	} else {
		return nil
	}
}

func TestCPU(t *testing.T) {
	nesRom, err := rom.LoadNesRom("nestest.nes")
	require.NoError(t, err)
	cpuMemo := memo.NewMemo(nesRom, ppu.NewPPU(nesRom), pad.NewPad(), pad.NewPad())
	c := NewCPU(cpuMemo, true)
	c.Reset()
	c.register.PC = 0xC000
	c.register.P = 0x24 // only FLAG_I and FLAG_U enable
	logDiffer := NewLogDiffer("nestest.log")
	fmt.Println("code count: ", codeCount())
	for logDiffer.HasNext() {
		traceLog, err := c.ExecuteOneInstruction()
		if err != nil {
			fmt.Println("code count: ", codeCount())
		}
		require.NoError(t, err)
		err = logDiffer.Diff(traceLog, c)
		if err != nil {
			fmt.Println("code count: ", codeCount())
			require.NoError(t, err)
		}
	}
}
