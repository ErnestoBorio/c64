package c64

/*
	<io.go> Handles all writes and reads in the I/O memory mapped registers $D000..$DFFF

	Mirrored memory, like $D000..$D03F which is mirrored up to $D3FF, is only ever read and written from the original,
	unmirrored addresses, I.E. $D040 will never be read or written through these functions.
	Reading and writing to those mirror addresses is up to mods and has to be done through direct access to C64.IO[]
*/

// Read I/O ports
func (c64 *C64) readIO(address uint16) byte {
	if(address >= 0xD000 && address <= 0xD3FF) {
		// WIP test pending
		address &= 0x3F // de-mirror address, registers $D000..$D03F are repeated up to $D3FF
		if address >= 0x2F {
			return 0xFF // Unused bytes
		}
		return c64.IO[address]
	} else {
		// Generic IO read WIP
		// Transpose address into space $0..$FFF of IO bank
		return c64.IO[address & 0xFFF]
	}
}

// Write I/O ports
func (c64 *C64) writeIO(address uint16, value byte) {
	// VIC video registers
	if(address >= 0xD000 && address <= 0xD3FF) {
		// WIP test pending
		address &= 0x3F // de-mirror address, registers $D000..$D03F are repeated up to $D3FF
		c64.IO[address] = value

		switch address {
			case 0xD011: c64.writeD011(value)
			case 0xD012:
				c64.VIC.RasterIRQpointer &= 0b100000000 // Keep only bit 8
				c64.VIC.RasterIRQpointer |= int(value)  // Set all other bits as written
			// case $D02F..$D03F: unused bytes, write is just ignored
		}
	} else {
		// Generic IO write WIP
		// Transpose address into space $0..$FFF of IO bank
		c64.IO[address & 0xFFF] = value
	}
}

// WIP test pending
func (c64 *C64) writeD011(value byte) {
	c64.VIC.VerticalScroll = int(value & 0b111) // bits 0..2

	if value & 0b1000 == 0 { // bit 3
		c64.VIC.ScreenCharHeight = 24
	} else {
		c64.VIC.ScreenCharHeight = 25
	}

	if value & 0b10000 == 0 { // bit 4
		// 0 = Screen off, complete screen is covered by border
		c64.VIC.DisplayEnabled = false
	} else {
		// 1 = Screen on, normal screen contents are visible
		c64.VIC.DisplayEnabled = true
	}

	if value & 0b100000 == 0 { // bit 5
		c64.VIC.GraphicMode = Character
	} else {
		c64.VIC.GraphicMode = Bitmap
	}

	if value & 0b1000000 == 0 { // bit 6: Extended background mode
		c64.VIC.ExtendedBackGround = false
	} else {
		c64.VIC.ExtendedBackGround = true
	}

	if value & 0b10000000 == 0 { // bit 7: Bit 8 of the scanline number where to fire IRQ
		c64.VIC.RasterIRQpointer &= 0b011111111 // Unset bit 8
		//                      bits: 876543210
	} else {
		c64.VIC.RasterIRQpointer |= 0b100000000 // Set bit 8
	}
}
