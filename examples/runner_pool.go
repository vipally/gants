package main

import (
	"sync"
)

type runner struct {
	arg int64
}

func (r *runner) clean() *runner {
	return r
}

func (r *runner) run() {

}

type runnerPool struct {
	lock sync.Mutex

	pool  []*runner
	idGen uint64
}

func (p *runnerPool) Acquire(arg int64) *runner {
	p.lock.Lock()
	defer p.lock.Unlock()

	var r *runner
	if size := len(p.pool); size > 0 {
		r = p.pool[size-1]
		p.pool = p.pool[:size-1]
	} else {
		r = &runner{}
	}
	r.arg = arg

	return r
}

func (p *runnerPool) Recycle(t *runner) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pool = append(p.pool, t.clean())
}
