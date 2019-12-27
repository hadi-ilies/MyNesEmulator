package ui

type View interface {
	Start()
	Update()
	End()
}
