package c64

import (
	"os"
	"go64/cpu6502"
)

type word = uint16

// Commodore 64 virtual machine
type C64 struct {
	CPU cpu6502.CPU
	RAM [0x10000] byte
}

// Package-stored ROMs for all instantiated C64s
var ROMs struct {
	Basic   [0x2000] byte // 8kB Basic ROM
	Chargen [0x1000] byte // 4kB Character ROM
	Kernal  [0x2000] byte // 8kB Kernal ROM
}

// Loads Basic, Chargen and Kernal ROMs
func LoadRoms() {
	romFile,_ := os.Open("/Users/petruza/Source/etc/go/src/go64/c64/roms/basic")
	romFile.Read(ROMs.Basic[:])
	romFile.Close()
	
	romFile,_ = os.Open("/Users/petruza/Source/etc/go/src/go64/c64/roms/chargen")
	romFile.Read(ROMs.Chargen[:])
	romFile.Close()

	romFile,_ = os.Open("/Users/petruza/Source/etc/go/src/go64/c64/roms/kernal")
	romFile.Read(ROMs.Kernal[:])
	romFile.Close()
}

// Initialize package
func Init() {
	LoadRoms()
}

// Initialize the C64 VM instance
func (c64 *C64) Init() {
	c64.CPU = cpu6502.CPU{}
	c64.CPU.Init()
	// hook the 6502 CPU with the C64 memory
	c64.CPU.ReadMemory = c64.readMemory
	c64.CPU.WriteMemory = c64.writeMemory

	// Initial memory state
	c64.RAM[0] = 0x2F // 00101111
	c64.RAM[1] = 0x37 // 00110111
	c64.RAM[0x800] = 0 // Unused. (Must contain a value of 0 so that the BASIC program can be RUN.)
}

func (c64 *C64) Run() {
	c64.CPU.Reset()
	// Main loop
	for ;; {
		/* cycles :=  */c64.CPU.Step()

	}
}