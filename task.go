// MIT License
//
// Copyright (c) 2022 @gxlb
//
// Authors:
//   Ally Dale <vipally@gamil.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gants

import (
	"sync"
	"sync/atomic"
)

// TaskID is id of task
type TaskID = uint64

type task struct {
	f         func()
	id        uint64
	timestamp int64
}

func (t *task) clean() *task {
	t.f = nil
	t.id = 0
	t.timestamp = 0
	return t
}

//------------------------------------------------------------------------------

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
