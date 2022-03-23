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
	e  Executer
	id uint64
}

func (t *task) Execute() {
	t.e.Execute()
}

type taskPool struct {
	pool  []*task
	idGen uint64
}

func (p *taskPool) Acquire(e Executer) *task {
	id := atomic.AddUint64(&p.idGen, 1)
	return &task{
		e:  e,
		id: id,
	}
}

func (p *taskPool) Recycle(t *task) {

}

func (p *taskPool) Clean() {

}
