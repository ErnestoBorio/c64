package c64

import (
	"os"
	"go64/cpu6502"
)

type C64 struct {
	CPU cpu6502.Cpu
	RAM [0x10000] byte
}

var ROMs struct {
	Basic   [0x2000] byte // 8kB Basic ROM
	Chargen [0x1000] byte // 4kB Character ROM
	Kernal  [0x2000] byte // 8kB Kernal ROM
}

func Init() {
	LoadRoms()
}

// Loads Basic, Chargen and Kernal ROMs
func LoadRoms() {
	var romFile *os.File

	romFile,_ = os.Open("/Users/petruza/Source/etc/go/src/github.com/ernestoborio/c64/roms/basic")
	romFile.Read(ROMs.Basic[:])
	romFile.Close()
	
	romFile,_ = os.Open("/Users/petruza/Source/etc/go/src/github.com/ernestoborio/c64/roms/chargen")
	romFile.Read(ROMs.Chargen[:])
	romFile.Close()

	romFile,_ = os.Open("/Users/petruza/Source/etc/go/src/github.com/ernestoborio/c64/roms/kernal")
	romFile.Read(ROMs.Kernal[:])
	romFile.Close()
}