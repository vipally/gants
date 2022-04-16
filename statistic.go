// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// License that can be found in the LICENSE file.

package gants

import (
	"sync/atomic"
	"time"
)

type stat struct {
	CountStat       uint64
	MaxTaskTime     int64
	TooBusyCount    uint64
	FinishTaskCount uint64
	FreeTime        int64
	WorkingTime     int64
	BusyTime        int64

	//stat
	SumWorkerCount uint64
	SumListLen     uint64
	SumTaskCount   uint64
	SumTaskTime    uint64

	Status      uint32
	IdleWorkers int32
	CurWorkers  int32
	GoWorkers   int32
	CurListLen  uint32
	MaxWorkers  uint32
	MaxListLen  uint32

	sw stopwatch
	p  *Pool
}

func (s *stat) init(p *Pool) {
	s.p = p
	s.sw.Start()
}

func (s *stat) AddStat() {
	iwc := atomic.LoadInt32(&s.IdleWorkers)
	wc := int(atomic.LoadInt32(&s.CurWorkers))
	if iwc == 0 &&
		wc < s.p.opts.MaxWorkerCount {
		s.p.spawnWorker()
	}

	taskLen := s.p.taskLen()
	atomic.StoreUint32(&s.CurListLen, uint32(taskLen))
}

func (s *stat) AddWorker() {
}

func (s *stat) DecWorker() {
}

func (s *stat) AddTask(listLen int) {
}

func (s *stat) DecTask(listLen int) {
}

type StatResult struct {
	MaxListLen      uint32
	CountStat       uint32
	TooBusyCount    uint64
	FinishTaskCount uint64
	AvgWorkers      float32
	AvgListLen      float32
	FreeTime        time.Duration
	WorkingTime     time.Duration
	BusyTime        time.Duration
}

func (s *stat) Report() {
}
