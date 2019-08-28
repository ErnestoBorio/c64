package c64

import (
	"cpu6502"
	"os"
	"fmt"
	path "path/filepath"
)

type C64 struct {
	cpu cpu6502.Cpu
}

var ROMs struct {
	Basic   [0x2000] byte // 8kB Basic ROM
	Chargen [0x1000] byte // 4kB Character ROM
	Kernal  [0x2000] byte // 8kB Kernal ROM
}