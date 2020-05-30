package c64

// Plug a listener function for a raster line event
func (c64 *C64) PlugRasterMod(line int, handler *func(int)) {
	type rasterInfo struct {
		line int
		handler *func(int)
	}
	c64.Mods.raster = append(c64.Mods.raster, rasterInfo{line, handler})
}
