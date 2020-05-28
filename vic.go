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

// $D011 bit 0..2: Vertical scroll in pixels
func (c64 *C64) VerticalScroll() byte {
	return c64.IO[0x11] & 0b111
}

// $D011 bit 3: 24 | 25
func (c64 *C64) ScreenCharHeight() byte {
	if c64.IO[0x11] & 0b1000 == 0 {
		return 24
	}
	return 25
}

// $D011 bit 4
// 0 = Screen off, complete screen is covered by border
// 1 = Screen on, normal screen contents are visible
func (c64 *C64) DisplayEnabled() bool {
	if c64.IO[0x11] & 0b10000 != 0 {
		return true
	}
	return false
}


// $D011 bit 5 // Character | Bitmap
func (c64 *C64) GraphicMode() int {
	if c64.IO[0x11] & 0b100000 == 0 {
		return Character
	}
	return Bitmap
}

// $D011 bit 6: Extended background mode
func (c64 *C64) ExtendedBackGround() bool {
	if c64.IO[0x11] & 0b1000000 == 0 { 
		return false
	}
	return true
}

// Keep memory in sync with the new scanline number
func (c64 *C64) setScanline(newScanline int) {
	c64.scanline = newScanline
	c64.IO[0x12] = byte(newScanline) // Get lower 8 bits
	c64.IO[0x11] &= 0b01111111 // Clear bit 7
	c64.IO[0x11] |= (byte(newScanline & 0b100000000) >>1) // Get scanline bit 8 and push it as bit 7 in $D011

	// According to http://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt #3.5
	if c64.scanline >= 48 && c64.scanline <= 247 &&
		(byte(c64.scanline & 0b111) == c64.VerticalScroll()) &&
			c64.DisplayEnabled() {
				c64.BadLine = true
	} else {
		c64.BadLine = false
	}

func (vic *VIC) setRasterIRQline(line int) {
}
