package c64

/**
	File formats reading and parsing
*/

import (
	"io"
	"os"
	"fmt"
	"strings"
	"strconv"
)

type File struct {
	format string
	err error
	data []Filedata // Can be multiple if it's a T64 container
}

type Filedata struct {
	name string
	address uint16
	data []byte
}

// Loads a formatted file and returns the raw file(s) content and load address(es)
func ReadFile(file *os.File) File {
	// Get file size and name
	info, err := file.Stat()
	var filename string
	var size int64

	if err == nil {
		size = info.Size()
		filename = info.Name()
	} else {
		// This shouldn't happen, why would getting file info fail?
		size = 100*1024 // Hopefully allocate enough space to read whole file
		filename = "<no filename>"
	}
	// Read a 32 byte header that will hint what type of file it is
	headerBuf := make([]byte, 32)

	_, err = file.Read(headerBuf)
	if err != nil {
		return File {
			err: fmt.Errorf("Can't read file %s", filename)}
	}
	// Make header all uppercase to make case insensitive string comparison easier
	header := strings.ToUpper( string(headerBuf))
	
	if header[:8] == "C64FILE\x00" { /** P00 **/
		// It's a P00 file (file extension could be P** for * = digit)
		panic("P00 format not implemented yet.")
	} else { /** T64 **/
		// T64 header is not standard, at least should start with "C64" and include "tape" somewhere.
		if header[:3] == "C64" && strings.Contains( header, "TAPE") {
			// It's a T64 file
			panic("T64 format not implemented yet.")
		} else { /** PRG? **/
			// Probably a PRG file
			_,err = file.Seek(0, io.SeekStart)
			buffer := make([]byte, size)
			_,err = file.Read(buffer)
			if err != nil {
				return File {
					err: fmt.Errorf("Couldn't allocate buffer of size %d bytes to read file %s.",
					size, filename )}
			}
			machineCode, address := parseBasicLoader(buffer)
			return File {
				format: "PRG",
				err: nil,
				data: []Filedata { Filedata {
						name: filename, // Since PRG format doesn't have a C64 filename, return host OS filename
						address: address,
						data: machineCode }}}
		}
	}
}

// Parse a PRG raw file in search of a Basic loader that uses the SYS command 
// to jump to the address of the beginning of the machine language program.
func parseBasicLoader(data []byte) (machineCode []byte, address uint16) {
	// if address == 0x0801 most likely a BASIC loader
	// TODO: now blindly assuming program is single line with SYS, should correctly parse BASIC
	baseAddress := (uint16(data[1]) <<8) | uint16(data[0])
	// Skip first 6 bytes: load address, next line pointer and current line number
	if data[6] == 0x9E {
		// $9E token is Basic SYS command
		line := data[7:]
		strAddress := ""
		for {
			bite := line[0] // read next byte
			if bite == 0 {
				break
			}
			if '0' <= bite && bite <= '9' { // concat digit characters to form address
				strAddress += string(bite)
			}
			line = line[1:] // consume the byte just read
		}
		if strAddress == "" {
			panic("Weird BASIC Loader doesn't just 10 SYS 2064. (a)")
		}
		// discard Basic loader and only keep machine code program
		address,_ := strconv.Atoi(strAddress)
		delta := 2 + uint16(address) - baseAddress // offset from start of Basic program to start of machine code
		data = data[delta:]
		return data, uint16(address)
	} else {
		panic("Weird BASIC Loader doesn't just 10 SYS 2064. (b)")
	}
}