package scenes

import (
	"TTetris/src/lib"
	"TTetris/src/obj"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v3"
)

var selectedButton = 0
var isMenuBackgroundRunning = false

func OpenMenu() {
	var channel chan bool
	if obj.Config.Game.ShowBackground {
		generateBackground()
		channel = startBackgroundAnimation(&isMenuBackgroundRunning, drawMenu)
	}
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
					isMenuBackgroundRunning = false
					if obj.Config.Game.ShowBackground && channel != nil {
						<-channel
					}

					obj.MenuAudio.Stop()
					StartNewGame()
					obj.MenuAudio.Play()
					selectedButton = 0

					if obj.Config.Game.ShowBackground {
						channel = startBackgroundAnimation(&isMenuBackgroundRunning, drawMenu)
					}
				case 1:
					isMenuBackgroundRunning = false
					if obj.Config.Game.ShowBackground && channel != nil {
						<-channel
					}

					OpenConfig()
					selectedButton = 0

					if obj.Config.Game.ShowBackground {
						channel = startBackgroundAnimation(&isMenuBackgroundRunning, drawMenu)
					}
				case 2:
					return
				}

			// Mute
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
			isMenuBackgroundRunning = false
			if obj.Config.Game.ShowBackground && channel != nil {
				<-channel
			}

			lib.Width, lib.Height = lib.Screen.Size()

			if obj.Config.Game.ShowBackground {
				generateBackground()
				channel = startBackgroundAnimation(&isMenuBackgroundRunning, drawMenu)
			}
		}

		// Pintar
		drawMenu()
	}
}

func drawMenu() {
	lib.Screen.Clear()

	if obj.Config.Game.ShowBackground {
		drawBackground()
	}

	switch selectedButton {
	case 0:
		lib.Screen.PutStrStyled(lib.Width/2-4, lib.Height/2+2, "[ "+startButton+" ]", lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-3, lib.Height/2+5, optionsButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-2, lib.Height/2+8, exitButton, lib.DefaultStyle)
	case 1:
		lib.Screen.PutStrStyled(lib.Width/2-2, lib.Height/2+2, startButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-5, lib.Height/2+5, "[ "+optionsButton+" ]", lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-2, lib.Height/2+8, exitButton, lib.DefaultStyle)
	case 2:
		lib.Screen.PutStrStyled(lib.Width/2-2, lib.Height/2+2, startButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-3, lib.Height/2+5, optionsButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-4, lib.Height/2+8, "[ "+exitButton+" ]", lib.DefaultStyle)
	default:
		lib.Screen.PutStrStyled(lib.Width/2-2, lib.Height/2+2, startButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-3, lib.Height/2+5, optionsButton, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2-2, lib.Height/2+8, exitButton, lib.DefaultStyle)
	}

	lib.DrawString(lib.Width/2-28, lib.Height/2-9, logo, lib.DefaultStyle)

	lib.Screen.Show()
}

func drawBackground() {
	for iR, r := range background {
		for iC, c := range r {
			if backgroundShadow[iR][iC] != ' ' {
				lib.Screen.Put(iC, iR, "|", lib.BackgroundShadowStyles[backgroundShadow[iR][iC]-1])
			}

			if c == '[' || c == ']' {
				lib.Screen.Put(iC, iR, string(c), lib.BackgroundPieceStyle)
			}
		}
	}
}

func generateBackground() {
	background = make([][]byte, lib.Height)
	backgroundShadow = make([][]rune, lib.Height)
	for i := range background {
		background[i] = make([]byte, lib.Width)
		backgroundShadow[i] = make([]rune, lib.Width)
		for j := range backgroundShadow[i] {
			backgroundShadow[i][j] = ' '
		}
	}

	cantPieces := (lib.Height * lib.Width) / 200

	for range cantPieces {
		var id, x, y int
		var sprite []string
		var spriteHeight, spriteWidth int

		isRightOnTop := true
		for isRightOnTop {
			id = rand.Intn(7)
			sprite = obj.Pieces[id].GetSprite(rand.Intn(4))

			spriteHeight = len(sprite)
			spriteWidth = 0
			for _, row := range sprite {
				if len(row) > spriteWidth {
					spriteWidth = len(row)
				}
			}

			x, y = rand.Intn(lib.Width-spriteWidth-2), rand.Intn(lib.Height-spriteHeight)
			x++

			minY := y - 10
			maxY := y + spriteHeight

			minX := x - 1
			maxX := x + spriteWidth + 1

			isRightOnTop = false
			for i := minY; i < maxY; i++ {
				for j := minX; j < maxX; j++ {
					if i < 0 {
						if background[lib.Height+i][j] != 0 || backgroundShadow[lib.Height+i][j] != ' ' {
							isRightOnTop = true
							break
						}
					} else {
						if background[i][j] != 0 || backgroundShadow[i][j] != ' ' {
							isRightOnTop = true
							break
						}
					}
				}

				if isRightOnTop {
					break
				}
			}
		}

		for iR := spriteHeight - 1; iR >= 0; iR-- {
			for iC, c := range sprite[iR] {
				if c == '[' || c == ']' {
					for k := 10; k > 0; k-- {
						shadowY := y + iR - k
						if shadowY < 0 {
							backgroundShadow[lib.Height+shadowY][x+iC] = rune(k)
						} else {
							backgroundShadow[shadowY][x+iC] = rune(k)
						}
					}
				}
			}
		}

		for iR, r := range sprite {
			for iC, c := range r {
				if c != ' ' {
					background[y+iR][x+iC] = byte(c)
				}
			}
		}
	}
}

func startBackgroundAnimation(condition *bool, draw func()) chan bool {
	channel := make(chan bool)

	*condition = true
	loop := func() {
		lastExec := time.Now()
		for *condition {
			if time.Since(lastExec) >= time.Duration(50)*time.Millisecond {
				lastLineBackground := background[len(background)-1]
				lastLineShadow := backgroundShadow[len(backgroundShadow)-1]

				for y := len(background) - 2; y >= 0; y-- {
					background[y+1] = background[y]
					backgroundShadow[y+1] = backgroundShadow[y]
				}
				background[0] = lastLineBackground
				backgroundShadow[0] = lastLineShadow

				drawBackground()
				draw()
				lastExec = time.Now()
			}
		}

		channel <- true
	}

	go loop()
	return channel
}

// Sprites
var (
	background       [][]byte
	backgroundShadow [][]rune

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
	startButton   = "PLAY"
	optionsButton = "CONFIG"
	exitButton    = "EXIT"
)
