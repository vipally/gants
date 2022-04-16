// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// license that can be found in the LICENSE file.

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
