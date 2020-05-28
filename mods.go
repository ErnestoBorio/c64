package c64

// Just a namespace to contain everything Mod related
type Mods struct {
	raster []ModRasterInfo
}

type ModRasterInfo struct {
	line int
	handler *func(int)
}

// Plug a listener function for a raster line event
func (mods Mods) PlugRasterMod(line int, handler *func(int)) {
	mods.raster = append(mods.raster, ModRasterInfo{line, handler})
}
