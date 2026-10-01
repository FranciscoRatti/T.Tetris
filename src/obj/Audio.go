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

var (
	EnterAudio     buffer
	FireworksAudio audio
	FloorAudio     buffer
	GameAudio      loop
	GameOverAudio  audio
	KeyAudio       buffer
	LevelUpAudio   audio
	LineAudio      audio
	Line4Audio     audio
	MenuAudio      loop
	MoveAudio      buffer
	PauseAudio     audio
	ResumeAudio    audio
	RotateAudio    buffer
)

func InitSpeakers() {
	err := speaker.Init(44100, 4410)
	if err != nil {
		log.Fatal("Error al tomar parlantes", err)
	}

	EnterAudio = newBuffer(lib.RESOURCES_PATH + "audio/enter.wav")
	FireworksAudio = newAudio(lib.RESOURCES_PATH + "audio/fireworks.wav")
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

	ChangeEffectsVolumeWithoutChange(calculateCorrectValue(&Config.Volume.Effects.Value))
	ChangeMuteEffectsWithoutChange(Config.Volume.Effects.Mute)
	ChangeMusicVolumeWithoutChange(calculateCorrectValue(&Config.Volume.Music.Value))
	ChangeMuteMusicWithoutChange(Config.Volume.Music.Mute)

	speaker.Play(GameAudio.Volume)
	speaker.Play(MenuAudio.Volume)
}


// Audio ---------------------------------------------------------------------------------------------------------------

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
func (a audio) PlayAndWait() chan bool {
	if err := (*a.Audio).Seek(0); err != nil {
		log.Fatal(err)
	}

	channel := make(chan bool)

	go func() {
		speaker.PlayAndWait(a.Volume)
		channel <- true
	}()

	return channel
}

// Buffer --------------------------------------------------------------------------------------------------------------

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

// Loop ----------------------------------------------------------------------------------------------------------------

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

// Volumen -------------------------------------------------------------------------------------------------------------

func ChangeEffectsVolume(value float64) {
	Config.Volume.Effects.Value += value
	correctValue := calculateCorrectValue(&Config.Volume.Effects.Value)
	if correctValue != -1 {
		ChangeEffectsVolumeWithoutChange(correctValue)
	}
}

func ChangeEffectsVolumeWithoutChange(value float64) {
	speaker.Lock()

	EnterAudio.Volume.Volume = value
	FloorAudio.Volume.Volume = value
	KeyAudio.Volume.Volume = value
	LevelUpAudio.Volume.Volume = value
	LineAudio.Volume.Volume = value
	Line4Audio.Volume.Volume = value
	MoveAudio.Volume.Volume = value
	PauseAudio.Volume.Volume = value
	ResumeAudio.Volume.Volume = value
	RotateAudio.Volume.Volume = value

	speaker.Unlock()
}

func ChangeMusicVolume(value float64) {
	Config.Volume.Music.Value += value
	correctValue := calculateCorrectValue(&Config.Volume.Music.Value)
	if correctValue != -1 {
		ChangeMusicVolumeWithoutChange(correctValue)
	}

}

func ChangeMusicVolumeWithoutChange(value float64) {
	speaker.Lock()

	FireworksAudio.Volume.Volume = value
	GameAudio.Volume.Volume = value
	GameOverAudio.Volume.Volume = value
	MenuAudio.Volume.Volume = value

	speaker.Unlock()
}

func calculateCorrectValue(value *float64) float64 {
	if *value < 0.15 {
		*value = 0.15
		return -1
	} else if *value > 3.75 {
		*value = 3.75
		return -1
	}

	return math.Log10(*value) / math.Log10(2)
}

// Mute ----------------------------------------------------------------------------------------------------------------

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

	FireworksAudio.Volume.Silent = value
	GameAudio.Volume.Silent = value
	GameOverAudio.Volume.Silent = value
	MenuAudio.Volume.Silent = value

	speaker.Unlock()
}

var mute = false

func ChangeMuteAll() {
	if mute {
		mute = false
		ChangeMuteEffectsWithoutChange(Config.Volume.Effects.Mute)
		ChangeMuteMusicWithoutChange(Config.Volume.Music.Mute)
	} else {
		mute = true
		ChangeMuteEffectsWithoutChange(true)
		ChangeMuteMusicWithoutChange(true)
	}
}