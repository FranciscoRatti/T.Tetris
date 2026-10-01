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
