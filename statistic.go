package gants

import (
	"time"
)

type stat struct {
	Status          int32
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
	CountStat      uint32

	sw *stopwatch
}

func (s *stat) AddStat(workerCount int, listLen int) {
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
