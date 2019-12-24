package nescomponents

//game loader that the cartridge will use
const INESFileMagic = 0x1a53454e

//ines format header
type InesHeader struct {
	prgRomChunks byte // number of PRG-ROM banks (16KB each)
	chrRomChunks byte // number of CHR-ROM banks (8KB each)
	mapper1      byte // control bits
	mapper2      byte // control bits
	prgRamSize   byte // PRG-RAM size (x 8KB)
	tvSystem1    byte
	tvSystem2    byte
}
