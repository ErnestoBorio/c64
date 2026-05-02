package c64

/**
 *	<io.go> Handles all writes and reads in the I/O memory mapped registers $D000..$DFFF
 *	Mirrored memory, like $D000..$D03F which is mirrored up to $D3FF, is only ever read and written from the original,
 *	unmirrored addresses, I.E. $D040 will never be read from or written  to through these functions, trying to do so
 *	will actually read from or write to $D000.
 *	Reading and writing to those mirror addresses is up to mods and has to be done through direct access to C64.IO[]
 */

/** TODO
 *	Fijandome en VICE, si lees $D020, los 4 high bits vienen todos en 1.
 *	Sospecho que todos los bits no usados siempre se leen como 1, tener en cuenta.
 *	Como las direcciones tipo $D02F que siempre devuelven $FF.
 */

// Read I/O ports
// address is assumed to be $D000..$DFFF
func (c64 *C64) ReadIO(address uint16) byte {
	if address <= 0xD3FF { // VIC video registers
		address &= 0x3F // $D040..$D3FF are mirrors of $D000..$D03F
		if address >= 0x2F {
			return 0xFF //  $D02F..$D03F Unused bytes, always read as $FF
		} else if address >= 0x20 {
			// Registers $D020..$D02E are colors and only use low 4 bits. High bits always read as 1
			return c64.IO[address] | 0b11110000
		}
		return c64.IO[address]
	} else {
		// Generic IO read WIP
		return c64.IO[address&0xFFF]
	}
}

// Write I/O ports
// address is assumed to be $D000..$DFFF
func (c64 *C64) WriteIO(address uint16, value byte) {
	if address <= 0xD3FF { // VIC video registers
		address &= 0x3F // $D040..$D3FF are mirrors of $D000..$D03F
		switch address {
		case 0x11:
			c64.setRasterCompareHigh(value)
			c64.IO[address] = (c64.IO[address] & 0b10000000) | (value & 0b01111111)
		case 0x12:
			c64.setRasterCompareLow(value)
		case 0x19:
			c64.IO[address] &^= value & vicIRQFlagRaster
			c64.updateVICIRQLine()
		case 0x1A:
			c64.IO[address] = value & vicIRQFlagRaster
			c64.updateVICIRQLine()
		case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E:
			c64.IO[address] = 0b11110000 | value // higher 4 bits are always 1
		case 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F:
			// Unused bytes in $D02F..$D03F ignore writes.
		default:
			c64.IO[address] = value
		}
	} else { // Generic IO write WIP
		c64.IO[address&0xFFF] = value
		// if address == 0xDD00 {
		// 	c64.Vic.setBank(^(value & 0b11)) // Bitwise not of bits 0 & 1 turns 0,1,2,3 into 3,2,1,0
		// }
	}
}
