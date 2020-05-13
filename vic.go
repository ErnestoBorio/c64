package c64

const (
	Character = iota
	Bitmap
)

type VIC struct {
	c64 *C64 // Host system
	Scanline        int // Current rendering scanline
	Cycles2scanline int // How many cycles are left to reach the beginning of next scanline

	VerticalScroll     int  // $D011 bit 0..2: Vertical scroll in pixels
	ScreenCharHeight   int  // $D011 bit 3: 24 | 25
	DisplayEnabled     bool // $D011 bit 4
	GraphicMode        int  // $D011 bit 5: Character | Bitmap
	ExtendedBackGround bool // $D011 bit 6
	RasterIRQpointer   int  // write $D012 and $D011 bit 8: scanline number where to fire IRQ
}

func (vic *VIC) Init(c64 *C64) {
	vic.c64 = c64
}
