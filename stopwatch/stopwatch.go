//Package stopwatch implements some easy way to mesure time
package stopwatch

import (
	"bytes"
	"fmt"
	"time"
)

type Duration = time.Duration

func New(options ...Option) *StopWatch {
	p := &StopWatch{}

	opts := loadOptions(options...)
	if opts.StepCacheCap > 0 {
		p.durList = make([]WatchDuration, 0, opts.StepCacheCap)
	}
	p.Start()

	return p
}

//StopWatch is a useful time mesure tool object
type StopWatch struct {
	startTime time.Time
	lastDur   Duration //last watch duration
	durList   []WatchDuration
}

func (w *StopWatch) Start() time.Time {
	w.startTime = time.Now()
	w.durList = w.durList[:0]
	w.lastDur = 0
	return w.startTime
}

func (w *StopWatch) Stop() {
	w.startTime = time.Time{}
	w.durList = w.durList[:0]
	w.lastDur = 0
}

func (w *StopWatch) Clear() Duration {
	t := time.Now()
	d := t.Sub(w.startTime)
	w.durList = nil
	w.lastDur = 0
	return d
}

func (w *StopWatch) AddStepWatch(name string, times uint) WatchDuration {
	t := time.Now()
	var d WatchDuration
	if times <= 0 { //0 is forbinden
		times = 1
	}
	d.Init(name, t.Sub(w.startTime), times)

	w.lastDur = d.dur
	w.durList = append(w.durList, d)

	return d
}

func (w *StopWatch) ReportOnce(name string) string {
	dur := time.Now().Sub(w.startTime)
	stepDur := dur - w.lastDur

	w.lastDur = dur

	return fmt.Sprintf("%s:%s/%s", name, stepDur, dur)
}

func (w *StopWatch) ReportWatch(name string, times uint) string {
	dur := time.Now().Sub(w.startTime)
	stepDur := dur - w.lastDur
	if times <= 0 { //0 is forbinden
		times = 1
	}

	w.lastDur = dur

	atomicTime := stepDur / Duration(times)

	return fmt.Sprintf("Watch%d\t%s\t%s\t%s\t%s\t%d", len(w.durList), name, dur, stepDur, atomicTime, int64(atomicTime))
}

func (w *StopWatch) ReportAll() string {
	var buf bytes.Buffer
	buf.WriteString("StopWatch:\n")
	if nil != w.durList {
		lastDur := Duration(0)
		for i, v := range w.durList {
			buf.WriteString(v.Report(i+1, lastDur))
			buf.WriteString("\n")
			lastDur = v.dur
		}
	}

	return buf.String()
}

func (w *StopWatch) Count() int {
	return len(w.durList)
}

func (w *StopWatch) TellDuration() Duration {
	dur := time.Now().Sub(w.startTime)
	return dur
}

func (w *StopWatch) TellStepDuration() Duration {
	dur := time.Now().Sub(w.startTime)
	stepDur := dur - w.lastDur
	return stepDur
}

func (w *StopWatch) TellWatch(idx int) (d WatchDuration) {
	if idx >= 0 && idx < len(w.durList) {
		d = w.durList[idx]
	}
	return
}

func (w *StopWatch) TellAllWatch() []WatchDuration {
	return w.durList
}

// Restart restart the stop watch and return time duration from last start time
func (w *StopWatch) Restart() Duration {
	t := time.Now()
	d := t.Sub(w.startTime)
	w.durList = w.durList[:0]
	w.startTime = t
	w.lastDur = 0
	return d
}
