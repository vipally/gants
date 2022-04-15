// MIT License
//
// Copyright (c) 2022 @gxlb
//
// Authors:
//   Ally Dale <vipally@gamil.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
