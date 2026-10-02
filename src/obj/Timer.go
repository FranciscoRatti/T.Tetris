package obj

import "time"

type timer struct {
	sleep     float32
	loop      uint
	isRunning bool
	run       func()
}

var Timer timer
var lastExec time.Time
var firstExec time.Time
var channel chan bool

func (_ *timer) Initialize(sleep float32, run func()) {
	channel = make(chan bool)
	Timer = timer{sleep, 0, false, run}
	lastExec = time.Now()
	firstExec = time.Now()
}

func (t *timer) Start() {
	t.isRunning = true

	go func() {
		for t.isRunning {
			if time.Since(lastExec) >= time.Duration(t.sleep)*time.Millisecond {
				t.run()
				t.loop++
				lastExec = time.Now()
			}
		}
		channel <- true
	}()
}

func (t *timer) Stop() {
	t.isRunning = false
}

func (t *timer) StopAndWait() {
	if t.isRunning {
		t.isRunning = false
		<-channel
	}
}

func (t *timer) Restart(speed float32) {
	Timer.StopAndWait()
	Timer.sleep = speed
	Timer.Start()
}

func (t *timer) UpdateLastExec() {
	lastExec = time.Now()
}

func (t *timer) ChangeSleep(sleep float32) {
	t.sleep = sleep
}
func (t *timer) GetSleep() float32 {
	return t.sleep
}

func (t *timer) GetSinceFirstExec() time.Duration {
	return time.Since(firstExec)
}
func (t *timer) GetLoop() uint {
	return t.loop
}
