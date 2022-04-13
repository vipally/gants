package gants

type TaskID = uint64

type task struct {
	f           func()
	id          uint64
	executeTime int64
}

func (t *task) clean() *task {
	t.f = nil
	t.id = 0
	t.executeTime = 0
	return t
}
