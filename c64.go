package c64

import (
	"go64/cpu6502"
)

// Commodore 64 virtual machine
type C64 struct {
	CPU cpu6502.CPU
	RAM [0x10000]byte // Whole 64KB of RAM
	IO  [0x1000]byte  // TODO for now just store the bytes raw
}

// Make creates and initializes a C64 instance.
func Make() *C64 {
	c64 := new(C64)
	c64.Init()
	return c64
}

// Initialize the C64 VM instance
func (c64 *C64) Init() {
	c64.CPU = cpu6502.CPU{}
	c64.CPU.Init(c64.readMemory, c64.writeMemory)

	// Initial RAM state
	c64.RAM[0] = 0x2F   // cpu port direction: 00101111
	c64.RAM[1] = 0x37   // cpu port (bank switch): 00110111 (Basic, IO & Kernel)
	c64.RAM[0x0800] = 0 // Unused (Must contain a value of 0 so that the BASIC program can be RUN)

	// IO Registers, 0xD000 .. 0xDFFF
	c64.IO[0x16] = 0xC8 // Screen control register #2: 11001000
}

// Makes C64 set given address to execute next
func (c64 *C64) Jump(address uint16) {
	c64.CPU.PC = address
}

// Whether Basic ROM is switched on or not
func (c64 *C64) isBasicOn() bool {
	return c64.RAM[1]&3 == 3 // [$0001] bits: x11
}

// Whether Character generator ROM is switched on or not
func (c64 *C64) isChargenOn() bool {
	return c64.RAM[1]&7 >= 1 && c64.RAM[1]&7 <= 3 // [$0001] bits: 0xx but not 000
}

// Whether IO register bank is switched on or not
func (c64 *C64) isIOon() bool {
	return c64.RAM[1]&7 >= 5 && c64.RAM[1]&7 <= 7// [$0001] bits: 1xx but not 100
}

// Whether Kernal ROM is switched on or not
func (c64 *C64) isKernalOn() bool {
	return c64.RAM[1]&2 > 0 // [$0001] bits: x1x
}
