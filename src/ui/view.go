package ui

//View interface will be useful when i will create a menu and add features to my nes emulator
type View interface {
	Start()
	Update(dt float64)
	End()
}
