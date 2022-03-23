package gants

import (
	"sync/atomic"
)

type TaskID = uint64

// Executer is the executable interface.
type Executer interface {
	Execute()
}

// ExecuterFunc is the func of Executer.
type ExecuterFunc func()

// Execute impliments Executer.
func (f ExecuterFunc) Execute() {
	f()
}

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
	id := atomic.AddUint64(&p.idGen, 1)
	return &task{
		f:  f,
		id: id,
	}
}

func (p *taskPool) Recycle(t *task) {

}

func (p *taskPool) Clean() {

}
