package c64

/*
Handles all writes and reads in the I/O memory mapped registers $D000..$DFFF
*/

// Read I/O ports
func (c64 *C64) readIO(address uint16) byte {
	return c64.IO[address&0xFFF]
}

// Write I/O ports
func (c64 *C64) writeIO(address uint16, value byte) {
	c64.IO[address&0xFFF] = value
}
