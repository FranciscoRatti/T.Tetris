package obj

import "time"

type gameTimer struct {
	sleep     float32
	loop      uint
	isRunning bool
	run       func()
}

var Timer gameTimer
var lastExec time.Time
var firstExec time.Time

func (_ *gameTimer) Initialize(sleep float32, run func()) {
	Timer = gameTimer{sleep, 0, false, run}
	lastExec = time.Now()
	firstExec = time.Now()
}

func (t *gameTimer) Start() {
	t.isRunning = true

	go func() {
		for t.isRunning {
			if time.Since(lastExec) >= time.Duration(t.sleep)*time.Millisecond {
				t.run()
				t.loop++
				lastExec = time.Now()
			}
		}
	}()
}

func (t *gameTimer) Stop() {
	t.isRunning = false
}

func (t *gameTimer) Resume() {
	t.isRunning = true
	t.Start()
}

func (t *gameTimer) Restart(speed float32) {
	Timer.Stop()
	Timer.sleep = speed
	Timer.Start()
}

func (t *gameTimer) UpdateLastExec() {
	lastExec = time.Now()
}

func (t *gameTimer) ChangeSleep(sleep float32) {
	t.sleep = sleep
}
func (t *gameTimer) GetSleep() float32 {
	return t.sleep
}

func (t *gameTimer) GetSinceFirstExec() time.Duration {
	return time.Since(firstExec)
}
func (t *gameTimer) GetLoop() uint {
	return t.loop
}
