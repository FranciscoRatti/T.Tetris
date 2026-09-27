package obj

import (
	"TTetris/src/lib"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v3"
)

// Config

type config struct {
	Game struct {
		Shadow bool `json:"shadow"`
		Hold   bool `json:"hold"`
	} `json:"game"`
	KeyBinding struct {
		Right  string `json:"right"`
		Down   string `json:"down"`
		Left   string `json:"left"`
		Rotate string `json:"rotate"`
		Hold   string `json:"hold"`
		Floor  string `json:"floor"`
		Pause  string `json:"pause"`
		Mute   string `json:"mute"`
	} `json:"key_binding"`
	Volume struct {
		Music struct {
			Value float64 `json:"value"`
			Mute  bool    `json:"mute"`
		} `json:"music"`
		Effects struct {
			Value float64 `json:"value"`
			Mute  bool    `json:"mute"`
		} `json:"effects"`
	} `json:"volume"`
}

var Config = config{}

func ReadConfig() {
	file, err := os.Open(lib.CONFIG_PATH + "config.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewDecoder(file).Decode(&Config); err != nil {
		log.Fatal(err)
	}

	// Key Binding
	KeyRight = ParseKey(Config.KeyBinding.Right)
	KeyDown = ParseKey(Config.KeyBinding.Down)
	KeyLeft = ParseKey(Config.KeyBinding.Left)
	KeyRotate = ParseKey(Config.KeyBinding.Rotate)
	KeyHold = ParseKey(Config.KeyBinding.Hold)
	KeyFloor = ParseKey(Config.KeyBinding.Floor)
	KeyPause = ParseKey(Config.KeyBinding.Pause)
	KeyMute = ParseKey(Config.KeyBinding.Mute)

	// Sonido
	if Config.Volume.Effects.Value > 6 {
		Config.Volume.Effects.Value = 6
	} else if Config.Volume.Effects.Value < 0.25 {
		Config.Volume.Effects.Value = 0.25
	}

	if Config.Volume.Music.Value > 6 {
		Config.Volume.Music.Value = 6
	} else if Config.Volume.Music.Value < 0.25 {
		Config.Volume.Music.Value = 0.25
	}

	if err = file.Close(); err != nil {
		log.Fatal(err)
	}
}

func WriteConfig() {
	bytes, err := json.MarshalIndent(Config, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(lib.CONFIG_PATH+"config.json", bytes, 0644); err != nil {
		log.Fatal(err)
	}
}

// Keys

type Key struct {
	Key  tcell.Key
	Char string
	Name string
}

var (
	KeyRight  Key
	KeyDown   Key
	KeyLeft   Key
	KeyRotate Key
	KeyHold   Key
	KeyFloor  Key
	KeyPause  Key
	KeyMute   Key
)

var keyNames = func() map[string]tcell.Key {
	m := make(map[string]tcell.Key, len(tcell.KeyNames))
	for k, name := range tcell.KeyNames {
		m[strings.ToLower(name)] = k
	}
	return m
}()

func ParseKey(name string) Key {
	lower := strings.ToLower(name)

	switch lower {
	case "q", "w", "e", "r", "t", "y", "u", "i", "o", "p",
		"a", "s", "d", "f", "g", "h", "j", "k", "l", "ñ",
		"<", ">", "z", "x", "c", "v", "b", "n", "m", ",",
		";", ".", ":", "-", "_", " ":
		return Key{tcell.KeyRune, lower, lower}
	case "nil":
		return Key{tcell.KeyNUL, "", lower}
	default:
		return Key{keyNames[lower], "", lower}
	}
}

func (k1 Key) Equals(k2 tcell.Key, rune string) bool {
	if k1.Key == k2 && k1.Char == rune {
		return true
	}
	return false
}
