package gants

import (
	"sync"
)

type Pool struct {
	lk   sync.Mutex
	cond *sync.Cond
}
