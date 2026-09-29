package scenes

import (
	"TTetris/src/lib"
	"TTetris/src/obj"
	"math/rand"
	"slices"
	"strconv"

	"github.com/gdamore/tcell/v3"
)

var (
	isRunning  bool
	isHold     bool
	isGameOver bool
	isPause    bool

	currentPiece *obj.Piece
	nextPieces   [3]*obj.Piece
	holdPiece    *obj.Piece

	score int64
	lines uint64
	level = uint8(1)
)

func StartNewGame() {
	initializeGame()
	drawGame()
	obj.Timer.Start()

	// Bucle ------------------------------------------------
	for isRunning {

		// Eventos
		event := <-lib.Screen.EventQ()

		switch event := event.(type) {
		case *tcell.EventKey:
			k := event.Key()
			s := event.Str()

			if isGameOver {
				switch k {
				case tcell.KeyDown:
					obj.KeyAudio.Play()

					if selectedButton == 1 {
						selectedButton = 0
					} else {
						selectedButton = 1
					}
				case tcell.KeyUp:
					obj.KeyAudio.Play()

					if selectedButton == 0 {
						selectedButton = 1
					} else {
						selectedButton = 0
					}
				case tcell.KeyEnter:
					obj.EnterAudio.Play()

					switch selectedButton {
					case 0:
						initializeGame()
						obj.Timer.Start()
					case 1:
						isRunning = false
						obj.Timer.Stop()
						obj.GameAudio.Stop()
					}
				default:

				}
			} else if isPause {
				if k == tcell.KeyDown {
					obj.KeyAudio.Play()

					if selectedButton == 2 {
						selectedButton = 0
					} else {
						selectedButton++
					}
				} else if k == tcell.KeyUp {
					obj.KeyAudio.Play()

					if selectedButton == 0 {
						selectedButton = 2
					} else {
						selectedButton--
					}
				} else if k == tcell.KeyEnter {
					switch selectedButton {
					case 0:
						obj.ResumeAudio.Play()

						isPause = false
						obj.Timer.Resume()
					case 1:
						obj.EnterAudio.Play()

						initializeGame()
						obj.Timer.Start()
					case 2:
						obj.EnterAudio.Play()

						isRunning = false
						obj.Timer.Stop()
						obj.GameAudio.Stop()
					}
				} else if obj.KeyPause.Equals(k, s) {
					obj.ResumeAudio.Play()

					isPause = false
					obj.Timer.Resume()
				}
			} else {
				if obj.KeyRight.Equals(k, s) {
					obj.MoveAudio.Play()
					currentPiece.MoveRight(stationaryPieces)
				} else if obj.KeyDown.Equals(k, s) {
					obj.MoveAudio.Play()
					obj.Timer.UpdateLastExec()
					if !currentPiece.MoveDown(stationaryPieces) {
						onChangePiece()
						break
					} else {
						score++
					}
				} else if obj.KeyLeft.Equals(k, s) {
					obj.MoveAudio.Play()
					currentPiece.MoveLeft(stationaryPieces)
				} else if obj.KeyRotate.Equals(k, s) {
					obj.RotateAudio.Play()
					currentPiece.RotateRight(stationaryPieces)
				} else if obj.KeyHold.Equals(k, s) {
					if obj.Config.Game.Hold && !isHold {
						obj.KeyAudio.Play()

						isHold = true

						if holdPiece == nil {
							holdPiece = currentPiece
							currentPiece = nextPieces[0]
							currentPiece.SetPos(8, 0)
							currentPiece.SetRotation(0)
							nextPieces[0] = nextPieces[1]
							nextPieces[1] = nextPieces[2]
							nextPieces[2] = obj.NewPiece(nextPieces)
						} else {
							lastPiece := currentPiece
							currentPiece = holdPiece
							currentPiece.SetPos(8, 0)
							currentPiece.SetRotation(0)
							holdPiece = lastPiece
						}
					}
				} else if obj.KeyFloor.Equals(k, s) {
					obj.FloorAudio.Play()

					obj.Timer.UpdateLastExec()
					currentY := currentPiece.GetY()
					currentPiece.MoveFloor(stationaryPieces)
					score += int64(currentPiece.GetY() - currentY)

					onChangePiece()
				} else if obj.KeyPause.Equals(k, s) {
					obj.PauseAudio.Play()

					isPause = true
					selectedButton = 0
					obj.Timer.Stop()
				}
			}

			if obj.KeyMute.Equals(k, s) {
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
		case *tcell.EventResize:
			lib.Width, lib.Height = lib.Screen.Size()
		}

		// Pintar
		drawGame()
	}

	// Fin
	obj.Timer.Stop()
}

// Funciones -----------------------------------------

func drawGame() {
	lib.Screen.Clear()

	// Marco
	lib.DrawString(lib.Width/2-12, lib.Height/2-11, frame, lib.DefaultStyle)
	if !isPause {
		lib.DrawChars(lib.Width/2-10, lib.Height/2-11, stationaryPieces, stationaryColors)

		// Pieza
		currentPiece.Draw(stationaryPieces)

		// Siguiente
		lib.DrawString(lib.Width/2+14, lib.Height/2-11, []string{"NEXT"}, lib.DefaultStyle)
		var height int
		for _, p := range nextPieces {
			lib.DrawWithFrame(lib.Width/2+13, lib.Height/2-10+height, p.GetSprite(0), p.GetStyle())
			height += len(p.GetSprite(0)) + 2
		}

		// Hold
		if holdPiece != nil {
			lib.DrawString(lib.Width/2-18, lib.Height/2-11, []string{"HOLD"}, lib.DefaultStyle)
			var width int
			var length int
			for _, r := range holdPiece.GetSprite(0) {
				length = len(r)
				if length > width {
					width = length
				}
			}
			lib.DrawWithFrame((lib.Width/2-15)-width, lib.Height/2-10, holdPiece.GetSprite(0), holdPiece.GetStyle())
		}
	} else {

		// Pausa
		lib.DrawString(lib.Width/2-7, lib.Height/2-2, pauseMessage, lib.DefaultStyle)

		alto := 0
		switch selectedButton {
		case 0:
			alto = 0
		case 1:
			alto = 1
		case 2:
			alto = 2
		}
		lib.DrawString(lib.Width/2-5, lib.Height/2-1+alto, []string{">"}, lib.DefaultStyle)
	}

	// Nivel
	lib.DrawString(lib.Width/2-21, lib.Height/2-4, []string{
		" LEVEL",
		" _---_ ",
		"{ " + lib.AppendCero(strconv.Itoa(int(level)), 3) + " }",
		" ¯---¯ ",
	}, lib.DefaultStyle)

	// Lineas
	lib.DrawString(lib.Width/2-25, lib.Height/2+1, []string{
		"     LINES",
		" _-------_ ",
		"{ " + lib.AppendCero(strconv.Itoa(int((lines%1000000)/1000)), 3) + "." + lib.AppendCero(strconv.Itoa(int(lines%1000)), 3) + " }",
		" ¯-------¯ ",
	}, lib.DefaultStyle)

	// Puntuacion
	lib.DrawString(lib.Width/2-29, lib.Height/2+6, []string{
		"         SCORE",
		" _-----------_ ",
		"{ " + lib.AppendCero(strconv.Itoa(int(score%1000000000)/1000000), 3) + "." + lib.AppendCero(strconv.Itoa(int((score%1000000)/1000)), 3) + "." + lib.AppendCero(strconv.Itoa(int(score%1000)), 3) + " }",
		" ¯-----------¯ ",
	}, lib.DefaultStyle)

	// Game over
	if isGameOver {
		lib.DrawString(lib.Width/2-7, lib.Height/2-2, gameOverMessage, lib.GameOverStyle)

		alto := 0
		switch selectedButton {
		case 0:
			alto = 0
		case 1:
			alto = 1
		}
		lib.DrawString(lib.Width/2-5, lib.Height/2-1+alto, []string{">"}, lib.GameOverStyle)
	}

	lib.Screen.Show()
}

func initializeGame() {

	// Variables
	isRunning = true
	isGameOver = false
	isHold = false
	isPause = false

	currentPiece = &obj.Pieces[rand.Intn(7)]
	currentPiece.SetPos(8, 0)
	currentPiece.SetRotation(0)

	nextPieces = [3]*obj.Piece([]*obj.Piece{&obj.Pieces[rand.Intn(7)], &obj.Pieces[rand.Intn(7)], &obj.Pieces[rand.Intn(7)]})
	holdPiece = nil

	score = 0
	lines = 0
	level = 1

	stationaryPieces = make([]string, 20)
	for i := range stationaryPieces {
		stationaryPieces[i] = "                    "
	}

	stationaryColors = make([][]tcell.Style, 20)
	for i := range stationaryColors {
		stationaryColors[i] = make([]tcell.Style, 20)
		for j := range stationaryColors[i] {
			stationaryColors[i][j] = lib.DefaultStyle
		}
	}

	// Timer
	obj.Timer.Initialize(800, func() {

		// Cambiar pieza
		if !currentPiece.MoveDown(stationaryPieces) {
			onChangePiece()
		}

		drawGame()
	})

	// Audio
	obj.GameAudio.Play()
}

func onChangePiece() {
	isHold = false
	putStationaryPieces(*currentPiece)
	checkLines()
	checkGameOver()

	currentPiece = nextPieces[0]
	currentPiece.SetPos(8, 0)
	currentPiece.SetRotation(0)
	nextPieces[0] = nextPieces[1]
	nextPieces[1] = nextPieces[2]
	nextPieces[2] = obj.NewPiece(nextPieces)
}

func putStationaryPieces(piece obj.Piece) {
	x := piece.GetX()
	y := piece.GetY()

	for iR, r := range piece.GetActualSprite() {
		line := stationaryPieces[y+iR]
		row := line[:x]

		for iC, c := range r {
			if c == ' ' {
				row += string(line[x+iC])
			} else {
				row += string(c)
				stationaryColors[y+iR][x+iC] = piece.GetStyle()
			}
		}

		row += line[x+len(r):]

		stationaryPieces[y+iR] = row
	}
}

func checkLines() {
	var continuosLines int
	isLevelUp := false

	for r := 19; r >= 0; r-- {
		cells := 0
		for c := 0; c < 20; c += 2 {
			if stationaryPieces[r][c] == '[' {
				cells++
			} else {
				break
			}
		}

		if cells == 10 {
			lines++
			continuosLines++

			if lines%10 == 0 {
				level++
				isLevelUp = true
				if level == 30 {
					obj.Timer.ChangeSleep(float32(1 * 1000 / 60))
				} else if level >= 9 && level < 30 {
					obj.Timer.ChangeSleep(obj.Timer.GetSleep() - float32(1*1000/60))
				} else if level < 9 {
					obj.Timer.ChangeSleep(obj.Timer.GetSleep() - float32(5*1000/60))
				}
			}

			for r2 := r; r2 >= 0; r2-- {
				if r2 == 0 {
					stationaryPieces[r2] = "                    "
					for c2 := range 20 {
						stationaryColors[r2][c2] = lib.DefaultStyle
					}
				} else {
					stationaryPieces[r2] = stationaryPieces[r2-1]
					stationaryColors[r2] = slices.Clone(stationaryColors[r2-1])
				}
			}

			r++
		}
	}

	if continuosLines != 0 {
		if isLevelUp {
			obj.LevelUpAudio.Play()
		} else if continuosLines == 4 {
			obj.Line4Audio.Play()
		} else {
			obj.LineAudio.Play()
		}

		switch continuosLines {
		case 1:
			score += int64(40 * level)
		case 2:
			score += int64(100 * level)
		case 3:
			score += int64(300 * int(level))
		case 4:
			score += int64(1200 * int(level))
		}
	}
}

func checkGameOver() {
	if stationaryPieces[0] != "                    " {
		isGameOver = true
		obj.Timer.Stop()
		obj.GameAudio.Stop()
		obj.GameOverAudio.Play()
	}
}

var (
	frame = []string{
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<! . . . . . . . . . .!>",
		"<!====================!>",
		" \\/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\/",
	}

	stationaryPieces []string
	stationaryColors [][]tcell.Style

	gameOverMessage = []string{
		"┌=┤GAMEOVER├=┐",
		"║   RESTART  ║",
		"║   EXIT     ║",
		"└============┘",
	}

	pauseMessage = []string{
		"┌===┤MENU├===┐",
		"║   RESUME   ║",
		"║   RESTART  ║",
		"║   EXIT     ║",
		"└============┘",
	}
)
