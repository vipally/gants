package gants

import (
	"sync"
	"sync/atomic"
)

func NewPool(options ...Option) *Pool {
	opts := loadOptions(options...)
	return &Pool{
		opts: opts,
	}
}

type Pool struct {
	lk     sync.Mutex
	wCond  *sync.Cond
	chTask chan *Task
	tq     taskQueue // task queue buffer if chTak is full
	stat   stat
	opts   *Options

	workerIDGen int32
}

func (p *Pool) PushTask(t *Task) (err error) {
	if len(p.chTask) < cap(p.chTask) {
		p.chTask <- t
	} else {
		err = p.tq.Push(t)
	}
	p.wCond.Signal()
	return
}

func (p *Pool) Stop() {

}

func (p *Pool) Report() {

}

func (p *Pool) spawnWorker() {
	w := &goWorker{
		p:  p,
		id: int(atomic.AddInt32(&p.workerIDGen, 1)),
	}
	go w.run()
}
