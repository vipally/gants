package gants

import (
	"fmt"
	"strings"
	"time"
)

func nowTimestamp() int64 {
	return timestamp(time.Now())
}

func timestamp(t time.Time) int64 {
	return t.Unix()*1e3 + int64(t.Nanosecond()/1e6)
}

func newTimeHeap(size int) *timeHeap {
	h := &timeHeap{
		list: make([]*task, 0, size),
		tm:   time.NewTimer(time.Second),
	}
	h.tm.Stop()
	return h
}

type timeHeap struct {
	tm   *time.Timer
	list []*task
}

func (h *timeHeap) Push(s *task) error {
	h.list = append(h.list, s)
	if err := h.adjustUp(h.list, len(h.list)-1, s); err != nil {
		h.list = h.list[:h.Size()-1]
		return err
	}
	h.tm.Reset(h.TopDuration())
	return nil
}

func (h *timeHeap) WaitTimer() {
	<-h.tm.C
}

func (h *timeHeap) PopTimeout() ([]*task, bool) {
	// NOTE:
	// If PopTimeout is trigger by WaitTimer()
	// It is necessary to check if h.TopDuration()<=0 here.
	// Because if system clock is reset back during sleep,
	// it is probbly h.TopDuration()>0 when wake up.
	// In that case, it is necessary to sleep again until timeout.
	var ts []*task
	for {
		if dur := h.TopDuration(); dur <= 0 {
			if t, ok := h.Pop(); ok {
				ts = append(ts, t)
			}
		} else {
			break
		}
	}

	if dur := h.TopDuration(); dur > 0 {
		h.tm.Reset(dur)
	}
	return ts, len(ts) > 0
}

func (h *timeHeap) Pop() (*task, bool) {
	if h.Empty() {
		return nil, false
	}
	ret := h.list[0]
	if s := h.Size() - 1; s > 0 {
		h.list[0], h.list[s] = h.list[s], h.list[0]
		h.adjustDown(h.list[:h.Size()-1], 0, h.list[0])
	}
	h.list = h.list[:h.Size()-1]

	return ret, true
}

// adjust heap to select a proper hole to set v
func (h *timeHeap) adjustDown(b []*task, hole int, v *task) {
	// adjust heap to select a proper hole to set v
	for l := h.lchild(hole); l < len(b); l = h.lchild(hole) {
		c := l // index that need compare with hole

		// let the most proper child to compare with v
		if r := l + 1; r < len(b) && !h.cmp(b[r], b[l]) {
			c = r
		}

		if h.cmp(b[c], v) { //v is the most proper root, finish adjust
			break
		}

		//c is the most proper root, swap with hole, and continue adjust
		b[hole], hole = b[c], c
	}
	b[hole] = v // put v to last hole
}

// adjust heap to select a proper hole to set v
func (h *timeHeap) adjustUp(b []*task, hole int, v *task) error {
	for hole > 0 {
		if parent := h.parent(hole); !h.cmp(v, b[parent]) {
			b[hole], hole = b[parent], parent
		} else {
			break
		}
	}
	b[hole] = v //put v to last hole
	return nil
}

func (h *timeHeap) cmp(c, p *task) bool {
	return p.timestamp <= c.timestamp
}

// get parent index
func (h *timeHeap) parent(idx int) int {
	return (idx - 1) / 2
}

// get left child index
func (h *timeHeap) lchild(idx int) int {
	return 2*idx + 1
}

// TopDuration return duration from now to top timer
func (h *timeHeap) TopDuration() time.Duration {
	if t := h.Top(); t > 0 {
		return time.Duration(nowTimestamp()-t) * time.Millisecond
	}
	return time.Hour
}

func (h *timeHeap) Top() int64 {
	if h.Empty() {
		return -1
	}
	return h.list[0].timestamp
}

func (h *timeHeap) Size() int {
	return len(h.list)
}

func (h *timeHeap) Empty() bool {
	return h.Size() == 0
}

// check if b is a valid heap
func (h *timeHeap) CheckHeap() bool {
	for i := h.Size() - 1; i > 0; i-- {
		p := h.parent(i)
		if !h.cmp(h.list[i], h.list[p]) {
			return false
		}
	}
	return true
}

func (h *timeHeap) Show() string {
	var b strings.Builder
	b.WriteString("-------\n")
	for i, v := range h.list {
		b.WriteString(fmt.Sprintf("%d ", v.timestamp))
		if (i+1)&(i+2) == 0 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
