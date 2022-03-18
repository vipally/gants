package gants

import (
	"time"

	"github.com/gxlb/gants/stopwatch"
)

type stat struct {
	Status          int32
	CurWorkers      uint32
	MaxWorkers      uint32
	MaxListLen      uint32
	TooBusyCount    uint64
	FinishTaskCount uint64
	FreeTime        time.Duration
	BusyTime        time.Duration

	sw *stopwatch.StopWatch

	//stat
	SumWorkerCount uint64
	SumListLen     uint64
	CountStat      uint32
	// //caculate values
	// TotalTime  time.Duration
}

type StatResult struct {
	MaxListLen      uint32
	CountStat       uint32
	TooBusyCount    uint64
	FinishTaskCount uint64
	AvgWorkers      float32
	AvgListLen      float32
}

func (s *stat) AddStat(workerCount int, listLen int) {
}

func (s *stat) Report() {

}
