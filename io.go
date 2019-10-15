package c64

/*
Handles all writes and reads in the I/O memory mapped registers $D000..$DFFF
*/

// Read I/O ports
func (c64 *C64) readIO(address uint16) byte {
	return 0
}

// Write I/O ports
func (c64 *C64) writeIO(address uint16, value byte) {

}
