package scenes

import (
	"TTetris/src/lib"
	"TTetris/src/obj"

	"github.com/gdamore/tcell/v3"
)

var selectedButton = 0

func OpenMenu() {
	drawMenu()

	obj.MenuAudio.Play()

	// Bucle
	for {
		event := <-lib.Screen.EventQ()

		switch event := event.(type) {
		case *tcell.EventKey:
			k := event.Key()

			// Navegación
			switch k {
			case tcell.KeyDown:
				obj.KeyAudio.Play()

				if selectedButton == 2 {
					selectedButton = 0
				} else {
					selectedButton++
				}
			case tcell.KeyUp:
				obj.KeyAudio.Play()

				if selectedButton == 0 {
					selectedButton = 2
				} else {
					selectedButton--
				}
			case tcell.KeyEnter:
				obj.EnterAudio.Play()

				switch selectedButton {
				case 0:
					obj.MenuAudio.Stop()
					StartNewGame()
					obj.MenuAudio.Play()
				case 1:
					OpenConfig()
					obj.WriteConfig()
				case 2:
					return
				}

				selectedButton = 0
			default:
				if obj.KeyMute.Equals(k, event.Str()) {
					if lib.Mute {
						lib.Mute = false
						obj.ChangeMuteEffectsWithoutChange(obj.Config.Volume.Effects.Mute)
						obj.ChangeMuteMusicWithoutChange(obj.Config.Volume.Music.Mute)
					} else {
						lib.Mute = true
						obj.ChangeMuteEffectsWithoutChange(true)
						obj.ChangeMuteMusicWithoutChange(true)
					}
				}
			}
		case *tcell.EventResize:
			lib.Width, lib.Height = lib.Screen.Size()
		}

		// Pintar
		drawMenu()
	}
}

func drawMenu() {
	lib.Screen.Clear()

	switch selectedButton {
	case 0:
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+2, "[ "+startButton+" ]", lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+5, "  "+optionsButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+8, "  "+exitButton, lib.DefaultStyle)
	case 1:
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+2, "  "+startButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+5, "[ "+optionsButton+" ]", lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+8, "  "+exitButton, lib.DefaultStyle)
	case 2:
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+2, "  "+startButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+5, "  "+optionsButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+8, "[ "+exitButton+" ]", lib.DefaultStyle)
	default:
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+2, "  "+startButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+5, "  "+optionsButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+8, "  "+exitButton, lib.DefaultStyle)
	}

	lib.DrawString(lib.Width/2-28, lib.Height/2-9, logo, lib.DefaultStyle, lib.Screen)

	lib.Screen.Show()
}

// Sprites
var (
	logo = []string{
		"╔═╦═════════════════════════════════════════════════╦═╗",
		"╠═╝                                                 ╚═╣",
		"║   [][][]    [][][] [][][] [][][] [][][] [] [][][]   ║",
		"║     []        []   []       []   []  []    []       ║",
		"║     []        []   [][]     []   [][]   [] [][][]   ║",
		"║     []        []   []       []   []  [] []     []   ║",
		"║     []   []   []   [][][]   []   []  [] [] [][][]   ║",
		"╠═╗                                                 ╔═╣",
		"╚═╩═════════════════════════════════════════════════╩═╝",
	}

	// Buttons
	startButton   = " PLAY "
	optionsButton = "CONFIG"
	exitButton    = " EXIT "
)
