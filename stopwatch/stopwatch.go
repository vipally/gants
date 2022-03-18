//Package stopwatch implements some easy way to mesure time
package stopwatch

import (
	"bytes"
	"fmt"
	"time"
)

func New(options ...Option) *StopWatch {
	p := &StopWatch{}

	opts := loadOptions(options...)
	if opts.StepCacheCap > 0 {
		p.durList = make([]WatchDur, 0, opts.StepCacheCap)
	}
	if opts.AutoStart {
		p.Start()
	}
	return p
}

type WatchDur struct {
	name  string
	dur   time.Duration
	times uint
}

func (d *WatchDur) Init(name string, dur time.Duration, times uint) {
	if times <= 0 { //0 is forbinden
		times = 1
	}
	d.name = name
	d.dur = dur
	d.times = times
}

func (d *WatchDur) String() string {
	return fmt.Sprintf("%s\t%s\t%s", d.name, d.dur, d.dur/time.Duration(d.times))
}

func (d *WatchDur) Report(idx int, lastDur time.Duration) string {
	stepDur := d.dur - lastDur
	atomicTime := lastDur / time.Duration(d.times)

	return fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d", idx, d.name, d.dur, stepDur, atomicTime)
}

func (d *WatchDur) Duration() time.Duration {
	return d.dur
}

//StopWatch is a useful time mesure tool object
type StopWatch struct {
	startTime time.Time
	lastDur   time.Duration //last watch duration
	lastIdx   int           //last watch index
	durList   []WatchDur
}

func (w *StopWatch) Start() time.Time {
	w.startTime = time.Now()
	w.durList = nil
	w.lastDur = 0
	w.lastIdx = 0
	return w.startTime
}

func (w *StopWatch) Stop() time.Duration {
	t := time.Now()
	d := t.Sub(w.startTime)
	return d
}

func (w *StopWatch) Clear() time.Duration {
	t := time.Now()
	d := t.Sub(w.startTime)
	w.durList = nil
	w.lastDur = 0
	w.lastIdx = 0
	return d
}

func (w *StopWatch) AddStepWatch(name string, times uint) WatchDur {
	t := time.Now()
	var d WatchDur
	if times <= 0 { //0 is forbinden
		times = 1
	}
	d.Init(name, t.Sub(w.startTime), times)

	w.lastDur = d.dur
	w.lastIdx++
	w.durList = append(w.durList, d)

	return d
}

func (w *StopWatch) ReportOnce(name string) string {
	dur := time.Now().Sub(w.startTime)
	stepDur := dur - w.lastDur

	w.lastDur = dur
	w.lastIdx++

	return fmt.Sprintf("%s:%s", name, stepDur)
}

func (w *StopWatch) ReportWatch(name string, times uint) string {
	dur := time.Now().Sub(w.startTime)
	stepDur := dur - w.lastDur
	if times <= 0 { //0 is forbinden
		times = 1
	}

	w.lastDur = dur
	w.lastIdx++

	atomic_time := stepDur / time.Duration(times)

	return fmt.Sprintf("Watch%d\t%s\t%s\t%s\t%s\t%d", w.lastIdx, name, dur, stepDur, atomic_time, int64(atomic_time))
}

func (w *StopWatch) Count() int {
	return len(w.durList)
}

func (w *StopWatch) TellWatch(idx int) (d WatchDur) {
	if w.durList != nil {
		if idx >= 0 && idx < len(w.durList) {
			d = w.durList[idx]
		}
	}
	return
}

func (w *StopWatch) TellAllWatch() []WatchDur {
	return w.durList
}

func (w *StopWatch) Report() string {
	var buf bytes.Buffer
	buf.WriteString("StopWatch:\n")
	if nil != w.durList {
		lastDur := time.Duration(0)
		for i, v := range w.durList {
			buf.WriteString(v.Report(i+1, lastDur))
			buf.WriteString("\n")
			lastDur = v.dur
		}
	}

	return buf.String()
}

func (w *StopWatch) Restart() time.Duration {
	t := time.Now()
	d := t.Sub(w.startTime)
	w.durList = w.durList[:0]
	w.startTime = t
	w.lastDur = 0
	w.lastIdx = 0
	return d
}
