package c64

import (
	"github.com/Drean64/cpu6502"
)

type TVMode int

const (
	NTSC TVMode = iota
	PAL
)

const (
	CyclesPerScanline   = 63
	NTSCCyclesPerSecond = 1022727
	NTSCScanlines       = 262
	PALCyclesPerSecond  = 985248
	PALScanlines        = 312
	IRQCycles           = 7
)

type IRQSource byte

const (
	IRQSourceVIC IRQSource = 1 << iota
	IRQSourceCIA
)

// C64 models a Commodore 64 virtual machine
type C64 struct {
	CPU        cpu6502.CPU
	RAM        [0x10000]byte // Whole 64KB of RAM
	IO         [0x1000]byte  // @todo WIP for now just store the bytes raw
	Type       TVMode        // NTSC or PAL
	Vic        VIC
	irqSources IRQSource
}

// Make creates a C64 instance.
func Make(c64type TVMode) *C64 {
	c64 := &C64{Type: c64type} // PAL | NTSC
	c64.Init()
	return c64
}

func (c64 *C64) NTSC() bool {
	return c64.Type == NTSC
}

// CyclesPerFrame returns the number of CPU cycles in one video frame for the
// current video standard.
func (c64 *C64) CyclesPerFrame() int {
	return c64.totalScanlines() * CyclesPerScanline
}

func (c64 *C64) SoftReset() {
	c64.CPU.Reset()
}

func (c64 *C64) HardReset() {
	c64.Init()
	c64.CPU.Reset()
}

// Initialize the C64 VM instance
func (c64 *C64) Init() {
	c64.CPU.Init(c64.readMemory, c64.writeMemory)

	// Init memory. Mirrored memory has to be set by appropriate function calls. Non mirrored memory can be set directly
	// Initial RAM state
	c64.RAM[0] = 0b00101111 // cpu port direction
	c64.RAM[1] = 0b00110111 // cpu port (bank switch) Basic, IO & Kernel switched on
	c64.RAM[0x2B] = 0x01    // Start address of BASIC program
	c64.RAM[0x2C] = 0x08
	c64.RAM[0x37] = 0 // Pointer to end of BASIC area
	c64.RAM[0x38] = 0xA0
	c64.RAM[0x800] = 0     // Unused (Must contain a value of 0 so that the BASIC program can be RUN)
	c64.RAM[0xFFFC] = 0xE2 // Reset vector low byte
	c64.RAM[0xFFFD] = 0xFC // Reset vector high byte ($FCE2)

	// IO Registers, 0xD000 .. 0xDFFF
	c64.IO[0x11] = 0b00011011  // Screen control register #1
	c64.IO[0xD00] = 0b00111011 // VIC bank selection, RS232 and serial ports
	c64.VICInit()
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

func (c64 *C64) AssertIRQ(source IRQSource) {
	c64.irqSources |= source
}

func (c64 *C64) ClearIRQ(source IRQSource) {
	c64.irqSources &^= source
}

func (c64 *C64) IRQLine() bool {
	return c64.irqSources != 0
}

// advanceTiming advances machine subsystems using the number of CPU cycles just
// consumed by the current instruction.
func (c64 *C64) advanceTiming(cycles int) {
	c64.advanceVIC(cycles)
	// CIA, IRQ, and other machine timing will be advanced here.
}

// Step advances the whole C64 by one CPU instruction and returns the number of
// CPU cycles consumed.
func (c64 *C64) Step() int {
	if c64.IRQLine() && !c64.CPU.Status.NoInterrupt {
		c64.CPU.IRQ()
		c64.advanceTiming(IRQCycles)
		return IRQCycles
	}

	cyclesAdvanced := c64.CPU.Step()
	c64.advanceTiming(cyclesAdvanced)
	return cyclesAdvanced
}

// RunCycles advances the C64 until at least the requested number of CPU cycles
// have elapsed. The returned value is the actual number of cycles consumed,
// which may be greater than requested because execution only stops between
// instructions.
func (c64 *C64) RunCycles(cycles int) int {
	cyclesAdvanced := 0
	for cyclesAdvanced < cycles {
		cyclesAdvanced += c64.Step()
	}
	return cyclesAdvanced
}
