package gants

import (
	"sync/atomic"
	"time"
)

// newStopwatch create a stopwatch
func newStopwatch() *stopwatch {
	p := &stopwatch{
		startTime: time.Now().UnixNano(),
	}
	return p
}

// goroutine-safe time stopwatch
type stopwatch struct {
	startTime int64
}

// StartTime return start time of this stopwatch.
func (w *stopwatch) StartTime() int64 {
	return atomic.LoadInt64(&w.startTime)
}

// Start reset start time of this stopwatch.
func (w *stopwatch) Start() int64 {
	s := time.Now().UnixNano()
	atomic.StoreInt64(&w.startTime, s)
	return s
}

func (w *stopwatch) duration(reset bool) time.Duration {
	t := time.Now().UnixNano()
	dur := time.Duration(t - atomic.LoadInt64(&w.startTime))
	if reset {
		atomic.StoreInt64(&w.startTime, t)
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
