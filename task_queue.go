// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// License that can be found in the LICENSE file.

package gants

// taskQueue is lockfree queue of task
type taskQueue struct {
	d    []*task
	head int64
	tail int64
}

func (q *taskQueue) Init(size int) {
	q.grow(size)
}

func (q *taskQueue) Push(t *task) {
	q.d[int(q.tail)] = t
	if q.tail++; q.tail >= int64(q.Cap()) {
		q.tail = 0
		if q.tail == q.head { //tail catch up head, buffer full
			q.grow(0)
		}
	}
}

// PushFO push a task that with first out priority
func (q *taskQueue) PushFO(t *task) {
}

func (q *taskQueue) Pop() (t *task, ok bool) {
	if ok = q.head != q.tail; ok {
		t = q.d[q.head]
		if q.head++; q.head >= int64(q.Cap()) && q.head != q.tail {
			q.head = 0
		}
	}
	return
}

func (q *taskQueue) Cap() int {
	return len(q.d)
}

func (q *taskQueue) Len() int {
	return (int(q.tail) + q.Cap() - int(q.head)) % q.Cap()
}

func (q *taskQueue) grow(size int) {
	oldCap := q.Cap()
	d := q.d
	q.d = make([]*task, q.nextCap(oldCap, size))
	h := copy(q.d, d[q.head:])
	t := copy(q.d[:h], d[:q.tail])
	q.head, q.tail = 0, int64(h+t)
}

// nextCap get the next buffer size when grow
func (q *taskQueue) nextCap(oldCap int, newCap int) int {
	const qTooLarge = 4096
	switch {
	case newCap > oldCap:
		return newCap
	case oldCap < qTooLarge: // little size, 2*oldCap=>2^(n+1)
		newCap := oldCap * 2
		// if newCap!=2^x, thenewCap newCap=>2^(x+1), eg: 3=>6=>8
		// loop to remove the lowest binary digit 1, eg: 10110=>10100=>10000
		for t := 2 * (newCap & (newCap - 1)); t != 0; t &= (t - 1) {
			newCap = t
		}
		if newCap <= 0 {
			newCap = 16
		}
		return newCap
	default:
		// large size, grow by qTooLarge, at least +50%*qTooLarge
		const x = qTooLarge // 4096
		return ((oldCap+x/2)/x + 1) * x
	}
}
