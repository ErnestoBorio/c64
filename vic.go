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
	scanline        int // Current rendering scanline
	BadLine         bool // Whether the current scanline is a bad line
	Cycles2scanline int // How many cycles are left to reach the beginning of next scanline
	rasterIRQline   int // Next scanline at which fire an IRQ
}

func (vic *VIC) Init(c64 *C64) {
	vic.RAM = &c64.RAM
	vic.IO  = &c64.IO
	vic.BadLine = false
	vic.rasterIRQline = 0
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

// Keep memory in sync with the new scanline number
func (vic *VIC) setScanline(newScanline int) {
	vic.scanline = newScanline
	vic.IO[0x12] = byte(newScanline) // Get lower 8 bits
	vic.IO[0x11] &= 0b01111111 // Clear bit 7
	vic.IO[0x11] |= (byte(newScanline & 0b100000000) >>1) // Get scanline bit 8 and push it as bit 7 in $D011
}

func (vic *VIC) Scanline() int {
	return vic.scanline
}

func (vic *VIC) setRasterIRQline(line int) {
}

func (vic *VIC) getRasterIRQline() int {
	return vic.rasterIRQline
}