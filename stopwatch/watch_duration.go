package stopwatch

import (
	"fmt"
)

type WatchDuration struct {
	name  string
	dur   Duration
	times uint
}

func (d *WatchDuration) Init(name string, dur Duration, times uint) {
	if times <= 0 { //0 is forbinden
		times = 1
	}
	d.name = name
	d.dur = dur
	d.times = times
}

func (d *WatchDuration) String() string {
	return fmt.Sprintf("%s\t%s\t%s", d.name, d.dur, d.dur/Duration(d.times))
}

func (d *WatchDuration) Report(idx int, lastDur Duration) string {
	stepDur := d.dur - lastDur
	atomicTime := lastDur / Duration(d.times)

	return fmt.Sprintf("%d\t%s\t%s\t%s\t%s", idx, d.name, d.dur, stepDur, atomicTime)
}

func (d *WatchDuration) Duration() Duration {
	return d.dur
}

func (d *WatchDuration) Name() string {
	return d.name
}
