package ringbuffer

import (
	"sync"
	"testing"
	"time"
)

func TestRingBuffer(t *testing.T) {
	b := NewRingBuffer(10)
	b.Debug(true)
	w := &worker{
		b: b,
	}
	for i := 0; i < 20; i++ {
		go w.run(i + 1)
	}
	w.wg.Wait()
}

type worker struct {
	b  *RingBuffer
	wg sync.WaitGroup
}

func (w *worker) run(id int) {
	w.wg.Add(1)
	defer w.wg.Done()
	for i := 0; i < 5; i++ {
		bid := w.b.ReserveWrite(id)
		time.Sleep(time.Millisecond)
		w.b.CommitWrite(id, bid)
		time.Sleep(time.Microsecond)

		bid2 := w.b.ReserveRead(id)
		time.Sleep(time.Millisecond)
		w.b.CommitRead(id, bid2)
		time.Sleep(time.Microsecond)
	}
}
