package nescomponents

//Controller is a nes controller
type Controller struct {
	buttons [8]byte
	strobe  byte
	index   byte
}

//NewController Controller constructor
func NewController() *Controller {
	return &Controller{}
}

//GetButton can be useful
func (controller *Controller) GetButton() [8]byte {
	return controller.buttons
}

//SetButtons allow me to set buttons throught the ui
func (controller *Controller) SetButtons(buttons [8]byte) {
	controller.buttons = buttons
}

func (controller *Controller) Read() byte {
	var value byte = 0

	if controller.index < 8 && controller.buttons[controller.index] == 0 {
		value = 1
	}
	controller.index++
	if controller.strobe == 1 {
		controller.index = 0
	}
	return value
}

func (controller *Controller) Write(value byte) {
	controller.strobe = value
	if controller.strobe == 1 {
		controller.index = 0
	}
}
