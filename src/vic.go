package c64

// WIP not tested

const (
	Character = iota
	Bitmap
)

type VIC struct {
}

// $D011 bit 0..2: Vertical scroll in pixels
func (c64 *C64) VerticalScroll() byte {
	return c64.IO[0x11] & 0b111
}

// $D011 bit 3: 24 | 25
func (c64 *C64) ScreenCharHeight() byte {
	if c64.IO[0x11]&0b1000 == 0 {
		return 24
	}
	return 25
}

// $D011 bit 4
// 0 = Screen off, complete screen is covered by border
// 1 = Screen on, normal screen contents are visible
func (c64 *C64) DisplayEnabled() bool {
	return c64.IO[0x11]&0b10000 != 0
}

// $D011 bit 5 // Character | Bitmap
func (c64 *C64) GraphicMode() int {
	if c64.IO[0x11]&0b100000 == 0 {
		return Character
	}
	return Bitmap
}

// $D011 bit 6: Extended background mode
func (c64 *C64) ExtendedBackGround() bool {
	return c64.IO[0x11]&0b1000000 != 0
}
