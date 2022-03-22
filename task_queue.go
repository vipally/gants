package gants

// taskQueue is lockfree queue of task
type taskQueue struct {
}

func (q taskQueue) Push(t *Task) error {
	return nil
}

func (q taskQueue) Pop() (*Task, bool) {
	return nil, false
}

func (q taskQueue) Len() int {
	return 0
}
