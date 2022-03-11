package rom

import (
	"errors"
	"fc-emulator/utils"
	"fmt"
	"io/ioutil"
)

const nesPrefix = "NES\x1A"

type NesRom struct {
	PrgCount     int // unit is 16k
	ChrCount     int //  unit is 8k
	MapperNumber int
	IsTrainer    bool
	HasBattery   bool
	PrgRom       []byte
	ChrRom       []byte
	Trainer      []byte
}

func (r *NesRom) String() string {
	return fmt.Sprintln("MapperNumber:", r.MapperNumber, "prgCount:", r.PrgCount,
		"chrCount:", r.ChrCount, "PrgRomSize(kb):", len(r.PrgRom)/utils.Kb, "ChrRomSize(kb):", len(r.ChrRom)/utils.Kb)
}

func LoadNesRom(filename string) (*NesRom, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if string(data[:4]) != nesPrefix {
		return nil, errors.New("nes format error")
	}
	if data[7]&0x0c == 0x08 {
		return nil, errors.New("NES2.0 is not currently supported")
	}
	rom := &NesRom{
		PrgCount:     int(data[4]),
		ChrCount:     int(data[5]),
		MapperNumber: int((data[7] & 0xF0) | (data[6] >> 4)),
		IsTrainer:    (data[6] & 0x04) != 0x00, // 第4bit
		HasBattery:   (data[6] & 0x02) != 0x00, // 第2bit
	}
	if rom.MapperNumber != 0 {
		return nil, utils.NewError("Not support Mapper ", rom.MapperNumber)
	}
	prgStartIndex := 16
	if rom.IsTrainer {
		rom.Trainer = data[16:512]
		prgStartIndex += 512
	} else {
		rom.Trainer = make([]byte, 512)
	}
	prgSize := rom.PrgCount * 16 * utils.Kb
	rom.PrgRom = data[prgStartIndex : prgStartIndex+prgSize]
	rom.ChrRom = data[prgStartIndex+prgSize:]
	return rom, nil
}
