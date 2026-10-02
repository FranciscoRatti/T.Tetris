package scenes

import (
	"TTetris/src/lib"
	"TTetris/src/obj"
	"math"

	"github.com/gdamore/tcell/v3"
)

var (
	selectedPane             int
	selectedButtonKeyBinding int
	selectedButtonVolumen    int
	selectedButtonGame       int
	isKeyBindFocus           bool
)

var isConfigBackgroundRunning bool

func OpenConfig() {
	var channel chan bool

	// Reiniciar valores
	selectedButtonKeyBinding = 0
	selectedButtonVolumen = 0
	selectedButtonGame = 0
	selectedPane = 1
	isKeyBindFocus = false

	if obj.Config.Game.ShowBackground && (lib.Width > 90 || lib.Height > 30) {
		channel = startBackgroundAnimation(&isConfigBackgroundRunning, drawConfig)
	} else {
		drawConfig()
	}

	for {
		event := <-lib.Screen.EventQ()

		switch event := event.(type) {
		case *tcell.EventKey:
			k := event.Key()

			if isKeyBindFocus { // Cambiar tecla
				isKeyBindFocus = false

				var key obj.Key
				var text string

				if k == tcell.KeyRune {
					text = event.Str()
				} else {
					text = event.Name()
				}

				key = obj.ParseKey(text)

				switch selectedButtonKeyBinding {
				case 0:
					obj.KeyRight = key
					obj.Config.KeyBinding.Right = text
				case 1:
					obj.KeyDown = key
					obj.Config.KeyBinding.Down = text
				case 2:
					obj.KeyLeft = key
					obj.Config.KeyBinding.Left = text
				case 3:
					obj.KeyRotate = key
					obj.Config.KeyBinding.Rotate = text
				case 4:
					obj.KeyHold = key
					obj.Config.KeyBinding.Hold = text
				case 5:
					obj.KeyFloor = key
					obj.Config.KeyBinding.Floor = text
				case 6:
					obj.KeyPause = key
					obj.Config.KeyBinding.Pause = text
				case 7:
					obj.KeyMute = key
					obj.Config.KeyBinding.Mute = text
				}
			} else { // Desplazamiento
				switch k {
				case tcell.KeyDown: // Abajo
					obj.KeyAudio.Play()
					switch selectedPane {
					case 1:
						if selectedButtonKeyBinding == 7 {
							selectedButtonKeyBinding = 0
						} else {
							selectedButtonKeyBinding++
						}
					case 2:
						if selectedButtonVolumen == 3 {
							selectedButtonVolumen = 0
						} else {
							selectedButtonVolumen++
						}
					case 3:
						if selectedButtonGame == 2 {
							selectedButtonGame = 0
						} else {
							selectedButtonGame++
						}
					}
				case tcell.KeyUp: // Arriba
					obj.KeyAudio.Play()
					switch selectedPane {
					case 1:
						if selectedButtonKeyBinding == 0 {
							selectedButtonKeyBinding = 7
						} else {
							selectedButtonKeyBinding--
						}
					case 2:
						if selectedButtonVolumen == 0 {
							selectedButtonVolumen = 3
						} else {
							selectedButtonVolumen--
						}
					case 3:
						if selectedButtonGame == 0 {
							selectedButtonGame = 2
						} else {
							selectedButtonGame--
						}
					}
				case tcell.KeyRight: // Derecha
					if selectedPane == 2 {
						obj.KeyAudio.Play()
						switch selectedButtonVolumen {
						case 0:
							obj.ChangeMusicVolume(0.15)
						case 2:
							obj.ChangeEffectsVolume(0.15)
						}
					}
				case tcell.KeyLeft: // Izquierda
					if selectedPane == 2 {
						obj.KeyAudio.Play()
						switch selectedButtonVolumen {
						case 0:
							obj.ChangeMusicVolume(-0.15)
						case 2:
							obj.ChangeEffectsVolume(-0.15)
						}
					}
				case tcell.KeyEnter: // Enter
					switch selectedPane {
					case 1:
						isKeyBindFocus = true
					case 2:
						switch selectedButtonVolumen {
						case 1:
							obj.EnterAudio.Play()
							obj.ChangeMuteMusic()
						case 3:
							obj.EnterAudio.Play()
							obj.ChangeMuteEffects()
						}
					case 3:
						switch selectedButtonGame {
						case 0:
							obj.Config.Game.Shadow = !obj.Config.Game.Shadow
						case 1:
							obj.Config.Game.Hold = !obj.Config.Game.Hold
						case 2:
							obj.Config.Game.ShowBackground = !obj.Config.Game.ShowBackground
							if obj.Config.Game.ShowBackground {
								generateBackground()
								channel = startBackgroundAnimation(&isConfigBackgroundRunning, drawConfig)
							} else {
								isConfigBackgroundRunning = false
								if channel != nil {
									<-channel
								}
							}
						}
					}
				case tcell.KeyEscape: // Salir
					isConfigBackgroundRunning = false
					obj.EnterAudio.Play()
					return
				case tcell.KeyTab: // Cambiar de panel
					obj.KeyAudio.Play()
					if selectedPane < 3 {
						selectedPane++
					} else {
						selectedPane = 1
					}
				default:

					// Cambiar de panel
					switch event.Str() {
					case "1":
						obj.KeyAudio.Play()
						selectedPane = 1
					case "2":
						obj.KeyAudio.Play()
						selectedPane = 2
					case "3":
						obj.KeyAudio.Play()
						selectedPane = 3
					}

					if obj.KeyMute.Equals(k, event.Str()) { // Mute
						obj.ChangeMuteAll()
					}
				}
			}
		case *tcell.EventResize:
			isConfigBackgroundRunning = false
			if obj.Config.Game.ShowBackground && channel != nil {
				<-channel
			}

			lib.Width, lib.Height = lib.Screen.Size()

			if obj.Config.Game.ShowBackground {
				generateBackground()
				if lib.Width > 90 || lib.Height > 30 {
					channel = startBackgroundAnimation(&isConfigBackgroundRunning, drawConfig)
				}
			}
		}

		drawConfig()
	}
}

// Draw ----------------------------------------------------------------------------------------------------------------

func drawConfig() {
	lib.Screen.Clear()

	// Background
	if obj.Config.Game.ShowBackground && (lib.Width > 90 || lib.Height > 30) {
		drawBackground()

		minY := lib.Height/2 - 15
		minX := lib.Width/2 - 45
		for i := range 30 {
			for j := range 90 {
				lib.Screen.Put(j+minX, i+minY, " ", lib.DefaultStyle)
			}
		}
	}

	// Coordenadas
	minX := int(math.Max(0, float64(lib.Width/2-45)))
	maxX := int(math.Min(float64(lib.Width-1), float64(lib.Width/2+44)))

	minY := int(math.Max(0, float64(lib.Height/2-15)))
	maxY := int(math.Min(float64(lib.Height-1), float64(lib.Height/2+14)))
	differenceY := float64(maxY - minY)

	spacingY9 := differenceY / 9
	spacingY18 := differenceY / 18

	spacingY16 := differenceY / 16
	spacingY32 := differenceY / 32

	spacingY8 := (differenceY + 2) / 8

	spacingX4 := float64(maxX-minX) / 4

	// KeyBinding
	i := 0
	x := int(math.Round(float64(minX)+spacingX4)) - 7
	for y := float64(minY+1) + spacingY18; y <= float64(maxY-1) && i < 8; y += spacingY9 {
		if isKeyBindFocus && i == selectedButtonKeyBinding {
			lib.Screen.PutStrStyled(x, int(math.Round(y)), keyBindingSprites[i], lib.DefaultStyle)
		} else {
			lib.Screen.PutStrStyled(x, int(math.Round(y)), keyBindingSprites[i]+*keyBindingNames[i], lib.DefaultStyle)
		}
		i++
	}

	// Volume
	i = 0
	x = lib.Width/2 + int(math.Round(spacingX4))
	for y := float64(minY+1) + spacingY32; y <= float64(lib.Height/2) && i < 7; y += spacingY16 {
		lib.Screen.PutStrStyled(x-len(volumeSprites[i])/2, int(math.Round(y)), volumeSprites[i], lib.DefaultStyle)
		i++
	}
	lib.Screen.PutStrStyled(
		x-len(volumeSprite)/2+int(math.Round(obj.Config.Volume.Music.Value/0.15)),
		int(math.Round(float64(minY+1)+spacingY32+spacingY16)),
		"<>", lib.DefaultStyle)
	if obj.Config.Volume.Music.Mute {
		lib.Screen.PutStrStyled(
			x-len(noMuteSprite)/2+1,
			int(math.Round(float64(minY+1)+spacingY32+spacingY16*2)),
			"x", lib.DefaultStyle)
	}
	lib.Screen.PutStrStyled(
		x-len(volumeSprite)/2+int(math.Round(obj.Config.Volume.Effects.Value/0.15)),
		int(math.Round(float64(minY+1)+spacingY32+spacingY16*5)),
		"<>", lib.DefaultStyle)
	if obj.Config.Volume.Effects.Mute {
		lib.Screen.PutStrStyled(
			x-len(noMuteSprite)/2+1,
			int(math.Round(float64(minY+1)+spacingY32+spacingY16*6)),
			"x", lib.DefaultStyle)
	}

	// Game
	i = 0
	for y := float64(lib.Height/2+1) + spacingY16; y <= float64(maxY-1) && i < 3; y += spacingY8 {
		lib.Screen.PutStrStyled(x-len(gameSprites[i])/2, int(math.Round(y)), gameSprites[i], lib.DefaultStyle)
		i++
	}
	if obj.Config.Game.Shadow {
		lib.Screen.PutStrStyled(
			x-len(gameSprites[0])/2+1,
			int(math.Round(float64(lib.Height/2+1)+spacingY16)),
			"x", lib.DefaultStyle)
	}
	if obj.Config.Game.Hold {
		lib.Screen.PutStrStyled(
			x-len(gameSprites[1])/2+1,
			int(math.Round(float64(lib.Height/2+1)+spacingY16+spacingY8)),
			"x", lib.DefaultStyle)
	}
	if obj.Config.Game.ShowBackground {
		lib.Screen.PutStrStyled(
			x-len(gameSprites[1])/2+1,
			int(math.Round(float64(lib.Height/2+1)+spacingY16+spacingY8*2)),
			"x", lib.DefaultStyle)
	}

	// Marco
	lib.Screen.Put(minX, minY, "┌", lib.DefaultStyle)
	lib.Screen.Put(lib.Width/2, minY, "┬", lib.DefaultStyle)
	lib.Screen.Put(maxX, minY, "┐", lib.DefaultStyle)
	lib.Screen.Put(lib.Width/2, lib.Height/2, "├", lib.DefaultStyle)
	lib.Screen.Put(maxX, lib.Height/2, "┤", lib.DefaultStyle)
	lib.Screen.Put(minX, maxY, "└", lib.DefaultStyle)
	lib.Screen.Put(lib.Width/2, maxY, "┴", lib.DefaultStyle)
	lib.Screen.Put(maxX, maxY, "┘", lib.DefaultStyle)

	drawHorizontalLine(minX+1, lib.Width/2-1, minY, "─")
	drawHorizontalLine(minX+1, lib.Width/2-1, maxY, "─")
	drawHorizontalLine(lib.Width/2+1, maxX-1, minY, "─")
	drawHorizontalLine(lib.Width/2+1, maxX-1, lib.Height/2, "─")
	drawHorizontalLine(lib.Width/2+1, maxX-1, maxY, "─")

	drawVerticalLine(minY+1, maxY-1, minX, "│")
	drawVerticalLine(minY+1, lib.Height/2-1, lib.Width/2, "│")
	drawVerticalLine(minY+1, lib.Height/2-1, maxX, "│")
	drawVerticalLine(lib.Height/2+1, maxY-1, lib.Width/2, "│")
	drawVerticalLine(lib.Height/2+1, maxY-1, maxX, "│")

	lib.Screen.PutStrStyled(minX+((lib.Width/2-minX)/2)-7, minY, "[ KEY BINDING ]", lib.DefaultStyle)
	lib.Screen.PutStrStyled(maxX-(maxX-lib.Width/2)/2-5, minY, "[ VOLUME ]", lib.DefaultStyle)
	lib.Screen.PutStrStyled(maxX-(maxX-lib.Width/2)/2-6, lib.Height/2, "[ GENERAL ]", lib.DefaultStyle)

	// Selección
	switch selectedPane {
	case 1: // KeyBinding
		drawHorizontalLine(minX+1, lib.Width/2-1, minY, "═")
		drawHorizontalLine(minX+1, lib.Width/2-1, maxY, "═")

		drawVerticalLine(minY+1, maxY-1, minX, "║")
		drawVerticalLine(minY+1, maxY-1, lib.Width/2, "║")

		lib.Screen.PutStrStyled(minX+((lib.Width/2-minX)/2)-7, minY, "[ KEY BINDING ]", lib.SelectedStyle)

		lib.Screen.Put(minX, minY, "╔", lib.DefaultStyle)
		lib.Screen.Put(lib.Width/2, minY, "╗", lib.DefaultStyle)
		lib.Screen.Put(minX, maxY, "╚", lib.DefaultStyle)
		lib.Screen.Put(lib.Width/2, maxY, "╝", lib.DefaultStyle)

		// Teclas
		y := minY + 1 + int(math.Round(spacingY18+spacingY9*float64(selectedButtonKeyBinding)))
		x = int(math.Round(float64(minX)+spacingX4) - 7)
		if isKeyBindFocus {
			lib.Screen.PutStrStyled(x+9, y, "[ _ ]", lib.SelectedStyle)
		} else {
			lib.Screen.PutStrStyled(x+9, y, "[ "+*keyBindingNames[selectedButtonKeyBinding]+" ]", lib.SelectedStyle)
		}
	case 2: // Volume
		drawHorizontalLine(lib.Width/2+1, maxX-1, minY, "═")
		drawHorizontalLine(lib.Width/2+1, maxX-1, lib.Height/2, "═")

		drawVerticalLine(minY+1, lib.Height/2-1, lib.Width/2, "║")
		drawVerticalLine(minY+1, lib.Height/2-1, maxX, "║")

		lib.Screen.PutStrStyled(maxX-(maxX-lib.Width/2)/2-5, minY, "[ VOLUME ]", lib.SelectedStyle)
		lib.Screen.PutStrStyled(maxX-(maxX-lib.Width/2)/2-6, lib.Height/2, "[ GENERAL ]", lib.DefaultStyle)

		lib.Screen.Put(lib.Width/2, minY, "╔", lib.DefaultStyle)
		lib.Screen.Put(maxX, minY, "╗", lib.DefaultStyle)
		lib.Screen.Put(lib.Width/2, lib.Height/2, "╚", lib.DefaultStyle)
		lib.Screen.Put(maxX, lib.Height/2, "╝", lib.DefaultStyle)

		x = lib.Width/2 + int(math.Round(spacingX4))
		switch selectedButtonVolumen {
		case 0: // Volumen de musica
			lib.Screen.PutStrStyled(
				x-len(volumeSprite)/2+int(math.Round(obj.Config.Volume.Music.Value/0.15)),
				int(math.Round(float64(minY+1)+spacingY32+spacingY16)),
				"<>", lib.SelectedStyle)
		case 1: // Mute de musica
			lib.Screen.PutStrStyled(
				x-len(noMuteSprite)/2,
				int(math.Round(float64(minY+1)+spacingY32+spacingY16*2)),
				"[ ]", lib.SelectedStyle)
			if obj.Config.Volume.Music.Mute {
				lib.Screen.PutStrStyled(
					x-len(noMuteSprite)/2+1,
					int(math.Round(float64(minY+1)+spacingY32+spacingY16*2)),
					"x", lib.SelectedStyle)
			}
		case 2: // Volumen de efectos
			lib.Screen.PutStrStyled(
				x-len(volumeSprite)/2+int(math.Round(obj.Config.Volume.Effects.Value/0.15)),
				int(math.Round(float64(minY+1)+spacingY32+spacingY16*5)),
				"<>", lib.SelectedStyle)
		case 3: // Mute de efectos
			lib.Screen.PutStrStyled(
				x-len(noMuteSprite)/2,
				int(math.Round(float64(minY+1)+spacingY32+spacingY16*6)),
				"[ ]", lib.SelectedStyle)
			if obj.Config.Volume.Effects.Mute {
				lib.Screen.PutStrStyled(
					x-len(noMuteSprite)/2+1,
					int(math.Round(float64(minY+1)+spacingY32+spacingY16*6)),
					"x", lib.SelectedStyle)
			}
		}
	case 3: // Game
		drawHorizontalLine(lib.Width/2+1, maxX-1, lib.Height/2, "═")
		drawHorizontalLine(lib.Width/2+1, maxX-1, maxY, "═")

		drawVerticalLine(lib.Height/2+1, maxY-1, lib.Width/2, "║")
		drawVerticalLine(lib.Height/2+1, maxY-1, maxX, "║")

		lib.Screen.PutStrStyled(maxX-(maxX-lib.Width/2)/2-6, lib.Height/2, "[ GENERAL ]", lib.SelectedStyle)

		lib.Screen.Put(lib.Width/2, lib.Height/2, "╔", lib.DefaultStyle)
		lib.Screen.Put(maxX, lib.Height/2, "╗", lib.DefaultStyle)
		lib.Screen.Put(lib.Width/2, maxY, "╚", lib.DefaultStyle)
		lib.Screen.Put(maxX, maxY, "╝", lib.DefaultStyle)

		x := lib.Width/2 + int(math.Round(spacingX4)) - len(gameSprites[0])/2
		switch selectedButtonGame {
		case 0: // Mostrar sombras
			lib.Screen.PutStrStyled(
				x,
				int(math.Round(float64(lib.Height/2+1)+spacingY16)),
				"[ ]", lib.SelectedStyle)
			if obj.Config.Game.Shadow {
				lib.Screen.PutStrStyled(
					x+1,
					int(math.Round(float64(lib.Height/2+1)+spacingY16)),
					"x", lib.SelectedStyle)
			}
		case 1: // Activar hold
			lib.Screen.PutStrStyled(
				x,
				int(math.Round(float64(lib.Height/2+1)+spacingY16+spacingY8)),
				"[ ]", lib.SelectedStyle)
			if obj.Config.Game.Hold {
				lib.Screen.PutStrStyled(
					x+1,
					int(math.Round(float64(lib.Height/2+1)+spacingY16+spacingY8)),
					"x", lib.SelectedStyle)
			}
		case 2: // Mostrar fondo
			lib.Screen.PutStrStyled(
				x,
				int(math.Round(float64(lib.Height/2+1)+spacingY16+spacingY8*2)),
				"[ ]", lib.SelectedStyle)
			if obj.Config.Game.ShowBackground {
				lib.Screen.PutStrStyled(
					x+1,
					int(math.Round(float64(lib.Height/2+1)+spacingY16+spacingY8*2)),
					"x", lib.SelectedStyle)
			}
		}
	}

	lib.Screen.Show()
}

func drawHorizontalLine(min, max, y int, char string) {
	for ; min <= max; min++ {
		lib.Screen.Put(min, y, char, lib.DefaultStyle)
	}
}

func drawVerticalLine(min, max, x int, char string) {
	for ; min <= max; min++ {
		lib.Screen.Put(x, min, char, lib.DefaultStyle)
	}
}

// Sprites -------------------------------------------------------------------------------------------------------------

var (
	keyBindingSprites = []string{
		"RIGHT  :   ",
		"DOWN   :   ",
		"LEFT   :   ",
		"ROTATE :   ",
		"HOLD   :   ",
		"FLOOR  :   ",
		"PAUSE  :   ",
		"MUTE   :   ",
	}
	keyBindingNames = []*string{
		&obj.KeyRight.Name,
		&obj.KeyDown.Name,
		&obj.KeyLeft.Name,
		&obj.KeyRotate.Name,
		&obj.KeyHold.Name,
		&obj.KeyFloor.Name,
		&obj.KeyPause.Name,
		&obj.KeyMute.Name,
	}

	volumeSprite  = "{--------------------------}"
	noMuteSprite  = "[ ] MUTE"
	volumeSprites = []string{
		"- MUSIC -",
		volumeSprite,
		noMuteSprite,
		"",
		"- EFFECTS -",
		volumeSprite,
		noMuteSprite,
	}

	gameSprites = []string{
		"[ ] SHOW SHADOW    ",
		"[ ] ENABLE HOLD    ",
		"[ ] SHOW BACKGROUND",
	}
)
