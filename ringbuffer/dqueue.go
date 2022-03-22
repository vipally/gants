package ringbuffer

import (
	"sync"
)

type Buffer struct {
	lock sync.Mutex
	head int
	tail int
	size int
}
