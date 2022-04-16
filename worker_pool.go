// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// License that can be found in the LICENSE file.

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
