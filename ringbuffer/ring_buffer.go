package ringbuffer

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var maxTry = 1
var sumTry = 0
var countTry = 0

// NewRingBuffer create a ringbuffer.
func NewRingBuffer(size int) *RingBuffer {
	p := &RingBuffer{}
	if err := p.Init(size); err != nil {
		panic(err)
	}
	return p
}

// BufferId is the id of a buffer.
type BufferId uint64

// RingBuffer is goroutine-safe cycle buffer.
// It is designed as busy share buffer with lots of readers and writers.
// RingBuffer must runs under parallelism mode(runtime.GOMAXPROCS >= 4).
// It enables enable real-parallel R/W on busy shared buffers.
// see:
//   http://ifeve.com/ringbuffer
//   http://mechanitis.blogspot.com/2011/06/dissecting-disruptor-whats-so-special.html
type RingBuffer struct {
	debug      bool
	totalWait  int64
	size       int32      // buffer size, readonly
	waitReadR  *sync.Cond // waitlist that are wating read
	waitWriteR *sync.Cond // waitlist that are wating write
	waitReadC  *sync.Cond // waitlist that are wating read commit
	waitWriteC *sync.Cond // waitlist that are wating write commit
	rReserve   uint64     // Read reserve, mutable
	rCommit    uint64     // Read commit, mutable
	wReserve   uint64     // Write reserve, mutable
	wCommit    uint64     // Write commit, mutable
}

// Debug switch debug options.
func (rb *RingBuffer) Debug(enable bool) {
	rb.debug = enable
}

func (rb *RingBuffer) log(wid int, name string) func() {
	start := time.Now()
	//fmt.Printf("%s %s start\n", start, name)
	deferFun := func() {
		end := time.Now()
		cost := end.Sub(start)
		if cost > time.Millisecond*200 {
			panic(cost.String())
		}
		totalCost := atomic.AddInt64(&rb.totalWait, int64(cost))
		fmt.Printf("%s %s end, cost=%s totalCost=%s\n", start, name, cost, time.Duration(totalCost))
	}
	return deferFun
}

// Init ringbuffer with size.
// It is not goroutine-safe.
// RingBuffer must runs under parallelism mode(runtime.GOMAXPROCS >= 4).
func (rb *RingBuffer) Init(size int) error {
	const needCpus = 4
	if runtime.GOMAXPROCS(0) < needCpus {
		return fmt.Errorf("RingBuffer: requires parallelism(runtime.GOMAXPROCS >= %d)", needCpus)
	}

	if rb.waitReadR == nil {
		rb.waitReadR = sync.NewCond(new(sync.Mutex))
		rb.waitWriteR = sync.NewCond(new(sync.Mutex))
		rb.waitReadC = sync.NewCond(new(sync.Mutex))
		rb.waitWriteC = sync.NewCond(new(sync.Mutex))
	}

	return rb.Resize(size)
}

// Resize change size of ringbuffer.
func (rb *RingBuffer) Resize(size int) error {
	if size <= 0 {
		return fmt.Errorf("RingBuffer: invalid size %d", size)
	}
	atomic.StoreInt32(&rb.size, int32(size))
	return nil
}

// Size return size of ringbuffer.
func (rb *RingBuffer) Size() int {
	return int(atomic.LoadInt32(&rb.size))
}

// BufferIndex returns logic index of buffer by id.
func (rb *RingBuffer) BufferIndex(id uint64) int {
	return int(id % uint64(rb.Size()))
}

// Show shows status of ring buffer
func (rb *RingBuffer) Show() string {
	return fmt.Sprintf("%s rR=%d rC=%d wR=%d wC=%d",
		time.Now().Format("2006-01-02T15:04:05.999"),
		atomic.LoadUint64(&rb.rReserve),
		atomic.LoadUint64(&rb.rCommit),
		atomic.LoadUint64(&rb.wReserve),
		atomic.LoadUint64(&rb.wCommit),
	)
}

// ReserveWrite returns next avable id for write.
// It will wait if ringbuffer is full.
// It is goroutine-safe.
func (rb *RingBuffer) ReserveWrite(wid int) (bid uint64) {
	bid = atomic.AddUint64(&rb.wReserve, 1) - 1

	if rb.debug {
		fn := rb.log(wid, "ReserveWrite")
		defer fn()
	}

	try := 0
	for {
		try++
		if rb.debug {
			logTry("ReserveWrite", try, wid, rb.Show())
		}

		dataStart := atomic.LoadUint64(&rb.rCommit)
		maxW := dataStart + uint64(rb.size)
		if bid < maxW { //no conflict, reserve ok
			break
		}
		runtime.Gosched()
		//buffer full, wait as writer in order to awake by another reader
		//rb.waitWriteR.L.Lock()
		//rb.waitWriteR.Wait()
		//rb.waitWriteR.L.Unlock()
	}
	updateTry(try)

	return
}

// CommitWrite commit writer event for id.
// It will wait if previous writer id havn't commit.
// It will awake on reader wait list after commit OK.
// It is goroutine-safe.
func (rb *RingBuffer) CommitWrite(wid int, bid uint64) {
	newId := bid + 1

	if rb.debug {
		fn := rb.log(wid, "CommitWrite")
		defer fn()
	}

	try := 0
	for {
		try++
		if rb.debug {
			logTry("CommitWrite", try, wid, rb.Show())
		}

		if atomic.CompareAndSwapUint64(&rb.wCommit, bid, newId) { //commit OK
			rb.waitReadR.Broadcast()  //wakeup reader
			rb.waitWriteC.Broadcast() //wakeup write committer
			break
		}

		runtime.Gosched()
		//commit fail, wait as reader in order to wakeup by another writer
		//rb.waitWriteC.L.Lock()
		//rb.waitWriteC.Wait()
		//rb.waitWriteC.L.Unlock()
	}
	updateTry(try)
}

// ReserveRead returns next avable id for read.
// It will wait if ringbuffer is empty.
// It is goroutine-safe.
func (rb *RingBuffer) ReserveRead(wid int) (bid uint64) {
	bid = atomic.AddUint64(&rb.rReserve, 1) - 1

	if rb.debug {
		fn := rb.log(wid, "ReserveRead")
		defer fn()
	}

	try := 0
	for {
		try++
		if rb.debug {
			logTry("ReserveRead", try, wid, rb.Show())
		}

		w := atomic.LoadUint64(&rb.wCommit)
		if bid < w { //no conflict, reserve ok
			break
		}

		runtime.Gosched()
		//buffer empty, wait as reader in order to wakeup by another writer
		//rb.waitReadR.L.Lock()
		//rb.waitReadR.Wait()
		//rb.waitReadR.L.Unlock()
	}
	updateTry(try)

	return
}

// CommitRead commit reader event for id.
// It will wait if previous reader id havn't commit.
// It will awake on writer wait list after commit OK.
// It is goroutine-safe.
func (rb *RingBuffer) CommitRead(wid int, bid uint64) {
	newId := bid + 1

	if rb.debug {
		fn := rb.log(wid, "CommitRead")
		defer fn()
	}

	try := 0
	for {
		try++
		if rb.debug {
			logTry("CommitRead", try, wid, rb.Show())
		}

		if atomic.CompareAndSwapUint64(&rb.rCommit, bid, newId) {
			rb.waitWriteR.Broadcast() //wakeup writer
			rb.waitReadC.Broadcast()  //wakeup read committer
			break
		}

		runtime.Gosched()
		//commit fail, wait as writer in order to wakeup by another reader
		//rb.waitReadC.L.Lock()
		//rb.waitReadC.Wait()
		//rb.waitReadC.L.Unlock()
	}
	updateTry(try)
}

func updateTry(try int) {
	if sumTry > 0 {
		countTry++
	}
	sumTry += try
	if try > maxTry {
		maxTry = try
		fmt.Printf("************maxTry=%d\n", maxTry)
	}
}

func logTry(name string, try int, wid int, msg string) {
	if try > 1 {
		fmt.Printf("%s try=%d wid=%d maxTry=%d avgTry=%d %s\n", name, try, wid, maxTry, sumTry/countTry, msg)
		if try > 100 {
			//panic(try)
		}
	}
}
