package ui

type View interface {
	Start()
	Update(dt float64)
	End()
}
