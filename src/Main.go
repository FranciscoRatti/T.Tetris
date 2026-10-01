package main

import (
	"TTetris/src/lib"
	"TTetris/src/obj"
	"TTetris/src/scenes"
	"log"
	"os"
	"runtime"

	"github.com/gdamore/tcell/v3"
	"github.com/gopxl/beep/speaker"
)

// Main
func main() {

	// Variables
	length := len(os.Args)
	if length > 0 {
		for i := 0; i < length; i++ {
			switch os.Args[i] {
			case "--resources":
				if lib.RESOURCES_PATH != "" {
					log.Fatal("Error: Parámetro --resources duplicado")
				}
				i++
				if length < i {
					log.Fatal("Error: No se especifico un path luego de --resources")
				}

				lib.RESOURCES_PATH = os.Args[i]
			case "--config":
				if lib.CONFIG_PATH != "" {
					log.Fatal("Error: Parámetro --config duplicado")
				}
				i++
				if length < i {
					log.Fatal("Error: No se especifico un path luego de --config")
				}

				lib.CONFIG_PATH = os.Args[i]
			case "--var":
				if lib.VAR_PATH != "" {
					log.Fatal("Error: Parámetro --var duplicado")
				}
				i++
				if length < i {
					log.Fatal("Error: No se especifico un path luego de --var")
				}

				lib.VAR_PATH = os.Args[i]
			}
		}
	}

	if lib.RESOURCES_PATH == "" {
		if runtime.GOOS == "linux" {
			lib.RESOURCES_PATH = "/usr/share/ttetris/"
		} else if runtime.GOOS == "windows" {
			lib.RESOURCES_PATH = os.Getenv("ProgramFiles") + "\\TTetris\\data\\"
		}
	}

	if lib.CONFIG_PATH == "" {
		if runtime.GOOS == "linux" {
			lib.CONFIG_PATH = os.Getenv("HOME") + "/.config/ttetris/"
		} else if runtime.GOOS == "windows" {
			lib.CONFIG_PATH = os.Getenv("AppData") + "\\TTetris\\"
		}
	}

	if lib.VAR_PATH == "" {
		if runtime.GOOS == "linux" {
			lib.VAR_PATH = "/var/lib/ttetris/"
		} else if runtime.GOOS == "windows" {
			lib.VAR_PATH = os.Getenv("AppData") + "\\TTetris\\"
		}
	}

	// Init ------------------------------------------------------------------------------------------------------------
	obj.ReadConfig()
	obj.ReadScoreboard()
	obj.InitSpeakers()

	// Inicializar pantalla
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal("Error al tomar pantalla", err)
	}
	if err := s.Init(); err != nil {
		log.Fatal("Error al tomar pantalla", err)
	}

	lib.Width, lib.Height = s.Size()
	s.SetStyle(lib.DefaultStyle)
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
		os.Exit(0)
	}
	defer quit()

	lib.Screen = s

	// Abrir menu
	scenes.OpenMenu()

	// Final
	obj.WriteConfig()
	obj.WriteScoreboard()
	speaker.Close()
}
