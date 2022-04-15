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
	"sync/atomic"
	"time"
)

type stat struct {
	CountStat       uint64
	Status          uint32
	IdleWorkers     int32
	CurWorkers      int32
	GoWorkers       int32
	CurListLen      uint32
	MaxWorkers      uint32
	MaxListLen      uint32
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
