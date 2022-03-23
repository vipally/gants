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
	wCond  *sync.Cond
	chTask chan *task
	tq     taskQueue // task queue buffer if chTak is full
	tp     taskPool  // pool of task object
	stat   stat
	opts   *Options

	workerIDGen int32
}

func (p *Pool) Push(e Executer) TaskID {
	t := p.tp.Acquire(e)
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
