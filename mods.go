package c64

type ReaderMod interface {
	Read(address uint16, value byte)
}

type WriterMod interface {
	Write(address uint16, value byte)
}

type ReaderModInfo struct {
	fromAddress uint16
	toAddress uint16
	mod *ReaderMod
}

type WriterModInfo struct {
	fromAddress uint16
	toAddress uint16
	mod *WriterMod
}

// Holds lists of listener Mods
type Mods struct {
	Readers []ReaderModInfo
	Writers []WriterModInfo
}

// Plug a Mod that listens to memory reads in this specific address interval
// Enter address twice for a single address
func (c64 *C64) plugReaderMod(fromAddress uint16, toAddress uint16, mod *ReaderMod) {
	c64.Mods.Readers = append(c64.Mods.Readers, ReaderModInfo { fromAddress, toAddress, mod	})
}

// Plug a Mod that listens to memory writes in this specific address interval
// Enter address twice for a single address
func (c64 *C64) plugWriterMod(fromAddress uint16, toAddress uint16, mod *WriterMod) {
	c64.Mods.Writers = append(c64.Mods.Writers, WriterModInfo { fromAddress, toAddress, mod	})
}

