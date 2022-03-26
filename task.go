package gants

import (
	"sync/atomic"
)

type TaskID = uint64

type task struct {
	f  func()
	id uint64
}

// func (t *task) Execute() {
// 	t.f()
// }

type taskPool struct {
	pool  []*task
	idGen uint64
}

func (p *taskPool) Acquire(f func()) *task {
	id := p.NextID()
	return &task{
		f:  f,
		id: id,
	}
}

func (p *taskPool) NextID() TaskID {
	return atomic.AddUint64(&p.idGen, 1)
}

func (p *taskPool) Recycle(t *task) {

}

func (p *taskPool) Clean() {

}
