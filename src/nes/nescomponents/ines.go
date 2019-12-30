package nescomponents

//game loader that the cartridge will use
const INESFileMagic = 0x1a53454e

//ines format header
type InesHeader struct {
	// PrgRomChunks byte // number of PRG-ROM banks (16KB each)
	// ChrRomChunks byte // number of CHR-ROM banks (8KB each)
	// Mapper1      byte // control bits
	// Mapper2      byte // control bits
	// PrgRamSize   byte // PRG-RAM size (x 8KB)
	// TvSystem1    byte
	// TvSystem2    byte

	Magic        uint32  // iNES magic number
	PrgRomChunks byte    // number of PRG-ROM banks (16KB each)
	ChrRomChunks byte    // number of CHR-ROM banks (8KB each)
	Mapper1      byte    // control bits
	Mapper2      byte    // control bits
	PrgRamSize   byte    // PRG-RAM size (x 8KB)
	_            [7]byte // unused padding
}
