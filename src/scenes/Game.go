package scenes

import (
	"TTetris/src/lib"
	"TTetris/src/obj"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var (
	isRunning         bool
	isHold            bool
	isGameOver        bool
	isPause           bool
	isFireworkShowing bool

	currentPiece *obj.Piece
	nextPieces   [3]*obj.Piece
	holdPiece    *obj.Piece

	score            int64
	lines            uint64
	level            uint8
	newEntryPosition byte
	newEntryName     string
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
				if newEntryPosition != 0 {
					if k == tcell.KeyEsc {
						newEntryPosition = 0
					} else if k == tcell.KeyEnter {
						if len(newEntryName) == 3 {
							needReorder := 0

							if obj.Scoreboard.First.Name == newEntryName {
								if obj.Scoreboard.First.Score < score {
									obj.Scoreboard.First.Score = score
									needReorder = 1
								}
							} else if obj.Scoreboard.Second.Name == newEntryName {
								if obj.Scoreboard.Second.Score < score {
									obj.Scoreboard.Second.Score = score
									needReorder = 2
								}
							} else if obj.Scoreboard.Third.Name == newEntryName {
								if obj.Scoreboard.Third.Score < score {
									obj.Scoreboard.Third.Score = score
									needReorder = 3
								}
							} else if obj.Scoreboard.Fourth.Name == newEntryName {
								if obj.Scoreboard.Fourth.Score < score {
									obj.Scoreboard.Fourth.Score = score
									needReorder = 4
								}
							} else if obj.Scoreboard.Fifth.Name == newEntryName {
								if obj.Scoreboard.Fifth.Score < score {
									obj.Scoreboard.Fifth.Score = score
									needReorder = 5
								}
							} else {
								if newEntryPosition < 5 {
									obj.Scoreboard.Fifth.Name = obj.Scoreboard.Fourth.Name
									obj.Scoreboard.Fifth.Score = obj.Scoreboard.Fourth.Score
								}
								if newEntryPosition < 4 {
									obj.Scoreboard.Fourth.Name = obj.Scoreboard.Third.Name
									obj.Scoreboard.Fourth.Score = obj.Scoreboard.Third.Score
								}
								if newEntryPosition < 3 {
									obj.Scoreboard.Third.Name = obj.Scoreboard.Second.Name
									obj.Scoreboard.Third.Score = obj.Scoreboard.Second.Score
								}
								if newEntryPosition < 2 {
									obj.Scoreboard.Second.Name = obj.Scoreboard.First.Name
									obj.Scoreboard.Second.Score = obj.Scoreboard.First.Score
								}

								switch newEntryPosition {
								case 1:
									obj.Scoreboard.First.Name = newEntryName
									obj.Scoreboard.First.Score = score
								case 2:
									obj.Scoreboard.Second.Name = newEntryName
									obj.Scoreboard.Second.Score = score
								case 3:
									obj.Scoreboard.Third.Name = newEntryName
									obj.Scoreboard.Third.Score = score
								case 4:
									obj.Scoreboard.Fourth.Name = newEntryName
									obj.Scoreboard.Fourth.Score = score
								case 5:
									obj.Scoreboard.Fifth.Name = newEntryName
									obj.Scoreboard.Fifth.Score = score
								}
							}

							if needReorder > 1 {
								done := false
								for i := needReorder; i > 1 && !done; i-- {
									switch i {
									case 5:
										if obj.Scoreboard.Fourth.Score < obj.Scoreboard.Fifth.Score {
											auxScore := obj.Scoreboard.Fourth.Score
											obj.Scoreboard.Fourth.Score = obj.Scoreboard.Fifth.Score
											obj.Scoreboard.Fifth.Score = auxScore

											auxName := obj.Scoreboard.Fourth.Name
											obj.Scoreboard.Fourth.Name = obj.Scoreboard.Fifth.Name
											obj.Scoreboard.Fifth.Name = auxName
										} else {
											done = true
										}
									case 4:
										if obj.Scoreboard.Third.Score < obj.Scoreboard.Fourth.Score {
											auxScore := obj.Scoreboard.Third.Score
											obj.Scoreboard.Third.Score = obj.Scoreboard.Fourth.Score
											obj.Scoreboard.Fourth.Score = auxScore

											auxName := obj.Scoreboard.Third.Name
											obj.Scoreboard.Third.Name = obj.Scoreboard.Fourth.Name
											obj.Scoreboard.Fourth.Name = auxName
										} else {
											done = true
										}
									case 3:
										if obj.Scoreboard.Second.Score < obj.Scoreboard.Third.Score {
											auxScore := obj.Scoreboard.Second.Score
											obj.Scoreboard.Second.Score = obj.Scoreboard.Third.Score
											obj.Scoreboard.Third.Score = auxScore

											auxName := obj.Scoreboard.Second.Name
											obj.Scoreboard.Second.Name = obj.Scoreboard.Third.Name
											obj.Scoreboard.Third.Name = auxName
										} else {
											done = true
										}
									case 2:
										if obj.Scoreboard.First.Score < obj.Scoreboard.Second.Score {
											auxScore := obj.Scoreboard.First.Score
											obj.Scoreboard.First.Score = obj.Scoreboard.Second.Score
											obj.Scoreboard.Second.Score = auxScore

											auxName := obj.Scoreboard.First.Name
											obj.Scoreboard.First.Name = obj.Scoreboard.Second.Name
											obj.Scoreboard.Second.Name = auxName
										} else {
											done = true
										}
									}
								}
							}

							newEntryPosition = 0
						}
					} else {
						if len(newEntryName) < 3 {
							upper := strings.ToUpper(s)
							switch upper {
							case "Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P",
								"A", "S", "D", "F", "G", "H", "J", "K", "L",
								"Z", "X", "C", "V", "B", "N", "M":
								newEntryName += upper
							}
						}
						if k == tcell.KeyBackspace && len(newEntryName) > 0 {
							newEntryName = newEntryName[:len(newEntryName)-1]
						}
					}
				} else {
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
							obj.EnterAudio.Play()
							isRunning = false
						}
					}
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
	isFireworkShowing = false
	obj.GameAudio.Stop()
}

// Funciones -----------------------------------------

func drawGame() {
	lib.Screen.Clear()

	// Fireworks
	if isFireworkShowing {
		for i, index := range fireworkFramesIndex {
			if index >= 0 {
				for iR, r := range fireworkFrames[index] {
					for iC, c := range r {
						if fireworkStyles[index][iR][iC] != nil {
							lib.Screen.Put(fireworksX[i]+iC, fireworksY[i]+iR, string(c), *fireworkStyles[index][iR][iC])
						}
					}
				}
			}
		}
	}

	// Nivel
	lib.DrawString(lib.Width/2-21, lib.Height/2-4, levelFrame, lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2-19, lib.Height/2-2, lib.AppendBlank(strconv.Itoa(int(level)), 3), lib.DefaultStyle)

	// Lineas
	lib.DrawString(lib.Width/2-25, lib.Height/2+1, linesFrame, lib.DefaultStyle)
	var linesString string
	if lines >= 1000 {
		linesString = lib.AppendBlank(strconv.Itoa(int((lines%1000000)/1000)), 3) + "." + lib.AppendBlank(strconv.Itoa(int(lines%1000)), 3)
	} else if lines >= 1 {
		linesString = "    " + lib.AppendBlank(strconv.Itoa(int(lines%1000)), 3)
	} else {
		linesString = "      0"
	}
	lib.Screen.PutStrStyled(lib.Width/2-23, lib.Height/2+3, linesString, lib.DefaultStyle)

	// Score
	lib.DrawString(lib.Width/2-29, lib.Height/2+6, scoreFrame, lib.DefaultStyle)
	var scoreString string
	if score >= 1000000 {
		scoreString = lib.AppendBlank(strconv.Itoa(int(score%1000000000)/1000000), 3) + "." + lib.AppendBlank(strconv.Itoa(int((score%1000000)/1000)), 3) + "." + lib.AppendBlank(strconv.Itoa(int(score%1000)), 3)
	} else if score >= 1000 {
		scoreString = "    " + lib.AppendBlank(strconv.Itoa(int((score%1000000)/1000)), 3) + "." + lib.AppendBlank(strconv.Itoa(int(score%1000)), 3)
	} else if score >= 1 {
		scoreString = "        " + lib.AppendBlank(strconv.Itoa(int(score%1000)), 3)
	} else {
		scoreString = "          0"
	}
	lib.Screen.PutStrStyled(lib.Width/2-27, lib.Height/2+8, scoreString, lib.DefaultStyle)

	// Scoreboard
	lib.DrawString(lib.Width/2+14, lib.Height/2+3, scoreboardFrame, lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+4, obj.Scoreboard.First.Name, lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+4, strconv.FormatInt(obj.Scoreboard.First.Score, 10), lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+5, obj.Scoreboard.Second.Name, lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+5, strconv.FormatInt(obj.Scoreboard.Second.Score, 10), lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+6, obj.Scoreboard.Third.Name, lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+6, strconv.FormatInt(obj.Scoreboard.Third.Score, 10), lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+7, obj.Scoreboard.Fourth.Name, lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+7, strconv.FormatInt(obj.Scoreboard.Fourth.Score, 10), lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+8, obj.Scoreboard.Fifth.Name, lib.DefaultStyle)
	lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+8, strconv.FormatInt(obj.Scoreboard.Fifth.Score, 10), lib.DefaultStyle)

	if newEntryPosition != 0 {
		lib.Screen.PutStrStyled(lib.Width/2+18, lib.Height/2+3, "NEW__ENTRY", lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2+14, lib.Height/2+3+int(newEntryPosition), ">", lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+3+int(newEntryPosition), "___", lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+3+int(newEntryPosition), newEntryName, lib.DefaultStyle)
		lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+3+int(newEntryPosition), strconv.FormatInt(score, 10), lib.DefaultStyle)

		if newEntryPosition < 5 {
			lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+8, obj.Scoreboard.Fourth.Name, lib.DefaultStyle)
			lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+8, strconv.FormatInt(obj.Scoreboard.Fourth.Score, 10), lib.DefaultStyle)
		}
		if newEntryPosition < 4 {
			lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+7, obj.Scoreboard.Third.Name, lib.DefaultStyle)
			lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+7, strconv.FormatInt(obj.Scoreboard.Third.Score, 10), lib.DefaultStyle)
		}
		if newEntryPosition < 3 {
			lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+6, obj.Scoreboard.Second.Name, lib.DefaultStyle)
			lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+6, strconv.FormatInt(obj.Scoreboard.Second.Score, 10), lib.DefaultStyle)
		}
		if newEntryPosition < 2 {
			lib.Screen.PutStrStyled(lib.Width/2+19, lib.Height/2+5, obj.Scoreboard.First.Name, lib.DefaultStyle)
			lib.Screen.PutStrStyled(lib.Width/2+25, lib.Height/2+5, strconv.FormatInt(obj.Scoreboard.First.Score, 10), lib.DefaultStyle)
		}
	}

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
		lib.DrawString(lib.Width/2-7, lib.Height/2-2, pauseFrame, lib.DefaultStyle)

		i := 0
		switch selectedButton {
		case 0:
			i = 0
		case 1:
			i = 1
		case 2:
			i = 2
		}
		lib.DrawString(lib.Width/2-5, lib.Height/2-1+i, []string{">"}, lib.DefaultStyle)
	}

	// Game over
	if isGameOver {
		lib.DrawString(lib.Width/2-7, lib.Height/2-2, gameOverFrame, lib.GameOverStyle)

		if newEntryPosition == 0 {
			alto := 0
			switch selectedButton {
			case 0:
				alto = 0
			case 1:
				alto = 1
			}
			lib.DrawString(lib.Width/2-5, lib.Height/2-1+alto, []string{">"}, lib.GameOverStyle)
		}
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

	newEntryPosition = 0
	newEntryName = ""

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

	for i := range 10 {
		fireworkFramesIndex[i] = -1
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
		channel := obj.GameOverAudio.PlayAndWait()

		if score > obj.Scoreboard.First.Score {
			newEntryPosition = 1
		} else if score > obj.Scoreboard.Second.Score {
			newEntryPosition = 2
		} else if score > obj.Scoreboard.Third.Score {
			newEntryPosition = 3
		} else if score > obj.Scoreboard.Fourth.Score {
			newEntryPosition = 4
		} else if score > obj.Scoreboard.Fifth.Score {
			newEntryPosition = 5
		}

		go func() {
			<-channel
			obj.FireworksAudio.Play()

			fireworksPosX := make([]int, 10)
			fireworksPosY := make([]int, 10)
			isFireworkShowing = true

			// Mostrar 10 cohetes
			for i, t := range fireworkTimes {

				// Espera
				time.Sleep(t)

				if !isFireworkShowing {
					break
				}

				// Buscar posición
				var x, y int

				isRightOnTop := true
				for isRightOnTop {
					isRightOnTop = false

					x = rand.Intn(90)
					y = rand.Intn(40)

					if (x >= 22 && x <= 46) && (y >= 3 && y <= 25) {
						isRightOnTop = true
						continue
					}

					if i != 0 {
						for j := int(math.Max(0, float64(i-5))); j < i-1; j++ {
							if (x >= fireworksPosX[j] && x <= fireworksPosX[j]+13) &&
								(y >= fireworksPosY[j] && y <= fireworksPosY[j]+13) {
								isRightOnTop = true
								continue
							}
						}
					}
				}

				fireworksPosX[i] = x
				fireworksPosY[i] = y

				go func() {
					index := i
					fireworksX[index] = lib.Width/2 + (x - 45) - 7
					fireworksY[index] = lib.Height/2 + (y - 20) - 7

					now := time.Now()
					for j := range fireworkFrames {
						if !isFireworkShowing {
							break
						}

						fireworkFramesIndex[index] = j
						drawGame()

						for time.Since(now) < 100*time.Millisecond {
							time.Sleep(10 * time.Millisecond)
						}
						now = time.Now()
					}
					fireworkFramesIndex[index] = -1
					drawGame()
					if index == 9 {
						isFireworkShowing = false
					}
				}()
			}
		}()
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

	levelFrame = []string{
		" LEVEL ",
		" _---_ ",
		"{     }",
		" ¯---¯ ",
	}
	linesFrame = []string{
		"     LINES ",
		" _-------_ ",
		"{         }",
		" ¯-------¯ ",
	}
	scoreFrame = []string{
		"         SCORE ",
		" _-----------_ ",
		"{             }",
		" ¯-----------¯ ",
	}

	scoreboardFrame = []string{
		"____SCOREBOARD____",
		"  1°     :",
		"  2°     :",
		"  3°     :",
		"  4°     :",
		"  5°     :",
		"¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯",
	}
	gameOverFrame = []string{
		"┌=┤GAMEOVER├=┐",
		"║   RESTART  ║",
		"║   EXIT     ║",
		"└============┘",
	}
	pauseFrame = []string{
		"┌===┤MENU├===┐",
		"║   RESUME   ║",
		"║   RESTART  ║",
		"║   EXIT     ║",
		"└============┘",
	}

	fireworkFramesIndex = make([]int, 10)
	fireworksX          = make([]int, 10)
	fireworksY          = make([]int, 10)

	fireworkTimes = [10]time.Duration{ // 10
		100 * time.Millisecond,  // 100
		1600 * time.Millisecond, // 1.700
		200 * time.Millisecond,  // 1.900
		180 * time.Millisecond,  // 2.080
		200 * time.Millisecond,  // 2.280
		220 * time.Millisecond,  // 2.500
		300 * time.Millisecond,  // 2.800
		1100 * time.Millisecond, // 3.900
		1700 * time.Millisecond, // 5.600
		3300 * time.Millisecond, // 8.900
	}

	fireworkFrames = [][]string{
		{
			"",
			"",
			"",
			"",
			"      X",
		}, {
			"",
			"",
			"",
			"     \\|/",
			"    > O <",
			"     /|\\",
		}, {
			"",
			"",
			"   \\  |  /",
			"    \\ | /",
			"  ---{O}---",
			"    / | \\",
			"   /  |  \\",
		}, {
			"",
			"   _-===-_",
			"  \\ .   . /",
			" \\  . + .  /",
			" | +  X  + |",
			" /  ' + '  \\",
			"  / '   ' \\",
			"    -===- ",
		}, {
			"",
			"   _-===-_",
			"  \\ .   . /",
			" \\  . + .  /",
			" | +  X  + |",
			" /  ' + '  \\",
			"  / '   ' \\",
			"    -===- ",
		}, {
			"",
			"   _-===-_",
			"  \\ .   . /",
			" \\  . + .  /",
			" | +  X  + |",
			" /  ' + '  \\",
			"  / '   ' \\",
			"    -===- ",
		}, {
			"   |  -      ",
			" -  -   |  - ",
			"     |    -  ",
			"  |        | ",
			"-           |",
			"  |      |   ",
			"       -  |  ",
			"-   | -   -  ",
			"   -    |    ",
		}, {
			"   -  |      ",
			" |  |   -  | ",
			"     -    |  ",
			"  -        - ",
			"|           -",
			"  -      -   ",
			"       |  -  ",
			"|   - |   |  ",
			"   |    -    ",
		}, {
			"",
			"   |  -      ",
			" -  -   |  - ",
			"     |    -  ",
			"  |        | ",
			"-           |",
			"  |      |   ",
			"       -  |  ",
			"-   | -   -  ",
			"   -    |    ",
		}, {
			"",
			"   -  |      ",
			" |  |   -  | ",
			"     -    |  ",
			"  -        - ",
			"|           -",
			"  -      -   ",
			"       |  -  ",
			"|   - |   |  ",
			"   |    -    ",
		}, {
			"",
			"",
			"   |  -      ",
			" -  -   |  - ",
			"     |    -  ",
			"  |        | ",
			"-           |",
			"  |      |   ",
			"       -  |  ",
			"-   | -   -  ",
			"   -    |    ",
		}, {
			"",
			"",
			"   -  |      ",
			" |  |   -  | ",
			"     -    |  ",
			"  -        - ",
			"|           -",
			"  -      -   ",
			"       |  -  ",
			"|   - |   |  ",
			"   |    -    ",
		}, {
			"",
			"",
			"",
			"   |  -      ",
			" -  -   |  - ",
			"          -  ",
			"  |        | ",
			"-           |",
			"  |          ",
			"       -  |  ",
			"      -   -  ",
			"   -    |    ",
		}, {
			"",
			"",
			"",
			"   -  |      ",
			" |  |   -  | ",
			"          |  ",
			"  -        - ",
			"|           -",
			"  -          ",
			"       |  -  ",
			"      |   |  ",
			"   |    -    ",
		}, {
			"",
			"",
			"",
			"",
			"   |  -      ",
			" -  -   |  - ",
			"             ",
			"  |        | ",
			"            |",
			"  |          ",
			"       -     ",
			"          -  ",
			"   -    |    ",
		}, {
			"",
			"",
			"",
			"",
			"   -  |      ",
			" |  |   -  | ",
			"             ",
			"  -        - ",
			"            -",
			"  -          ",
			"       |     ",
			"          |  ",
			"   |    -    ",
		},
	}

	fireworkStyles = [][][]*tcell.Style{
		{
			nil,
			nil,
			nil,
			nil,
			{n, n, n, n, n, n, &r},
		}, {
			nil,
			nil,
			nil,
			{n, n, n, n, n, &p, &g, &b},
			{n, n, n, n, &w, n, &r, n, &p},
			{n, n, n, n, n, &p, &w, &g},
		}, {
			nil,
			nil,
			{n, n, n, &p, n, n, &g, n, n, &b},
			{n, n, n, n, &b, n, &w, n, &p},
			{n, n, &w, &g, &p, &r, &r, &r, &b, &g, &p},
			{n, n, n, n, &p, n, &w, n, &b},
			{n, n, n, &p, n, n, &g, n, n, &w},
		}, {
			nil,
			{n, n, n, &p, &g, &b, &w, &p, &g, &w},
			{n, n, &p, n, &b, n, n, n, &p, n, &g},
			{n, &b, n, n, &p, n, &g, n, &w, n, n, &b},
			{n, &w, n, &p, n, n, &r, n, n, &g, n, &p},
			{n, &p, n, n, &w, n, &p, n, &w, n, n, &b},
			{n, n, &p, n, &b, n, n, n, &p, n, &g},
			{n, n, n, &p, &g, &b, &p, &b, &p, &w},
		}, {
			nil,
			{n, n, n, &b, &w, &p, &g, &b, &w, &g},
			{n, n, &b, n, &p, n, n, n, &b, n, &w},
			{n, &p, n, n, &b, n, &w, n, &g, n, n, &p},
			{n, &g, n, &b, n, n, &r, n, n, &w, n, &b},
			{n, &b, n, n, &g, n, &b, n, &g, n, n, &p},
			{n, n, &b, n, &p, n, n, n, &b, n, &w},
			{n, n, n, &b, &w, &p, &b, &p, &b, &g},
		}, {
			nil,
			{n, n, n, &p, &g, &b, &w, &p, &g, &w},
			{n, n, &p, n, &b, n, n, n, &p, n, &g},
			{n, &b, n, n, &p, n, &g, n, &w, n, n, &b},
			{n, &w, n, &p, n, n, &r, n, n, &g, n, &p},
			{n, &p, n, n, &w, n, &p, n, &w, n, n, &b},
			{n, n, &p, n, &b, n, n, n, &p, n, &g},
			{n, n, n, &p, &g, &b, &p, &b, &p, &w},
		}, {
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			nil,
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			nil,
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			nil,
			nil,
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		}, {
			nil,
			nil,
			nil,
			nil,
			{n, n, n, &p, n, n, &g, n, n, n, n, n, n},
			{n, &b, n, n, &w, n, n, n, &r, n, n, &p, n},
			{n, n, n, n, n, &b, n, n, n, n, &w, n, n},
			{n, n, &r, n, n, n, n, n, n, n, n, &g, n},
			{&p, n, n, n, n, n, n, n, n, n, n, n, &w},
			{n, n, &r, n, n, n, n, n, n, &g, n, n, n},
			{n, n, n, n, n, n, n, &b, n, n, &p, n, n},
			{&r, n, n, n, &g, n, &b, n, n, n, &w, n, n},
			{n, n, n, &p, n, n, n, n, &g, n, n, n, n},
		},
	}
)

var n *tcell.Style = nil
var r = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 0, 0))
var g = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 255, 0))
var b = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 0, 255))
var w = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 255, 255))
var p = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 0, 255))
