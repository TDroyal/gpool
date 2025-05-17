package gpool

// cache worker and reuse goroutine
type workerQueue interface {
	len() int
	put(*worker)
	get() *worker
	getCleanList() []*worker
}
