package gpool

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrInvalidPoolCapacity = errors.New("the capacity of pool is invalid")
)

const (
	taskChanNumber = 1
)

type Pool interface {
	Submit(Task) error
}

// goroutine pool with blocking function
type pool struct {
	capacity    int32     // the capacity of the pool
	running     int32     // the running number of workers
	workerCache sync.Pool // object pool of workers
	cond        *sync.Cond
	lock        sync.Locker
}

// new a pool with capacity
func NewPool(size int32, opts ...option) (Pool, error) {
	if size < 1 {
		return nil, ErrInvalidPoolCapacity
	}

	gp := new(pool)

	// default config
	gp.capacity = size
	gp.lock = &sync.Mutex{}
	gp.cond = sync.NewCond(gp.lock)

	// todo workerCache
	gp.workerCache.New = func() any {
		return &worker{
			gp:       gp,
			taskChan: make(chan Task, taskChanNumber),
		}
	}

	// personalization
	for _, opt := range opts {
		opt(gp)
	}

	return gp, nil
}

// option pattern
type option func(*pool)

// capacity
func WithCapacity(capacity int32) option {
	return func(gp *pool) {
		gp.capacity = capacity
	}
}

// submit task to pool
func (gp *pool) Submit(task Task) error {

	// retrive a worker
	w, err := gp.retrieveWorker()
	if err != nil {
		return err
	}

	// put task to workerchan
	w.inputTask(task)

	return nil
}

// retrieve a worker from pool
func (gp *pool) retrieveWorker() (*worker, error) {
	gp.lock.Lock()

retry:
	if gp.running < gp.capacity {
		gp.lock.Unlock()
		w := gp.workerCache.Get().(*worker) // sync.Pool is concurrent safe
		w.run()                             // start the worker
		return w, nil
	}

	// // block and wait for an available worker
	gp.cond.Wait()
	/*
		Wait is to unlock first, then wait to be awakened; after being awakened, perform the lock (lock grabbing) operation
	*/

	goto retry
}

func (gp *pool) retriveWorker(work *worker) { // retrive worker to object pool
	gp.workerCache.Put(work)
}

// atomic option
func (gp *pool) addRunning(delta int32) int32 {
	return atomic.AddInt32(&gp.running, delta)
}
