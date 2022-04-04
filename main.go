package main

import (
	"fc-emulator/emu"
	"fc-emulator/ppu"
	"fc-emulator/rom"
	"fc-emulator/utils"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
	"log"
	"os"
	"time"
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

func PaletteTabItem(name string, getPaletteFn func() ppu.Palette) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		//fmt.Println(time.Now(), "w", w, "h: ", h)
		return getPaletteFn().Draw()
	})
	return container.NewTabItem(name, raster)
}

func GameScreenTabItem(screenFn func() image.Image) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		//fmt.Println(time.Now(), "w", w, "h: ", h)
		return screenFn()
	})
	return container.NewTabItem("Game", raster)
}

func RomInfoTabItem(nesRom *rom.NesRom) *container.TabItem {
	h := nesRom.Header
	f6 := h.Flag6
	f6Info := fmt.Sprintf("HasBattery: %v \n  MirrorMode: %s \n  Trainer: %v \n FourScreenMode: %v \n  MapperLowerVersion:  %v \n ",
		f6.HasBattery, rom.MirrorModeNameMap[f6.MirrorMode], f6.Trainer, f6.FourScreenMode, f6.MapperLowerVersion)
	headerMsg := fmt.Sprintf("MapperNumber: %d \n prgCount: %d \n chrCount: %d \n Flag6: %s \n ",
		h.MapperNumber, h.PrgCount, h.ChrCount, f6Info)
	romMsg := fmt.Sprintf("PrgRomSize(kb): %d \n ChrRomSize(kb): %d  \n ", len(nesRom.PrgRom)/utils.Kb, len(nesRom.ChrRom)/utils.Kb)
	return container.NewTabItem("Rom Info", widget.NewTextGridFromString(headerMsg+"\n"+romMsg))
}

func ctrlStatusStr(ppuCtrl ppu.PPUCTRL) string {
	msg := fmt.Sprintf(`NameTable BaseAddr: 0x%X
SpritePatternTableAddressFor88Mode: 0x%X
BackgroundPatternTableAddress: 0x%X
SpritePatternTableAddress:  0x%X
SpriteSize: %v
MasterSlaveSelect: %v
CanGenerateNMIBreakAtStartOfVerticalBlankingInterval: %v`,
		ppuCtrl.NameTableBaseAddress(),
		ppuCtrl.SpritePatternTableAddressFor88Mode(),
		ppuCtrl.BackgroundPatternTableAddress(),
		ppuCtrl.SpritePatternTableAddress(),
		ppuCtrl.SpriteSize(),
		ppuCtrl.MasterSlaveSelect(),
		ppuCtrl.CanGenerateNMIBreakAtStartOfVerticalBlankingInterval())
	return msg
}
func PPUCtrlStatusTabItem(pu *ppu.PPUImpl) *container.TabItem {
	txtWidget := widget.NewTextGridFromString(ctrlStatusStr(pu.Register.PPUCTRL))
	go func() {
		c := time.Tick(1 * time.Second)
		for {
			<-c
			txtWidget.SetText(ctrlStatusStr(pu.Register.PPUCTRL))
		}
	}()

	return container.NewTabItem("PPU CTRL Status", txtWidget)
}

func main() {
	emulator := setupEmulator()
	pu := emulator.PPU.(*ppu.PPUImpl)

	screenTabItem := GameScreenTabItem(pu.Render)
	bgPaletteTabItem := PaletteTabItem("BG Palette", pu.BgPalette)
	spritePaletteTabItem := PaletteTabItem("Sprite Palette", pu.SpritePalette)

	myApp := app.New()
	win := myApp.NewWindow("Raster")
	tabs := container.NewAppTabs(
		screenTabItem,
		bgPaletteTabItem,
		spritePaletteTabItem,
		RomInfoTabItem(emulator.Rom),
		PPUCtrlStatusTabItem(pu),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	emulator.FrameCallback = screenTabItem.Content.Refresh
	go func() {
		emulator.Start()
	}()

	win.SetContent(tabs)
	win.Resize(fyne.NewSize(480, 400))
	win.ShowAndRun()
}
func main2() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")

	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
		container.NewTabItem("Tab 2", widget.NewLabel("World!")),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.Size{
		Width:  500,
		Height: 400,
	})

	myWindow.ShowAndRun()
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
	fmt.Println("update", formatted)
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
