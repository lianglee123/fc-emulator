package main

import (
	"fc-emulator/emu"
	"flag"
	"log"
)

var fileName = flag.String("nes", "./static/mario.nes", "nes file name")

func main() {
	flag.Parse()
	if fileName == nil || len(*fileName) == 0 {
		log.Fatal("please specific nes file name")
	}
	emulator := emu.NewEmu()
	err := emulator.Load(*fileName)
	if err != nil {
		log.Fatal("load nes file fail: ", err)
	}
	emulator.Start()
}
