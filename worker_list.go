package gpool

// idle worker list
type workerList struct {
	wr []*worker
}

func newWorkerList(size int32) workerQueue {
	return &workerList{
		wr: make([]*worker, 0, size),
	}
}

func (l *workerList) getCleanList() []*worker {
	length := l.len()
	if length < 4 {
		return nil
	}

	length /= 4
	cworker := l.wr[:length]
	l.wr = l.wr[length:]

	return cworker
}

func (l *workerList) len() int {
	return len(l.wr)
}

func (l *workerList) put(w *worker) {
	l.wr = append(l.wr, w)
}

func (l *workerList) get() *worker {
	if l.len() == 0 {
		return nil
	}

	w := l.wr[0]
	l.wr[0] = nil // avoid memory leak
	l.wr = l.wr[1:]

	return w
}
