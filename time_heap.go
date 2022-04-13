package gants

import (
	"fmt"
	"strings"
	"time"
)

func nowTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func newTimeHeap(size int) *timeHeap {
	return &timeHeap{
		list: make([]*task, 0, size),
	}
}

type timeHeap struct {
	list []*task
}

func (h *timeHeap) Push(s *task) error {
	h.list = append(h.list, s)
	if err := h.adjustUp(h.list, len(h.list)-1, s); err != nil {
		h.list = h.list[:h.Size()-1]
		return err
	}
	return nil
}

func (h *timeHeap) PopIfTimeout(now int64) (*task, bool, time.Duration) {
	return nil, false, 0
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
		c := l                                            // index that need compare with hole
		if r := l + 1; r < len(b) && !h.cmp(b[r], b[l]) { // let the most proper child to compare with v
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
