package gants

import (
	"time"
)

func newStopwatch() *stopwatch {
	p := &stopwatch{}
	p.Start()
	return p
}

type stopwatch struct {
	startTime time.Time
}

func (w *stopwatch) StartTime() time.Time {
	return w.startTime
}

func (w *stopwatch) Start() time.Time {
	w.startTime = time.Now()
	return w.StartTime()
}

func (w *stopwatch) duration(reset bool) time.Duration {
	t := time.Now()
	dur := t.Sub(w.startTime)
	if reset {
		w.startTime = t
	}
	return dur
}

func (w *stopwatch) Restart() time.Duration {
	return w.duration(true)
}

func (w *stopwatch) Duration() time.Duration {
	return w.duration(false)
}
