package obj

import (
	"TTetris/src/lib"
	"encoding/gob"
	"log"
	"os"
)

type scoreboard struct {
	First, Second, Third, Fourth, Fifth player
}

type player struct {
	Name  string
	Score int64
}

var Scoreboard = scoreboard{}

func (s scoreboard) AddEntry(name string, score int64, index byte) {
	needReorder := 0

	if Scoreboard.First.Name == name {
		if Scoreboard.First.Score < score {
			Scoreboard.First.Score = score
			needReorder = 1
		}
	} else if Scoreboard.Second.Name == name {
		if Scoreboard.Second.Score < score {
			Scoreboard.Second.Score = score
			needReorder = 2
		}
	} else if Scoreboard.Third.Name == name {
		if Scoreboard.Third.Score < score {
			Scoreboard.Third.Score = score
			needReorder = 3
		}
	} else if Scoreboard.Fourth.Name == name {
		if Scoreboard.Fourth.Score < score {
			Scoreboard.Fourth.Score = score
			needReorder = 4
		}
	} else if Scoreboard.Fifth.Name == name {
		if Scoreboard.Fifth.Score < score {
			Scoreboard.Fifth.Score = score
			needReorder = 5
		}
	} else {
		if index < 5 {
			Scoreboard.Fifth.Name = Scoreboard.Fourth.Name
			Scoreboard.Fifth.Score = Scoreboard.Fourth.Score
		}
		if index < 4 {
			Scoreboard.Fourth.Name = Scoreboard.Third.Name
			Scoreboard.Fourth.Score = Scoreboard.Third.Score
		}
		if index < 3 {
			Scoreboard.Third.Name = Scoreboard.Second.Name
			Scoreboard.Third.Score = Scoreboard.Second.Score
		}
		if index < 2 {
			Scoreboard.Second.Name = Scoreboard.First.Name
			Scoreboard.Second.Score = Scoreboard.First.Score
		}

		switch index {
		case 1:
			Scoreboard.First.Name = name
			Scoreboard.First.Score = score
		case 2:
			Scoreboard.Second.Name = name
			Scoreboard.Second.Score = score
		case 3:
			Scoreboard.Third.Name = name
			Scoreboard.Third.Score = score
		case 4:
			Scoreboard.Fourth.Name = name
			Scoreboard.Fourth.Score = score
		case 5:
			Scoreboard.Fifth.Name = name
			Scoreboard.Fifth.Score = score
		}
	}

	if needReorder > 1 {
		done := false
		for i := needReorder; i > 1 && !done; i-- {
			switch i {
			case 5:
				if Scoreboard.Fourth.Score < Scoreboard.Fifth.Score {
					auxScore := Scoreboard.Fourth.Score
					Scoreboard.Fourth.Score = Scoreboard.Fifth.Score
					Scoreboard.Fifth.Score = auxScore

					auxName := Scoreboard.Fourth.Name
					Scoreboard.Fourth.Name = Scoreboard.Fifth.Name
					Scoreboard.Fifth.Name = auxName
				} else {
					done = true
				}
			case 4:
				if Scoreboard.Third.Score < Scoreboard.Fourth.Score {
					auxScore := Scoreboard.Third.Score
					Scoreboard.Third.Score = Scoreboard.Fourth.Score
					Scoreboard.Fourth.Score = auxScore

					auxName := Scoreboard.Third.Name
					Scoreboard.Third.Name = Scoreboard.Fourth.Name
					Scoreboard.Fourth.Name = auxName
				} else {
					done = true
				}
			case 3:
				if Scoreboard.Second.Score < Scoreboard.Third.Score {
					auxScore := Scoreboard.Second.Score
					Scoreboard.Second.Score = Scoreboard.Third.Score
					Scoreboard.Third.Score = auxScore

					auxName := Scoreboard.Second.Name
					Scoreboard.Second.Name = Scoreboard.Third.Name
					Scoreboard.Third.Name = auxName
				} else {
					done = true
				}
			case 2:
				if Scoreboard.First.Score < Scoreboard.Second.Score {
					auxScore := Scoreboard.First.Score
					Scoreboard.First.Score = Scoreboard.Second.Score
					Scoreboard.Second.Score = auxScore

					auxName := Scoreboard.First.Name
					Scoreboard.First.Name = Scoreboard.Second.Name
					Scoreboard.Second.Name = auxName
				} else {
					done = true
				}
			}
		}
	}
}

func ReadScoreboard() {
	file, err := os.OpenFile(lib.VAR_PATH+"scoreboard.obj", os.O_RDONLY, 0644)
	if err != nil {
		WriteScoreboard()
		return
	}

	if err = gob.NewDecoder(file).Decode(&Scoreboard); err != nil {
		WriteScoreboard()
		return
	}

	if err = file.Close(); err != nil {
		log.Fatal(err)
	}
}

func WriteScoreboard() {
	file, err := os.OpenFile(lib.VAR_PATH+"scoreboard.obj", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if err = gob.NewEncoder(file).Encode(Scoreboard); err != nil {
		log.Fatal(err)
	}

	if err = file.Close(); err != nil {
		log.Fatal(err)
	}
}

func RestartScoreboard() {
	Scoreboard.First.Name = "   "
	Scoreboard.First.Score = 0
	Scoreboard.Second.Name = "   "
	Scoreboard.Second.Score = 0
	Scoreboard.Third.Name = "   "
	Scoreboard.Third.Score = 0
	Scoreboard.Fourth.Name = "   "
	Scoreboard.Fourth.Score = 0
	Scoreboard.Fifth.Name = "   "
	Scoreboard.Fifth.Score = 0
}
