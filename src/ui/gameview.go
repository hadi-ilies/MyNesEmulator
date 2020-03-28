package ui

import (
	"image"
	//	"os"

	oglEncap "./openglencapsulation" // import and rename the package openglencapsulation to oglEncap
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/hadi-ilies/MyNesEmulator/src/nes"
)

//GameView struct that reprensent the gameview
type GameView struct {
	nes     *nes.Nes
	ui      *Ui // lol there is no inerittance in golang, I am a noob ':(
	texture uint32
	frames  []image.Image
}

//NewGameView gameview constructor
func NewGameView(ui *Ui, nes *nes.Nes) View {
	var gameView GameView

	gameView.texture = oglEncap.CreateTexture()
	gameView.nes = nes
	gameView.ui = ui
	return &gameView
}

func (gameView *GameView) Start() {
	gl.ClearColor(0, 0, 0, 1)
	// view.director.SetTitle(view.title)
	// view.console.SetAudioChannel(view.director.audio.channel)
	// view.console.SetAudioSampleRate(view.director.audio.sampleRate)
	gameView.ui.GetWindow().SetKeyCallback(gameView.onKey) // todo getWindow can be removed
	// load state
	// if err := view.console.LoadState(savePath(view.hash)); err == nil {
	// 	return
	// } else {
	gameView.nes.Reset() //init nes
	// //}
	// // load sram
	//cartridge := gameView.nes.GetComponents().GetCartridge()
	// if cartridge.Battery != 0 {
	// 	if sram, err := readSRAM(sramPath(gameView.hash)); err == nil {
	// 		cartridge.SRAM = sram
	// 	}
	// }
}

func (gameView *GameView) Update(dt float64) {
	if dt > 1 {
		dt = 0
	}
	// if joystickReset(glfw.Joystick1) {
	// 	view.director.ShowMenu()
	// }
	// if joystickReset(glfw.Joystick2) {
	// 	view.director.ShowMenu()
	// }
	// if readKey(window, glfw.KeyEscape) {
	// 	view.director.ShowMenu()
	// }
	updateControllers(gameView.ui.GetWindow(), gameView.nes) // todo code this func
	gameView.nes.Run(dt)
	gl.BindTexture(gl.TEXTURE_2D, gameView.texture)
	oglEncap.SetTexture(gameView.nes.PixelBuffer()) //todo code the buffer
	// println("VIEW")
	// os.Exit(0)
	gameView.drawBuffer(gameView.ui.GetWindow().GetFramebufferSize())
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (gameView *GameView) End() {
	gameView.ui.GetWindow().SetKeyCallback(nil)
	// view.console.SetAudioChannel(nil)
	// view.console.SetAudioSampleRate(0)
	// // save sram
	// cartridge := view.console.Cartridge
	// if cartridge.Battery != 0 {
	// 	writeSRAM(sramPath(view.hash), cartridge.SRAM)
	// }
	// // save state
	// view.console.SaveState(savePath(view.hash))
}

//will be useful when i will emulate controllers and physics interactions with my nes
/** PRIVATE METHODS **/
func (view *GameView) onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press { // if key is pressed
		switch key { //check which key has been pressed
		//case glfw.KeySpace:
		//	screenshot(view.console.Buffer())
		case glfw.KeyR: // if i pressed r key i will restart my nes
			view.nes.Reset()
		}
	}
}

func (view *GameView) drawBuffer(bufferWidth int, bufferHeight int) {
	padding := 0
	s1 := float32(bufferWidth) / 256
	s2 := float32(bufferHeight) / 240
	f := float32(1 - padding)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func readKey(window *glfw.Window, key glfw.Key) byte {
	if (window.GetKey(key) == glfw.Press) == true {
		return 0
	}
	return 1
}

func readKeys(window *glfw.Window) [8]byte {
	var result [8]byte

	result[nes.KeyA] = readKey(window, glfw.KeyA)
	result[nes.KeyB] = readKey(window, glfw.KeyS)
	result[nes.KeySelect] = readKey(window, glfw.KeyLeftShift)
	result[nes.KeyStart] = readKey(window, glfw.KeyEnter)
	result[nes.KeyUp] = readKey(window, glfw.KeyUp)
	result[nes.KeyDown] = readKey(window, glfw.KeyDown)
	result[nes.KeyLeft] = readKey(window, glfw.KeyLeft)
	result[nes.KeyRight] = readKey(window, glfw.KeyRight)
	// println()
	// for _, value := range result {
	// 	print(" ", value, " ")
	// }
	// println()
	return result
}

func updateControllers(window *glfw.Window, nes *nes.Nes) {
	nes.SetButtonToController(readKeys(window))
}
