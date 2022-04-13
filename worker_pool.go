package gants

import (
	"sync/atomic"
)

type workerPool struct {
	p     *Pool
	idGen int32
}

func (wp *workerPool) Init(p *Pool) {
	wp.p = p
}

func (wp *workerPool) Acquire() *goWorker {
	return wp.Spawn()
}

func (wp *workerPool) Spawn() *goWorker {
	w := &goWorker{
		p:  wp.p,
		id: atomic.AddInt32(&wp.idGen, 1),
	}
	return w
}
