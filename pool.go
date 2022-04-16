// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// license that can be found in the LICENSE file.

package gants

import (
	"sync"
	"time"
)

type Duration = time.Duration

func NewPool(options ...Option) *Pool {
	opts := loadOptions(options...)
	return &Pool{
		opts:   opts,
		ChExit: make(chan struct{}),
		chTask: make(chan *task, opts.MaxWorkerCount),
	}
}

// Pool is task pool
type Pool struct {
	ChExit chan struct{} // channel for notify exit for outside pool

	wCond  *sync.Cond //Idle worker condition variable
	chTask chan *task
	tq     taskQueue  // task queue buffer if chTak is full
	tp     taskPool   // pool of task object
	wp     workerPool // pool of workers
	stat   stat
	opts   *Options
}

// Push summit a task for scheduled worker.
func (p *Pool) Push(f func()) TaskID {
	t := p.tp.Acquire(f)
	var ok bool
	select {
	case p.chTask <- t:
		ok = true
	default:
	}
	if !ok {
		p.tq.Push(t)
	}

	p.stat.AddStat()
	p.wakeupWorker()

	return t.id
}

// PushWithTimeout push a task that with timeout
func (p *Pool) PushWithTimeout(f func(), timeout Duration) {

}

// PushDelay push a task that execute after duration
func (p *Pool) PushDelay(f func(), delay Duration) TaskID {
	return 0
}

// Go execute a heavy task directly by special worker without schedule.
func (p *Pool) Go(f func(p *Pool)) TaskID {
	w := p.wp.Spawn()
	id := p.tp.NextID()
	go w.Go(f)
	return id
}

func (p *Pool) Stop() {
	close(p.ChExit)
	close(p.chTask)
}

func (p *Pool) GracefulStop() {
	close(p.ChExit)
	close(p.chTask)
}

func (p *Pool) taskLen() int {
	return len(p.chTask) + p.tq.Len()
}

func (p *Pool) Report() {

}

func (p *Pool) wakeupWorker() {
	p.wCond.Signal()
}

func (p *Pool) spawnWorker() {
	w := p.wp.Spawn()
	go w.Run()
}
