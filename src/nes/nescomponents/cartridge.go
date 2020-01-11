package nescomponents

import (
	"encoding/binary"
	"io"
	"os"
)

//Comunication with main BUS
func (cartridge *Cartridge) CpuWrite(address uint16, data byte) bool {
	if address < 0x2000 {
		return cartridge.Mapper.Write(address, data)
	}
	return false
}

func (cartridge *Cartridge) CpuRead(address uint16, data *byte) bool {
	if address < 0x2000 {
		*data = cartridge.Mapper.Read(address)
		if *data != 0 {
			return true
		}
	}
	return false
}

//Comunication  with the second "PPU" BUS
func (cartridge *Cartridge) PpuRead(address uint16, data *byte) bool {
	address = address % 0x4000
	if address < 0x2000 {
		*data = cartridge.Mapper.Read(address)
		if *data != 0 {
			return true
		}
	}
	return false
}

func (cartridge *Cartridge) PpuWrite(address uint16, data byte) bool {
	address = address % 0x4000
	if address < 0x2000 {
		return cartridge.Mapper.Write(address, data)
	}
	return false
}

//The game "la cartouche"
type Cartridge struct {
	Mapper     Mapper
	prg        []byte // PRG-ROM banks
	chr        []byte // CHR-ROM banks
	sram       []byte // Save RAM
	mapperType byte   // mapper type
	mirror     byte   // mirroring mode
	battery    byte   // battery present
}

func NewCartridge(filename string) *Cartridge {
	// call ines struct and load file

	//create ines header struct
	var sHeader InesHeader
	var cartridge Cartridge
	var err error

	// open the game .nes
	file, err := os.Open(filename)
	if err != nil {
		println("call usage and exit")
	}
	defer file.Close()
	// read file header
	sHeader = InesHeader{}
	//insert data inside sheader
	err = binary.Read(file, binary.LittleEndian, &sHeader)
	if err != nil {
		println("call usage and exit")
	}
	//todo check header

	//get mapperId
	cartridge.mapperType = ((sHeader.Mapper2 >> 4) << 4) | (sHeader.Mapper1 >> 4)

	//get mirror mode
	cartridge.mirror = (sHeader.Mapper1 & 1) | (((sHeader.Mapper1 >> 3) & 1) << 1)

	// battery-backed RAM
	cartridge.battery = (sHeader.Mapper1 >> 1) & 1

	// read trainer if present (unused)
	if sHeader.Mapper1&0x04 == 4 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			println("call usage and exit")
		}
	}

	// read prg-rom bank(s)

	cartridge.prg = make([]byte, int(sHeader.PrgRomChunks)*16384) //number mentioned // http://wiki.nesdev.com/w/index.php/INES // http://nesdev.com/NESDoc.pdf (page 28)

	if _, err := io.ReadFull(file, cartridge.prg); err != nil {
		println("call usage and exit")
	}

	// read chr-rom bank(s)

	// provide chr-rom/ram if not in file
	if sHeader.ChrRomChunks == 0 {
		sHeader.ChrRomChunks = 1
	}
	//make funtion allow memory allocation just like malloc
	cartridge.chr = make([]byte, int(sHeader.ChrRomChunks)*8192) //number mentioned // http://wiki.nesdev.com/w/index.php/INES // http://nesdev.com/NESDoc.pdf (page 28)
	if _, err := io.ReadFull(file, cartridge.chr); err != nil {
		println("call usage and exit")
	}

	//sram allocation
	cartridge.sram = make([]byte, 0x2000)

	//load the mapper

	var maperr error
	cartridge.Mapper, maperr = NewMapper(&cartridge)
	if maperr != nil {
		println("call usage and exit")
	}
	return &cartridge
}
