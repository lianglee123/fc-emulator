package ui

import (
	"fc-emulator/emu"
	"fc-emulator/pad"
	"fc-emulator/ppu"
	"fc-emulator/rom"
	"fc-emulator/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image"
	"time"
)

func PaletteTabItem(name string, getPaletteFn func() ppu.Palette) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		return getPaletteFn().Draw()
	})
	return container.NewTabItem(name, raster)
}

func GameScreenTabItem(screenFn func() image.Image) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		return screenFn()
	})
	raster.ScaleMode = canvas.ImageScalePixels
	return container.NewTabItem("Game", raster)
}

func SpriteTableTabItem(source func() image.Image) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		return source()
	})
	raster.ScaleMode = canvas.ImageScalePixels
	return container.NewTabItem("Sprite Table", raster)
}

func SpriteTabItem(source func() image.Image) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		fmt.Println("sprites refresh")
		return source()
	})
	raster.ScaleMode = canvas.ImageScalePixels
	return container.NewTabItem("Sprites", raster)
}

func PaletteColorTabItem() *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		return ppu.SystemPaletteColor()
	})
	raster.ScaleMode = canvas.ImageScalePixels
	return container.NewTabItem("Sprites", raster)
}

func BgPatternTableTabItem(source func() image.Image) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		return source()
	})
	raster.ScaleMode = canvas.ImageScalePixels
	return container.NewTabItem("Bg Pattern", raster)
}

func SpritePatternTableTabItem(source func() image.Image) *container.TabItem {
	raster := canvas.NewRaster(func(w, h int) image.Image {
		return source()
	})
	raster.ScaleMode = canvas.ImageScalePixels
	return container.NewTabItem("Sprite Pattern", raster)
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

func ctrlStatusStr(ppuCtrl ppu.PPUCTRL) string {
	msg := fmt.Sprintf(`
NameTable BaseAddr: 0x%X
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

type UIConfig struct {
	Width  int
	Height int
}

func (c *UIConfig) Clean() {
	if c.Width <= 0 {
		c.Width = 480
	}
	if c.Height <= 0 {
		c.Height = 400
	}
}

func NewUIWin(emulator *emu.Emu, config *UIConfig) fyne.Window {
	config.Clean()
	pu := emulator.PPU.(*ppu.PPUImpl)
	gameTabItem := GameScreenTabItem(pu.Render)
	bgPaletteTabItem := PaletteTabItem("BG Palette", pu.BgPalette)
	spritePaletteTabItem := PaletteTabItem("Sprite Palette", pu.SpritePalette)
	bgPatternTableTabItem := BgPatternTableTabItem(pu.DrawBGPatternTable)
	spritePatternTableTabItem := SpritePatternTableTabItem(pu.DrawSpritePatternTable)
	spriteTableTabItem := SpriteTableTabItem(pu.DrawSpriteTable)
	spriteTabItem := SpriteTabItem(pu.RenderSprites)

	myApp := app.New()
	win := myApp.NewWindow("FC Emulator")
	tabs := container.NewAppTabs(
		gameTabItem,
		spriteTabItem,
		spriteTableTabItem,
		bgPaletteTabItem,
		spritePaletteTabItem,
		RomInfoTabItem(emulator.Rom),
		PPUCtrlStatusTabItem(pu),
		bgPatternTableTabItem,
		spritePatternTableTabItem,
		PaletteColorTabItem(),
	)

	win.Canvas().(desktop.Canvas).SetOnKeyDown(func(event *fyne.KeyEvent) {
		if tabs.Selected() != gameTabItem {
			return
		}
		keyMap := map[fyne.KeyName]pad.ButtonType{
			fyne.KeyW: pad.BUTTON_UP,
			fyne.KeyS: pad.BUTTON_DOWN,
			fyne.KeyA: pad.BUTTON_LEFT,
			fyne.KeyD: pad.BUTTON_RIGHT,
			fyne.KeyJ: pad.BUTTON_A,
			fyne.KeyK: pad.BUTTON_B,
			fyne.KeyU: pad.BUTTON_SELECT,
			fyne.KeyI: pad.BUTTON_START,
		}
		if button, ok := keyMap[event.Name]; ok {
			emulator.Pad1.UpdateButton(button, true)
		}
	})
	win.Canvas().(desktop.Canvas).SetOnKeyUp(func(event *fyne.KeyEvent) {
		if tabs.Selected() != gameTabItem {
			return
		}
		keyMap := map[fyne.KeyName]pad.ButtonType{
			fyne.KeyW: pad.BUTTON_UP,
			fyne.KeyS: pad.BUTTON_DOWN,
			fyne.KeyA: pad.BUTTON_LEFT,
			fyne.KeyD: pad.BUTTON_RIGHT,
			fyne.KeyJ: pad.BUTTON_A,
			fyne.KeyK: pad.BUTTON_B,
			fyne.KeyU: pad.BUTTON_SELECT,
			fyne.KeyI: pad.BUTTON_START,
		}
		if button, ok := keyMap[event.Name]; ok {
			emulator.Pad1.UpdateButton(button, false)
		}
	})

	tabs.SetTabLocation(container.TabLocationLeading)
	emulator.FrameCallback = func() {
		tabs.Refresh()
	}
	win.SetContent(tabs)
	win.Resize(fyne.NewSize(float32(config.Width), float32(config.Height)))
	return win
}
