package nescomponents

import "fmt"

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte) bool
	Step()
}

func NewMapper(cartridge *Cartridge) (*Mapper, error) {
	//load appropriate mapper
	switch cartridge.mapperType {
	// case 0:
	// 	return NewMapper2(cartridge), nil
	case 1:
		mapper := NewMapper1(cartridge)
		return &mapper, nil
		// case 2:
		// 	return NewMapper2(cartridge), nil
		// case 3:
		// 	return NewMapper3(cartridge), nil
		// case 4:
		// 	return NewMapper4(console, cartridge), nil
		// case 7:
		// 	return NewMapper7(cartridge), nil
		// case 225:
		// 	return NewMapper225(cartridge), nil
	}
	err := fmt.Errorf("unsupported mapper: %d", cartridge.mapperType)
	return nil, err
}
