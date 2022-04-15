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

type workerStatus = int32

const (
	statusWorking workerStatus = iota
	statusIdle
	statusGo
)

type goWorker struct {
	p               *Pool
	statusStartTime int64
	taskStartTime   int64
	id              int32
	status          workerStatus
}

func (w *goWorker) Run() {
	var t *task
	atomic.StoreInt32(&w.status, statusWorking)
	atomic.StoreInt64(&w.statusStartTime, time.Now().UnixNano())
	for {
		t = nil
		select {
		case <-w.p.ChExit:
			return
		case t = <-w.p.chTask: //fetch task from channel
			if t == nil { // channel closed
				return
			}
		default:
		}
		if t == nil {
			t, _ = w.p.tq.Pop() // fetch task from task queue
		}
		if t != nil {
			atomic.StoreInt64(&w.taskStartTime, time.Now().UnixNano())
			t.f() // execute the task
			w.p.tp.Recycle(t)
		} else {
			atomic.StoreInt32(&w.status, statusIdle)
			atomic.StoreInt64(&w.statusStartTime, time.Now().UnixNano())
			w.p.wCond.L.Lock()
			w.p.wCond.Wait() // no task, wait
			w.p.wCond.L.Unlock()
			atomic.StoreInt32(&w.status, statusWorking)
			atomic.StoreInt64(&w.statusStartTime, time.Now().UnixNano())
		}
	}
}

func (w *goWorker) Go(f func(p *Pool)) {
	atomic.StoreInt32(&w.status, statusGo)
	atomic.StoreInt64(&w.statusStartTime, time.Now().UnixNano())
	f(w.p)
}
