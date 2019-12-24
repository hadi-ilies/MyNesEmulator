package nescomponents

import (
	"encoding/binary"
	"io"
	"os"
)

//Comunication with main BUS
func (cartridge *Cartridge) CpuWrite(address uint16, data byte) bool {
	return false
}

func (cartridge *Cartridge) CpuRead(address uint16) bool {
	return false
}

//Comunication  with the second "PPU" BUS
func (cartridge *Cartridge) PpuRead(address uint16) bool {
	return false
}

func (cartridge *Cartridge) PpuWrite(address uint16, data byte) bool {
	return false
}

//The game "la cartouche"
type Cartridge struct {
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
	cartridge.mapperType = ((sHeader.mapper2 >> 4) << 4) | (sHeader.mapper1 >> 4)

	//get mirror mode
	cartridge.mirror = (sHeader.mapper1 & 1) | (((sHeader.mapper1 >> 3) & 1) << 1)

	// battery-backed RAM
	cartridge.battery = (sHeader.mapper1 >> 1) & 1

	// read trainer if present (unused)
	if sHeader.mapper1&0x04 == 4 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			println("call usage and exit")
		}
	}

	// read prg-rom bank(s)

	cartridge.prg = make([]byte, int(sHeader.prgRomChunks)*16384) //number mentioned // http://wiki.nesdev.com/w/index.php/INES // http://nesdev.com/NESDoc.pdf (page 28)

	if _, err := io.ReadFull(file, cartridge.prg); err != nil {
		println("call usage and exit")
	}

	// read chr-rom bank(s)

	// provide chr-rom/ram if not in file
	if sHeader.chrRomChunks == 0 {
		sHeader.chrRomChunks = 1
	}
	//make funtion allow memory allocation just like malloc
	cartridge.chr = make([]byte, int(sHeader.chrRomChunks)*8192) //number mentioned // http://wiki.nesdev.com/w/index.php/INES // http://nesdev.com/NESDoc.pdf (page 28)
	if _, err := io.ReadFull(file, cartridge.chr); err != nil {
		println("call usage and exit")
	}

	//sram allocation
	cartridge.sram = make([]byte, 0x2000)

	return &cartridge
}
