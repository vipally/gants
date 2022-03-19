package gants

import (
	"time"
)

// newStopwatch create a stopwatch
func newStopwatch() *stopwatch {
	p := &stopwatch{
		startTime: time.Now(),
	}
	return p
}

type stopwatch struct {
	startTime time.Time
}

// StartTime return start time of this stopwatch.
func (w *stopwatch) StartTime() time.Time {
	return w.startTime
}

// Start reset start time of this stopwatch.
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

// Restart reset the start time and return duration since last start time.
func (w *stopwatch) Restart() time.Duration {
	return w.duration(true)
}

// Duration return duration since start time.
func (w *stopwatch) Duration() time.Duration {
	return w.duration(false)
}
