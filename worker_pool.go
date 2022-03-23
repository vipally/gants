package gants

import (
	"sync/atomic"
)

type workerPool struct {
	p           *Pool
	workerIDGen int32
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
		id: atomic.AddInt32(&wp.workerIDGen, 1),
	}
	return w
}
