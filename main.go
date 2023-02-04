package main

import (
	"fc-emulator/emu"
	"fc-emulator/ui"
	"flag"
	"log"
)

var nesFileName = flag.String("nes", "./static/balloon.nes", "nes file path")

func setupEmulator() *emu.Emu {
	flag.Parse()
	if nesFileName == nil || len(*nesFileName) == 0 {
		log.Fatal("please specific nes file path")
	}
	emulator := emu.NewEmu(&emu.EmuOpt{Debug: false})
	err := emulator.Load(*nesFileName)
	if err != nil {
		log.Fatal("load nes file fail: ", err)
	}
	return emulator
}

func main() {
	emulator := setupEmulator()
	win := ui.NewUIWin(emulator, &ui.UIConfig{Width: 480, Height: 400})
	go func() {
		emulator.Start()
	}()
	win.ShowAndRun()
}
