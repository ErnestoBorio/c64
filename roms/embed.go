package main

import (
	"fmt"
	"os"
)

/*
Generates Go source files with byte arrays with the content of the ROM binaries.
*/
func main() {
	var path, outPath string

	// Arg 1 is ROM binaries path
	if len(os.Args) >= 2 {
		path = os.Args[1]
	} else {
		path = "./"
	}

	// Arg 2 is .go files output path
	if len(os.Args) >= 3 {
		outPath = os.Args[2]
	} else {
		outPath = "../"
	}

	roms := [...]struct {
		name    string
		desc    string
		size    int
		varName string
	}{
		{name: "basic", size: 0x2000, varName: "ROMBasic", desc: "// BASIC interpreter ROM, 8KB"},
		{name: "chargen", size: 0x1000, varName: "ROMChargen", desc: "// Character generator ROM, 4KB"},
		{name: "kernal", size: 0x2000, varName: "ROMKernal", desc: "// Kernal ROM, 8KB"},
	}

	for _, rom := range roms {
		romFile, err := os.Open(path + rom.name)
		if err != nil {
			panic(path + rom.name + " ROM could not be opened")
		}
		buffer := make([]byte, rom.size)
		read, err := romFile.Read(buffer)
		if err != nil {
			panic(path + rom.name + " ROM could not be read")
		}
		size := rom.size
		stat, err := romFile.Stat()
		romFile.Close()
		if err == nil {
			size = int(stat.Size())
		}
		if read != size || read != rom.size {
			panic(fmt.Sprintf("File size (%d), ROM size (%d) and bytes read (%d) don't match", size, rom.size, read))
		}

		outf, err := os.Create(outPath + "rom." + rom.name + ".go")
		if err != nil {
			panic(outPath + "rom." + rom.name + ".go file could not be created")
		}
		out := fmt.Sprintf("package c64\n\n%s\nvar %s = [0x%X] byte {\n", rom.desc, rom.varName, rom.size)
		outf.WriteString(out)
		for i, bite := range buffer {
			outf.WriteString(fmt.Sprintf("0x%02X,", bite))
			if (i+1)%16 == 0 {
				outf.WriteString("\n")
			}
		}
		outf.WriteString("}\n")
		outf.Close()
	}
}
