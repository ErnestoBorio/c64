package c64

import "fmt"

func (c64 *C64) debug(msg string, address uint16, value byte) {
	hex := fmt.Sprintf("%04X", address)
	fmt.Printf("[%04X] %s $%s\n", c64.CPU.PC, msg, hex)
	return
}

func (c64 *C64) readMemory(address uint16) byte {
	c64.debug("Read ", address, 0)
	switch {
		case address <= 0x9FFF:
			return c64.readRAM(address)
		
		case address <= 0xBFFF: // Basic ROM
			if c64.RAM[1] & 3 == 3 { // bits 0 & 1 = 11
				return c64.readBasic(address)
			} else {
				return c64.readRAM(address)
			} 
		
		case address <= 0xCFFF:
			return c64.readRAM(address)
			
		case address <= 0xDFFF: // IO / Chargen
			if c64.RAM[1] & 3 == 0 { // bits 0 & 1 = 00
				return c64.readRAM(address)
			} else if c64.RAM[1] & 1 > 0 { // bit 2 = 1
				return c64.readIO(address)
			} else { // bit 2 = 0
				return c64.readChargen(address)
			}

		default: // Kernal $E000..$FFFF
			if c64.RAM[1] & 2 > 0 { // bit 1 = 1
				return c64.readKernal(address)
			} else { // bit 1 = 0
				return c64.readRAM(address)
			}
	}
}

func (c64 *C64) writeMemory(address uint16, value byte) {
	c64.debug("Write", address, value)
	switch {
		case address >= 0xD000 && address <= 0xDFFF:
			if c64.RAM[1] & 1 > 0 {
				c64.writeIO(address, value)
			} else {
				c64.writeRAM(address, value)
			}
		default:
			c64.writeRAM(address, value)
	}
}

// Reads RAM
func (c64 *C64) readRAM(address uint16) byte {
	return c64.RAM[address]
}

// Writes to RAM
func (c64 *C64) writeRAM(address uint16, value byte) {
	if address == 0x800 && value != 0 {
		panic("$800 Unused. (Must contain a value of 0 so that the BASIC program can be RUN.) but tried to write non zero.")
	}
	c64.RAM[address] = value
}

// Reads from Kernal ROM
func (c64 *C64) readKernal(address uint16) byte {
	// Mask address to make it in range of Kernal length
	address &= 0x1FFF
	return ROMs.Kernal[address]
}

// Reads from BASIC ROM
func (c64 *C64) readBasic(address uint16) byte {
	// Mask address to make it in range of Basic length
	return ROMs.Basic[address & 0x1FFF]
}

// Read character generator ROM
func (c64 *C64) readChargen(address uint16) byte {
	return ROMs.Chargen[address & 0xFFF]
}