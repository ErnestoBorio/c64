package c64

import (
	"crypto/sha512"
	"encoding/hex"
	"os"
	"testing"
)

func TestInstanciation(t *testing.T) {
	c64 := Make(NTSC)
	c64.Init()
	c64.CPU.Reset()
}

func TestRoms(t *testing.T) {
	const kernalSha512 = "c6ef3021ab08151bd93399ed8c2a97896cb1fb1e2820865622ce1a9169242b48315c71f58aba26b7f720e872b9b941e651378a7a3d99218e1c104d55e412d25c"
	const basicSha512 = "683c3ca9bf14d71b988b35381843a9d8f4e083254b45f2f2a27c1a1a3508090134156171de56604fe3be65ae0d4efbd24e0923d7f6c19f3449c25212711e2320"
	const chargenSha512 = "63b583016b46d5c7af11dd1b57fc721370fbfaf50b7ef2c15a7f40dd420dcbcaf7d3e173d7d42b88e6ed45336bd0cae67b921a15a6b8ee6d0a0925a3ba7caab2"

	kernalHash := sha512.Sum512(ROMKernal)
	if kernalSha512 != hex.EncodeToString(kernalHash[:]) {
		t.Errorf("Kernal ROM hash failed (`roms/kernal`)")
	}
	basicHash := sha512.Sum512(ROMBasic)
	if basicSha512 != hex.EncodeToString(basicHash[:]) {
		t.Errorf("Basic ROM hash failed (`roms/basic`)")
	}
	chargenHash := sha512.Sum512(ROMChargen)
	if chargenSha512 != hex.EncodeToString(chargenHash[:]) {
		t.Errorf("Chargen ROM hash failed (`roms/chargen`)")
	}
}

func TestRunC64(t *testing.T) {
	c64 := Make(NTSC)
	c64.Init()

	currentPath, err := os.Getwd()
	filePath := currentPath + "/../files/Montezuma's Revenge (1984)(Parker Brothers).prg"
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("error reading file: " + filePath)
	}
	jumpTo, err := c64.LoadPRG(file)
	if err != nil {
		t.Errorf("error loading PRG: %v", err)
	}
	c64.Jump(jumpTo)
}

func loadTestProgram(c64 *C64, address uint16, program []byte) {
	copy(c64.RAM[address:], program)
	c64.Jump(address)
}

func TestRunCyclesZeroDoesNothing(t *testing.T) {
	c64 := Make(NTSC)
	c64.Init()
	loadTestProgram(c64, 0x0800, []byte{0xEA}) // NOP

	startPC := c64.CPU.PC
	ran := c64.RunCycles(0)

	if ran != 0 {
		t.Fatalf("RunCycles(0) consumed %d cycles, want 0", ran)
	}
	if c64.CPU.PC != startPC {
		t.Fatalf("RunCycles(0) advanced PC to $%04X, want $%04X", c64.CPU.PC, startPC)
	}
}

func TestRunCyclesStopsAtInstructionBoundaries(t *testing.T) {
	c64 := Make(NTSC)
	c64.Init()
	loadTestProgram(c64, 0x0800, []byte{
		0xEA,             // NOP      2 cycles
		0xEA,             // NOP      2 cycles
		0x4C, 0x00, 0x08, // JMP $0800 3 cycles
	})

	ran := c64.RunCycles(1)
	if ran != 2 {
		t.Fatalf("first RunCycles consumed %d cycles, want 2", ran)
	}
	if c64.CPU.PC != 0x0801 {
		t.Fatalf("first RunCycles left PC at $%04X, want $0801", c64.CPU.PC)
	}

	ran = c64.RunCycles(2)
	if ran != 2 {
		t.Fatalf("second RunCycles consumed %d cycles, want 2", ran)
	}
	if c64.CPU.PC != 0x0802 {
		t.Fatalf("second RunCycles left PC at $%04X, want $0802", c64.CPU.PC)
	}

	ran = c64.RunCycles(3)
	if ran != 3 {
		t.Fatalf("third RunCycles consumed %d cycles, want 3", ran)
	}
	if c64.CPU.PC != 0x0800 {
		t.Fatalf("third RunCycles left PC at $%04X, want $0800", c64.CPU.PC)
	}
}

func TestRunCyclesOvershootsTargetWhenNeeded(t *testing.T) {
	c64 := Make(NTSC)
	c64.Init()
	loadTestProgram(c64, 0x0800, []byte{
		0xEA,             // NOP      2 cycles
		0xEA,             // NOP      2 cycles
		0x4C, 0x00, 0x08, // JMP $0800 3 cycles
	})

	ran := c64.RunCycles(3)
	if ran != 4 {
		t.Fatalf("RunCycles(3) consumed %d cycles, want 4 due to instruction-boundary overshoot", ran)
	}
	if c64.CPU.PC != 0x0802 {
		t.Fatalf("RunCycles(3) left PC at $%04X, want $0802", c64.CPU.PC)
	}
}

func TestStepAndRunCyclesAgree(t *testing.T) {
	stepC64 := Make(NTSC)
	stepC64.Init()
	loadTestProgram(stepC64, 0x0800, []byte{
		0xEA,             // NOP
		0xEA,             // NOP
		0x4C, 0x00, 0x08, // JMP $0800
	})

	runC64 := Make(NTSC)
	runC64.Init()
	loadTestProgram(runC64, 0x0800, []byte{
		0xEA,
		0xEA,
		0x4C, 0x00, 0x08,
	})

	stepCycles := stepC64.Step() + stepC64.Step() + stepC64.Step()
	runCycles := runC64.RunCycles(stepCycles)

	if runCycles != stepCycles {
		t.Fatalf("RunCycles consumed %d cycles, want %d", runCycles, stepCycles)
	}
	if runC64.CPU.PC != stepC64.CPU.PC {
		t.Fatalf("RunCycles left PC at $%04X, want $%04X", runC64.CPU.PC, stepC64.CPU.PC)
	}
}
