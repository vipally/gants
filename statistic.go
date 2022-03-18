package gants

import (
	"time"
)

type stat struct {
	CurWorkers      int32
	MaxWorkers      int32
	MaxListLen      int32
	TooBusyCount    uint64
	FinishTaskCount uint64
	FreeTime        time.Duration
	BusyTime        time.Duration

	//caculate values
	AvgWorkers float32
	AvgListLen float32
	TotalTime  time.Duration
}

func (s *stat) Report() {

}
