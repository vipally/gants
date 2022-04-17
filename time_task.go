// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// License that can be found in the LICENSE file.

package gants

import (
	"fmt"
	"strings"
	"time"
)

const timeUnit = time.Millisecond

func nowPrecisionTimestamp() int64 {
	return time.Now().UnixNano()
}

func nowTimestamp() int64 {
	return timestamp(time.Now())
}

func duration(when int64) Duration {
	return Duration(nowTimestamp()-when) * timeUnit
}

func timestamp(t time.Time) int64 {
	return t.Unix()*1e3 + int64(t.Nanosecond()/1e6)
}

func newTimeHeap(size int) *timeHeap {
	h := &timeHeap{
		b:  make([]*timeTask, 0, size),
		tm: time.NewTimer(time.Second),
	}
	h.tm.Stop()
	return h
}

type timeHeap struct {
	tm  *time.Timer
	b   []*timeTask
	tmp [8]*timeTask
}

func (h *timeHeap) PushDelay(s *timeTask, dur time.Duration) {
	s.when = nowTimestamp() + int64((dur+timeUnit-1)/timeUnit)
	h.Push(s)
}

func (h *timeHeap) PushAt(s *timeTask, t time.Time) {
	s.when = timestamp(t)
	h.Push(s)
}

func (h *timeHeap) Push(s *timeTask) {
	h.b = append(h.b, s)
	h.adjustUp(h.b, len(h.b)-1, s)
	if s == h.b[0] {
		h.tm.Reset(-duration(s.when))
	}
}

func (h *timeHeap) WaitTimeout() {
	<-h.tm.C
}

func (h *timeHeap) PopTimeout() ([]*timeTask, bool) {
	// NOTE:
	// PopTimeout is trigger by WaitTimeout()
	// It is necessary to check if h.Top()<=now here.
	// Because if system clock is reset back during sleep,
	// it is probbly h.Top()>now when wake up.
	// In that case, it is necessary to sleep again until timeout.
	var ts []*timeTask = h.tmp[:0]
	var dur time.Duration = -1
	for now := nowTimestamp(); ; {
		if t, ok := h.Top(); ok {
			if t <= now {
				if t, ok := h.Pop(); ok {
					ts = append(ts, t)
					continue
				}
			}
			dur = time.Duration(t-now) * timeUnit
		}

		break // end loop
	}

	if dur > 0 {
		h.tm.Reset(dur)
	}
	return ts, len(ts) > 0
}

func (h *timeHeap) Pop() (*timeTask, bool) {
	if h.Empty() {
		return nil, false
	}
	ret := h.b[0]
	if s := h.Size() - 1; s > 0 {
		h.b[0], h.b[s] = h.b[s], h.b[0]
		h.adjustDown(h.b[:h.Size()-1], 0, h.b[0])
	}
	h.b = h.b[:h.Size()-1]

	return ret, true
}

// adjust heap to select a proper hole to set v
func (h *timeHeap) adjustDown(b []*timeTask, hole int, v *timeTask) {
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
func (h *timeHeap) adjustUp(b []*timeTask, hole int, v *timeTask) {
	for hole > 0 {
		if parent := h.parent(hole); !h.cmp(v, b[parent]) {
			b[hole], hole = b[parent], parent
		} else {
			break
		}
	}
	b[hole] = v //put v to last hole
}

func (h *timeHeap) cmp(c, p *timeTask) bool {
	return p.when <= c.when
}

// get parent index
func (h *timeHeap) parent(idx int) int {
	return (idx - 1) / 2
}

// get left child index
func (h *timeHeap) lchild(idx int) int {
	return 2*idx + 1
}

func (h *timeHeap) Top() (int64, bool) {
	if h.Empty() {
		return -1, false
	}
	return h.b[0].when, true
}

func (h *timeHeap) Size() int {
	return len(h.b)
}

func (h *timeHeap) Empty() bool {
	return h.Size() == 0
}

// check if b is a valid heap
func (h *timeHeap) CheckHeap() bool {
	for i := h.Size() - 1; i > 0; i-- {
		p := h.parent(i)
		if !h.cmp(h.b[i], h.b[p]) {
			return false
		}
	}
	return true
}

func (h *timeHeap) Show() string {
	var b strings.Builder
	b.WriteString("-------\n")
	for i, v := range h.b {
		b.WriteString(fmt.Sprintf("%d ", v.when))
		if (i+1)&(i+2) == 0 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
