package c64

import (
	"github.com/ErnestoBorio/cpu6502"
)

const (
	NTSC = 0
	PAL  = 1
	CyclesPerScanline    = 63
	NTSC_cyclesPerSecond = 1022727
	NTSC_cyclesPerFrame  = 16506
	NTSC_scanlines       = 262
	PAL_cyclesPerSecond  = 985248
	PAL_cyclesPerFrame   = 19656
	PAL_scanlines        = 312
)

var Scanlines = [2]int{
	NTSC_scanlines, // [NTSC=0] = 262
	PAL_scanlines,  // [PAL=1]  = 312
}

// Commodore 64 virtual machine
type C64 struct {
	CPU  cpu6502.CPU
	RAM  [0x10000]byte // Whole 64KB of RAM
	IO   [0x1000]byte  // WIP for now just store the bytes raw
	Type int           // NTSC or PAL
	VIC  VIC
	Mods Mods
}

// Make creates and initializes a C64 instance.
func Make(c64type int) *C64 {
	c64 := new(C64)
	c64.Type = c64type // PAL | NTSC
	c64.Init()
	return c64
}

// Initialize the C64 VM instance
func (c64 *C64) Init() {
	c64.CPU = cpu6502.CPU{}
	c64.CPU.Init(c64.readMemory, c64.writeMemory)
	c64.VIC.Init(c64)

	c64.VIC.setScanline(0)
	c64.VIC.Cycles2scanline = CyclesPerScanline

	// Init memory. Mirrored memory has to be set by appropriate function calls. Non mirrored memory can be set directly
	// Initial RAM state
	c64.RAM[0] = 0b00101111 // cpu port direction
	c64.RAM[1] = 0b00110111 // cpu port (bank switch) Basic, IO & Kernel switched on
	c64.RAM[0x2B] = 0x01 // Start address of BASIC program
	c64.RAM[0x2C] = 0x08
	c64.RAM[0x37] =    0 // Pointer to end of BASIC area
	c64.RAM[0x38] = 0xA0
	c64.RAM[0x800] =   0 // Unused (Must contain a value of 0 so that the BASIC program can be RUN)

	// IO Registers, 0xD000 .. 0xDFFF
	c64.IO[0x11] = 0b00011011 // Screen control register #1
}

// Makes C64 set given address to execute next
func (c64 *C64) Jump(address uint16) {
	c64.CPU.PC = address
}

// Whether Basic ROM is switched on or not
func (c64 *C64) isBasicOn() bool {
	return c64.RAM[1]&0b11 == 0b11
}

// Whether Character generator ROM is switched on or not
func (c64 *C64) isChargenOn() bool {
	return (c64.RAM[1]&0b100 == 0) && (c64.RAM[1]&0b11 != 0) // [$0001] bits: 0xx but not 000
}

// Whether IO register bank is switched on or not
func (c64 *C64) isIOon() bool {
	return (c64.RAM[1]&0b100 != 0) && (c64.RAM[1]&0b11 != 0) // [$0001] bits: 1xx but not 100
}

// Whether Kernal ROM is switched on or not
func (c64 *C64) isKernalOn() bool {
	return c64.RAM[1]&0b10 != 0
}

func (c64 *C64) getMaxScanlines() int {
	return Scanlines[c64.Type]
}

// Advance one step, I.E. execute one instruction
func (c64 *C64) Step() {
	// According to http://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt #3.5
	if c64.VIC.scanline >= 48 && c64.VIC.scanline <= 247 &&
		(byte(c64.VIC.scanline & 0b111) == c64.VIC.VerticalScroll()) && c64.VIC.DisplayEnabled() {
			c64.VIC.BadLine = true
	}
	
	c64.VIC.Cycles2scanline -= c64.CPU.Step()

	if c64.VIC.BadLine && c64.VIC.Cycles2scanline <= 40 {
		// Steal the CPU 40 cycles from the end of the scanline (WIP is this right?)
		c64.VIC.Cycles2scanline -= 40
	}
	// WIP: sprites in this scanline also steal CPU cycles, see vic-ii.txt

	if c64.VIC.Cycles2scanline <= 0 {
		c64.VIC.Cycles2scanline += CyclesPerScanline
		c64.VIC.setScanline(c64.VIC.scanline+1)
		
		if c64.VIC.scanline >= c64.getMaxScanlines() {
			c64.VIC.setScanline(0)
		}
	}
}
