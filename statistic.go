package gants

import (
	"time"
)

type stat struct {
	CountStat       uint32
	Status          int32
	IdleWorkers     int32
	CurWorkers      uint32
	CurListLen      uint32
	MaxWorkers      uint32
	MaxListLen      uint32
	TooBusyCount    uint64
	FinishTaskCount uint64
	FreeTime        int64
	WorkingTime     int64
	BusyTime        int64

	//stat
	SumWorkerCount uint64
	SumListLen     uint64
	SumTaskCount   uint64

	sw stopwatch
}

func (s *stat) init() {
	s.sw.Start()
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
