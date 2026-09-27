package obj

import (
	"TTetris/src/lib"
	"log"
	"math"
	"os"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

func InitSpeakers() {
	err := speaker.Init(44100, 4410)
	if err != nil {
		log.Fatal("Error al tomar parlantes", err)
	}

	EnterAudio = newBuffer(lib.RESOURCES_PATH + "audio/enter.wav")
	FloorAudio = newBuffer(lib.RESOURCES_PATH + "audio/floor.wav")
	GameAudio = newLoop(lib.RESOURCES_PATH+"audio/game.wav", -1)
	GameOverAudio = newAudio(lib.RESOURCES_PATH + "audio/gameover.wav")
	KeyAudio = newBuffer(lib.RESOURCES_PATH + "audio/key.wav")
	LevelUpAudio = newAudio(lib.RESOURCES_PATH + "audio/levelup.wav")
	LineAudio = newAudio(lib.RESOURCES_PATH + "audio/line.wav")
	Line4Audio = newAudio(lib.RESOURCES_PATH + "audio/line4.wav")
	MenuAudio = newLoop(lib.RESOURCES_PATH+"audio/menu.wav", -1)
	MoveAudio = newBuffer(lib.RESOURCES_PATH + "audio/move.wav")
	PauseAudio = newAudio(lib.RESOURCES_PATH + "audio/pause.wav")
	ResumeAudio = newAudio(lib.RESOURCES_PATH + "audio/resume.wav")
	RotateAudio = newBuffer(lib.RESOURCES_PATH + "audio/rotate.wav")

	speaker.Play(GameAudio.Volume)
	speaker.Play(MenuAudio.Volume)

	ChangeMuteEffectsWithoutChange(Config.Volume.Effects.Mute)
	ChangeMuteMusicWithoutChange(Config.Volume.Music.Mute)
}

var (
	EnterAudio    buffer
	FloorAudio    buffer
	GameAudio     loop
	GameOverAudio audio
	KeyAudio      buffer
	LevelUpAudio  audio
	LineAudio     audio
	Line4Audio    audio
	MenuAudio     loop
	MoveAudio     buffer
	PauseAudio    audio
	ResumeAudio   audio
	RotateAudio   buffer
)

func ChangeEffectsVolume(value float64) {
	Config.Volume.Effects.Value += value
	correctValue := calculateCorrectValue(&Config.Volume.Effects.Value)
	if correctValue == -1 {
		return
	}

	speaker.Lock()

	EnterAudio.Volume.Volume = correctValue
	FloorAudio.Volume.Volume = correctValue
	KeyAudio.Volume.Volume = correctValue
	LevelUpAudio.Volume.Volume = correctValue
	LineAudio.Volume.Volume = correctValue
	Line4Audio.Volume.Volume = correctValue
	MoveAudio.Volume.Volume = correctValue
	PauseAudio.Volume.Volume = correctValue
	ResumeAudio.Volume.Volume = correctValue
	RotateAudio.Volume.Volume = correctValue

	speaker.Unlock()
}

func ChangeMusicVolume(value float64) {
	Config.Volume.Music.Value += value
	correctValue := calculateCorrectValue(&Config.Volume.Music.Value)
	if correctValue == -1 {
		return
	}

	speaker.Lock()

	GameAudio.Volume.Volume = correctValue
	GameOverAudio.Volume.Volume = correctValue
	MenuAudio.Volume.Volume = correctValue

	speaker.Unlock()
}

func calculateCorrectValue(value *float64) float64 {
	if *value < 0.01 {
		*value = 0.01
		return -1
	} else if *value > 6 {
		*value = 6
		return -1
	}

	return math.Log10(*value) / math.Log10(2)
}

func ChangeMuteEffects() {
	Config.Volume.Effects.Mute = !Config.Volume.Effects.Mute
	ChangeMuteEffectsWithoutChange(Config.Volume.Effects.Mute)
}

func ChangeMuteEffectsWithoutChange(value bool) {
	speaker.Lock()

	EnterAudio.Volume.Silent = value
	FloorAudio.Volume.Silent = value
	KeyAudio.Volume.Silent = value
	LevelUpAudio.Volume.Silent = value
	LineAudio.Volume.Silent = value
	Line4Audio.Volume.Silent = value
	MoveAudio.Volume.Silent = value
	PauseAudio.Volume.Silent = value
	ResumeAudio.Volume.Silent = value
	RotateAudio.Volume.Silent = value

	speaker.Unlock()
}

func ChangeMuteMusic() {
	Config.Volume.Music.Mute = !Config.Volume.Music.Mute
	ChangeMuteMusicWithoutChange(Config.Volume.Music.Mute)
}

func ChangeMuteMusicWithoutChange(value bool) {
	speaker.Lock()

	GameAudio.Volume.Silent = value
	GameOverAudio.Volume.Silent = value
	MenuAudio.Volume.Silent = value

	speaker.Unlock()
}

// Audio ------------------------------------------------------

type audio struct {
	Audio  *beep.StreamSeekCloser
	Volume *effects.Volume
}

func newAudio(path string) audio {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, _, err := wav.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	return audio{&streamer, volume}
}

func (a audio) Play() {
	if err := (*a.Audio).Seek(0); err != nil {
		log.Fatal(err)
	}

	speaker.Play(a.Volume)
}

// Buffered audio ------------------------------------------------------

type buffer struct {
	Audio  *beep.Buffer
	Volume *effects.Volume
}

func newBuffer(path string) buffer {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	buffered := beep.NewBuffer(format)
	buffered.Append(streamer)
	if err = streamer.Close(); err != nil {
		log.Fatal(err)
	}

	return buffer{buffered, volume}
}

func (a buffer) Play() {
	volume := &effects.Volume{
		Streamer: a.Audio.Streamer(0, a.Audio.Len()),
		Base:     2,
		Volume:   math.Log10(Config.Volume.Effects.Value) / math.Log10(2),
		Silent:   a.Volume.Silent,
	}
	speaker.Play(volume)
}

// Loop ------------------------------------------------------

type loop struct {
	Audio  *beep.Ctrl
	Volume *effects.Volume
	Count  int
}

func newLoop(path string, count int) loop {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, _, err := wav.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	control := &beep.Ctrl{Streamer: beep.Loop(count, streamer), Paused: true}

	volume := &effects.Volume{
		Streamer: control,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	return loop{control, volume, count}
}

func (a loop) Play() {
	speaker.Lock()
	a.Audio.Paused = false
	speaker.Unlock()
}

func (a loop) Stop() {
	speaker.Lock()
	a.Audio.Paused = true
	speaker.Unlock()
}
