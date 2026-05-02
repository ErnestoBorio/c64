package c64

// WIP not tested

const (
	Character = iota
	Bitmap
)

const (
	vicIRQFlagRaster = 0b00000001
)

type VIC struct {
	scanline           int
	cyclesIntoScanline int
	rasterCompare      int
}

// totalScanlines returns the number of raster lines per frame for the current
// video standard.
func (c64 *C64) totalScanlines() int {
	if c64.NTSC() {
		return NTSCScanlines
	}
	return PALScanlines
}

func (c64 *C64) syncRasterRegisters() {
	scanline := c64.Vic.scanline
	c64.IO[0x12] = byte(scanline)
	c64.IO[0x11] &= 0b01111111
	c64.IO[0x11] |= byte(scanline>>1) & 0b10000000
}

func (c64 *C64) setRasterCompareLow(value byte) {
	c64.Vic.rasterCompare = (c64.Vic.rasterCompare & 0x100) | int(value)
}

func (c64 *C64) setRasterCompareHigh(value byte) {
	c64.Vic.rasterCompare = (c64.Vic.rasterCompare & 0x0FF) | (int(value&0b10000000) << 1)
}

func (c64 *C64) updateVICIRQLine() {
	if c64.IO[0x19]&c64.IO[0x1A]&vicIRQFlagRaster != 0 {
		c64.AssertIRQ(IRQSourceVIC)
		return
	}
	c64.ClearIRQ(IRQSourceVIC)
}

func (c64 *C64) triggerRasterIRQ() {
	c64.IO[0x19] |= vicIRQFlagRaster
	c64.updateVICIRQLine()
}

func (c64 *C64) advanceVIC(cycles int) {
	if cycles <= 0 {
		return
	}

	c64.Vic.cyclesIntoScanline += cycles
	for c64.Vic.cyclesIntoScanline >= CyclesPerScanline {
		c64.Vic.cyclesIntoScanline -= CyclesPerScanline
		c64.Vic.scanline++
		if c64.Vic.scanline >= c64.totalScanlines() {
			c64.Vic.scanline = 0
		}
		if c64.Vic.scanline == c64.Vic.rasterCompare {
			c64.triggerRasterIRQ()
		}
	}

	c64.syncRasterRegisters()
}

func (c64 *C64) VICInit() {
	c64.Vic.scanline = 0
	c64.Vic.cyclesIntoScanline = 0
	c64.Vic.rasterCompare = 0
	c64.IO[0x11] = 0b00011011 // vertical scroll = 3, height = 25 rows, screen on, text mode, extended bg off
	c64.IO[0x19] = 0
	c64.IO[0x1A] = 0
	c64.syncRasterRegisters()
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
