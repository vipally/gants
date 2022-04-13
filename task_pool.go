package gants

import (
	"sync"
	"sync/atomic"
)

type taskPool struct {
	lock sync.Mutex

	pool  []*task
	idGen uint64
}

func (p *taskPool) Acquire(f func()) *task {
	p.lock.Lock()
	defer p.lock.Unlock()

	var r *task
	if size := len(p.pool); size > 0 {
		r = p.pool[size-1]
		p.pool = p.pool[:size-1]
	} else {
		r = &task{}
	}

	r.id = p.NextID()
	r.f = f
	return r
}

func (p *taskPool) Recycle(t *task) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pool = append(p.pool, t.clean())
}

func (p *taskPool) NextID() TaskID {
	return atomic.AddUint64(&p.idGen, 1)
}
