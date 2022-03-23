package gants

type goWorker struct {
	p  *Pool
	id int
}

func (w *goWorker) run() {
	var t *task
	for {
		t = nil
		select {
		case t = <-w.p.chTask: //fetch task from channel
			if t == nil { // channel closed
				return
			}
		default:
		}
		if t == nil {
			t, _ = w.p.tq.Pop() // fetch task from task queue
		}
		if t != nil {
			t.Execute() // execute the task
		} else {
			w.p.wCond.Wait() // no task, wait
		}
	}
}
