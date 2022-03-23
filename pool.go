package gants

import (
	"sync"
)

func NewPool(options ...Option) *Pool {
	opts := loadOptions(options...)
	return &Pool{
		opts:   opts,
		ChExit: make(chan struct{}),
		chTask: make(chan *task, opts.MaxWorkerCount),
	}
}

type Pool struct {
	wCond  *sync.Cond
	ChExit chan struct{} // channel for notify exit for outside pool
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

	p.wCond.Signal()
	return t.id
}

// Go execute a heavy task directly by special worker without schedule.
func (p *Pool) Go(f func()) TaskID {
	return 0
}

func (p *Pool) Stop() {
	close(p.ChExit)
	close(p.chTask)
}

func (p *Pool) Report() {

}

func (p *Pool) spawnWorker() {
	w := p.wp.Spawn()
	go w.Run()
}
