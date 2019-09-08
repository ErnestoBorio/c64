package c64

import "fmt"

func (c64 *C64) debug(dir string, address word, value byte) {
	hex := fmt.Sprintf("%04X", address)
	fmt.Printf("[%04X] %s $%s\n", c64.CPU.PC, dir, hex)
	return
}

func (c64 *C64) readMemory(address word) byte {
	c64.debug("Read", address, 0)
	switch {
		case address <= 0x9FFF:
			return c64.readRAM(address)
		
		case 0xA000 <= address && address <= 0xBFFF: // Basic ROM
			if c64.RAM[1] & 3 == 3 { // bits 0 & 1 = 11
				return c64.readBasic(address)
			} else {
				return c64.readRAM(address)
			} 
		
		case 0xC000 <= address && address <= 0xCFFF:
			return c64.readRAM(address)
			
		case 0xD000 <= address && address <= 0xDFFF: // IO / Chargen
			if c64.RAM[1] & 3 == 0 { // bits 0 & 1 = 00
				return c64.readRAM(address)
			} else if c64.RAM[1] & 1 > 0 { // bit 2 = 1
				return c64.readIO(address)
			} else { // bit 2 = 0
				return c64.readChargen(address)
			}

		case address >= 0xE000: // Kernal
			if c64.RAM[1] & 2 > 0 { // bit 1 = 1
				return c64.readKernal(address)
			} else { // bit 1 = 0
				return c64.readRAM(address)
			}

		default:
			panic("address not implemented")
	}
}

func (c64 *C64) writeMemory(address word, value byte) {
	c64.debug("Write", address, value)
	switch {
		case address <= 0xCFFF || address >= 0xE000: // All but I/O area
			c64.writeRAM(address, value)
	}

}

// Reads RAM
func (c64 *C64) readRAM(address word) byte {
	return c64.RAM[address]
}

// Writes to RAM
func (c64 *C64) writeRAM(address word, value byte) {
	if address == 0x800 && value != 0 {
		panic("$800 Unused. (Must contain a value of 0 so that the BASIC program can be RUN.) but tried to write non zero.")
	}
	c64.RAM[address] = value
}

// Reads from Kernal ROM
func (c64 *C64) readKernal(address word) byte {
	// Mask address to make it in range of Kernal length
	address &= 0x1FFF
	return ROMs.Kernal[address]
}

// Reads from BASIC ROM
func (c64 *C64) readBasic(address word) byte {
	// Mask address to make it in range of Basic length
	return ROMs.Basic[address & 0x1FFF]
}

// Read I/O ports
func (c64 *C64) writeIO(address word, value byte) {
	panic("Read I/O Not implemented yet")
}

// Write I/O ports
func (c64 *C64) readIO(address word) byte {
	panic("Write I/O Not implemented yet")
}

// Read character generator ROM
func (c64 *C64) readChargen(address word) byte {
	return ROMs.Chargen[address & 0xFFF]
}