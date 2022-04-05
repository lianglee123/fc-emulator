package ppu

import "fc-emulator/utils"

type PPUFlag uint8

type RegisterManager struct {
	PPUCTRL   PPUCTRL // $2000
	PPUMASK   uint8   // $2001
	PPUSTATUS uint8   // $2002
	OAMADDR   uint8   // $2003
	//OAMDATA   uint8      // $2004
	PPUSCROLL *PPUSCroll // $2005
	PPUADDR   *PPUADDR   // $2006
	OAMDMA    uint8      // $4014
}

func NewRegisterManager() *RegisterManager {
	return &RegisterManager{
		PPUCTRL:   0,
		PPUMASK:   0,
		PPUSTATUS: 0b10100000,
		OAMADDR:   0,
		//OAMDATA:   0,
		PPUSCROLL: &PPUSCroll{},
		PPUADDR:   &PPUADDR{},
		OAMDMA:    0,
	}
}

// STATUS
//7  bit  0
//---- ----
//VSO. ....
//|||| ||||
//|||+-++++- Least significant bits previously written into a PPU register
//|||        (due to register not being updated for this address)
//||+------- Sprite overflow. The intent was for this flag to be set
//||         whenever more than eight sprites appear on a scanline, but a
//||         hardware bug causes the actual behavior to be more complicated
//||         and generate false positives as well as false negatives; see
//||         PPU sprite evaluation. This flag is set during sprite
//||         evaluation and cleared at dot 1 (the second dot) of the
//||         pre-render line.
//|+-------- Sprite 0 Hit.  Set when a nonzero pixel of sprite 0 overlaps
//|          a nonzero background pixel; cleared at dot 1 of the pre-render
//|          line.  Used for raster timing.
//+--------- Vertical blank has started (0: not in vblank; 1: in vblank).
//Set at dot 1 of line 241 (the line *after* the post-render
//line); cleared after reading $2002 and at dot 1 of the
//pre-render line.

// 7  bit  0
//---- ----
//VPHB SINN
//|||| ||||
//|||| ||++- Base nametable address
//|||| ||    (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
//|||| |+--- VRAM address increment per CPU read/write of PPUDATA
//|||| |     (0: add 1, going across; 1: add 32, going down)
//|||| +---- Sprite pattern table address for 8x8 sprites
//||||       (0: $0000; 1: $1000; ignored in 8x16 mode)
//|||+------ Background pattern table address (0: $0000; 1: $1000)
//||+------- Sprite size (0: 8x8 pixels; 1: 8x16 pixels)
//|+-------- PPU master/slave select
//|          (0: read backdrop from EXT pins; 1: output color on EXT pins)
//+--------- Generate an NMI at the start of the
//           vertical blanking interval (0: off; 1: on)
func (r *RegisterManager) SpriteSize() {

}

type PPUADDR struct {
	SecondWrite          bool
	MostSignificantByte  byte
	LeastSignificantByte byte
}

func (p *PPUADDR) Write(val byte) {
	if p.SecondWrite {
		p.LeastSignificantByte = val
	} else {
		p.MostSignificantByte = val & 0x3F
	}
	p.SecondWrite = !p.SecondWrite
}

func (p *PPUADDR) Value() uint16 {
	if p.SecondWrite {
		panic("waiting write second byte, can not get value")
	}
	return uint16(p.MostSignificantByte)<<8 | uint16(p.LeastSignificantByte)
}

func (p *PPUADDR) Add(val byte) {
	v := (p.Value() + uint16(val)) & 0x3FFF
	p.LeastSignificantByte = byte(v & 0xFF)
	p.MostSignificantByte = byte(v >> 8)
}

type PPUSCroll struct {
	SecondWrite bool
	XScroll     byte
	YScroll     byte
}

func (p *PPUSCroll) Write(val byte) {
	if p.SecondWrite {
		p.YScroll = val
	} else {
		p.XScroll = val
	}
	p.SecondWrite = !p.SecondWrite
}

func (p *PPUSCroll) Value() (byte, byte) {
	if p.SecondWrite {
		panic("waiting write second byte, can not get value")
	}
	return p.XScroll, p.YScroll
}

//func (r *Register) PPUCTRL

// https://www.nesdev.org/wiki/PPU_registers#Controller_($2000)_%3E_write

//  * Common name: '''PPUCTRL'''
//  * Description: PPU control register
//  * Access: write
//
//  Various flags controlling PPU operation
//  7  bit  0
//  ---- ----
//  VPHB SINN
//  |||| ||||
//  |||| ||++- Base nametable address
//  |||| ||    (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
//  |||| |+--- VRAM address increment per CPU read/write of PPUDATA
//  |||| |     (0: add 1, going across; 1: add 32, going down)
//  |||| +---- Sprite pattern table address for 8x8 sprites
//  ||||       (0: $0000; 1: $1000; ignored in 8x16 mode)
//  |||+------ Background pattern table address (0: $0000; 1: $1000)
//  ||+------- Sprite size (0: 8x8 pixels; 1: 8x16 pixels)
//  |+-------- PPU master/slave select
//  |          (0: read backdrop from EXT pins; 1: output color on EXT pins)
//  +--------- Generate an [[NMI]] at the start of the Vertical blanking interval (0: off; 1: on)
//
//  Equivalently, bits 1 and 0 are the most significant bit of the scrolling coordinates (see [[PPU_nametables|Nametables]] and [[#PPUSCROLL|PPUSCROLL]]):
//  7  bit  0
//  ---- ----
//  .... ..YX
//  ||
//  |+- 1: Add 256 to the X scroll position
//  +-- 1: Add 240 to the Y scroll position
//
//  Another way of seeing the explanation above is that when you reach the end of a nametable, you must switch to the next one, hence, changing the nametable address.
//
//  [[PPU power up state|After power/reset]], writes to this register are ignored for about 30,000 cycles.
//
//  If the PPU is currently in vertical blank, and the [[#PPUSTATUS|PPUSTATUS]] ($2002) vblank flag is still set (1), changing the NMI flag in bit 7 of $2000 from 0 to 1 will immediately generate an NMI.
//  This can result in graphical errors (most likely a misplaced scroll) if the NMI routine is executed too late in the blanking period to finish on time.
//  To avoid this problem it is prudent to read $2002 immediately before writing $2000 to clear the vblank flag.
//
//  For more explanation of sprite size, see: [[Sprite size]]
//
//  ==== Master/slave mode and the EXT pins ====
//  When bit 6 of PPUCTRL is clear (the usual case), the PPU gets the [[PPU_palettes|palette index]] for the background color from the EXT pins. The stock NES grounds these pins, making palette index 0 the background color as expected. A secondary picture generator connected to the EXT pins would be able to replace the background with a different image using colors from the background palette, which could be used e.g. to implement parallax scrolling.
//
//  Setting bit 6 causes the PPU to output the lower four bits of the palette memory index on the EXT pins for each pixel (in addition to normal image drawing) - since only four bits are output, background and sprite pixels can't normally be distinguished this way. As the EXT pins are grounded on an unmodified NES, setting bit 6 is discouraged as it could potentially damage the chip whenever it outputs a non-zero pixel value (due to it effectively shorting Vcc and GND together). Looking at the relevant circuitry in [[Visual 2C02]], it appears that the [[PPU palettes|background palette hack]] would not be functional for output from the EXT pins; they would always output index 0 for the background color.
//
//  ==== Bit 0 race condition ====
//  Be very careful when writing to this register outside vertical blanking if you are using vertical mirroring (horizontal arrangement) or 4-screen VRAM.
//  For specific CPU-PPU alignments, [//forums.nesdev.org/viewtopic.php?p=112424#p112424 a write that starts] on [[PPU scrolling#At dot 257 of each scanline|dot 257]] will cause only the next scanline to be erroneously drawn from the left nametable.
//  This can cause a visible glitch, and it can also interfere with sprite 0 hit for that scanline (by being drawn with the wrong background).
//
//  The glitch has no effect in horizontal or one-screen mirroring.
//  Only writes that start on dot 257 and continue through dot 258 can cause this glitch: any other horizontal timing is safe.
//  The glitch specifically writes the value of open bus to the register, which will almost always be the upper byte of the address. Writing to this register or the mirror of this register at $2100 according to the desired nametable appears to be a [//forums.nesdev.org/viewtopic.php?p=230434#p230434 functional workaround].
//
//  This produces an occasionally [[Game bugs|visible glitch]] in ''Super Mario Bros.'' when the program writes to PPUCTRL at the end of game logic.
//  It appears to be turning NMI off during game logic and then turning NMI back on once the game logic has finished in order to prevent the NMI handler from being called again before the game logic finishes.
//  Another workaround is to use a software flag to prevent NMI reentry, instead of using the PPU's NMI enable.
//
type PPUCTRL byte

//  (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
func (r PPUCTRL) NameTableBaseAddress() uint16 {
	tmp := r & 0b11
	switch tmp {
	case 0:
		return 0x2000
	case 1:
		return 0x2400
	case 2:
		return 0x2800
	case 3:
		return 0x2c00
	default:
		panic("wrong value")
	}
}

func (r PPUCTRL) VRAMAddressIncrementPerCPUReadOrWriteOfPPUDATA() byte {
	if utils.GetBitFromRight(byte(r), 2) == 0 {
		return 1
	} else {
		return 32
	}
}

func (r PPUCTRL) SpritePatternTableAddressFor88Mode() byte {
	return utils.GetBitFromRight(byte(r), 3)
}

func (r PPUCTRL) BackgroundPatternTableAddress() uint16 {
	if utils.GetBitFromRight(byte(r), 4) == 0 {
		return 0x0000
	} else {
		return 0x1000
	}
}

func (r PPUCTRL) SprintPatternTableAddress() uint16 {
	if r.BackgroundPatternTableAddress() == 0x0000 {
		return 0x1000
	} else {
		return 0x0000
	}
}

func (r PPUCTRL) SpritePatternTableAddress() uint16 {
	if utils.GetBitFromRight(byte(r), 4) == 0 {
		return 0x1000
	} else {
		return 0x0000
	}
}

// Sprite size (0: 8x8 pixels; 1: 8x16 pixels)
func (r PPUCTRL) SpriteSize() byte {
	return utils.GetBitFromRight(byte(r), 5)
}

func (r PPUCTRL) SpriteHeight() byte {
	if r.SpriteSize() == 0 {
		return 8
	} else {
		return 16
	}
}

// PPU master/slave select (0: read backdrop from EXT pins; 1: output color on EXT pins)
func (r PPUCTRL) MasterSlaveSelect() byte {
	return utils.GetBitFromRight(byte(r), 6)
}

func (r PPUCTRL) CanGenerateNMIBreakAtStartOfVerticalBlankingInterval() bool {
	return utils.GetBitFromRight(byte(r), 7) == 1
}
