package c64

// WIP not tested
/*
	The VIC struct is not meant to be an independent entity, but rather a namespace of sorts to keep VIC-II related
	attributes and methods and not pollute C64
*/

const (
	Character = iota
	Bitmap
)

type VIC struct {
	// Pointers to the C64 memory:
	RAM             *[0x10000]byte
	IO              *[0x1000] byte
	Scanline        int // Current rendering scanline
	Cycles2scanline int // How many cycles are left to reach the beginning of next scanline
}

func (vic *VIC) Init(c64 *C64) {
	vic.RAM = &c64.RAM
	vic.IO  = &c64.IO
}

// $D011 bit 0..2: Vertical scroll in pixels
func (vic *VIC) VerticalScroll() byte {
	return vic.IO[0x11] & 0b111
}

// $D011 bit 3: 24 | 25
func (vic *VIC) ScreenCharHeight() byte {
	if vic.IO[0x11] & 0b1000 == 0 {
		return 24
	}
	return 25
}

// $D011 bit 4
// 0 = Screen off, complete screen is covered by border
// 1 = Screen on, normal screen contents are visible
func (vic *VIC) DisplayEnabled() bool {
	if vic.IO[0x11] & 0b10000 != 0 {
		return true
	}
	return false
}


// $D011 bit 5 // Character | Bitmap
func (vic *VIC) GraphicMode() int {
	if vic.IO[0x11] & 0b100000 == 0 {
		return Character
	}
	return Bitmap
}

// $D011 bit 6: Extended background mode
func (vic *VIC) ExtendedBackGround() bool {
	if vic.IO[0x11] & 0b1000000 == 0 { 
		return false
	}
	return true
}

// $D012 lower 8 bits of the scanline number where to fire IRQ
// $D011 bit 7: Bit 8 of previous number
func (vic *VIC) RasterIRQpointer() int {
	return int((vic.IO[0x11] & 0b10000000) << 1) | int(vic.IO[0x12])
}

func (vic *VIC) SetScanline(scanline int) {
	vic.IO[0x12] = byte(scanline) // get lower 8 bits
	vic.IO[0x11] &= 0b01111111 // clear bit 7
	vic.IO[0x11] |= byte(scanline & 0b10000000)
}