package c64

func (c64 *C64) readMemory(address uint16) byte {
	switch {
	case address <= 0x9FFF:
		return c64.readRAM(address)

	case address <= 0xBFFF: // Basic ROM
		if c64.isBasicOn() {
			return c64.readBasic(address)
		}
		return c64.readRAM(address)

	case address <= 0xCFFF:
		return c64.readRAM(address)

	case address <= 0xDFFF: // IO / Chargen
		if c64.isIOon() {
			return c64.readIO(address)
		}
		if c64.isChargenOn() {
			return c64.readChargen(address)
		}
		return c64.readRAM(address)

	default: // Kernal $E000..$FFFF
		if c64.isKernalOn() {
			return c64.readKernal(address)
		}
		return c64.readRAM(address)
	}
}

func (c64 *C64) writeMemory(address uint16, value byte) {
	if address >= 0xD000 && address <= 0xDFFF && c64.isIOon() {
		c64.writeIO(address, value)
	} else {
		c64.writeRAM(address, value)
	}
}

// Reads RAM
func (c64 *C64) readRAM(address uint16) byte {
	return c64.RAM[address]
}

// Writes to RAM
func (c64 *C64) writeRAM(address uint16, value byte) {
	c64.RAM[address] = value
}

// Reads from Kernal ROM
func (c64 *C64) readKernal(address uint16) byte {
	return ROMKernal[address&0x1FFF]
}

// Reads from BASIC ROM
func (c64 *C64) readBasic(address uint16) byte {
	return ROMBasic[address&0x1FFF]
}

// Read character generator ROM
func (c64 *C64) readChargen(address uint16) byte {
	return ROMChargen[address&0xFFF]
}
