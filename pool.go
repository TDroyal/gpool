package gpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidPoolCapacity = errors.New("the capacity of pool is invalid")
)

const (
	taskChanNumber  = 1
	defaultInterval = time.Second * 5
)

type Pool interface {
	Submit(Task) error
}

// goroutine pool with blocking function
type pool struct {
	capacity    int32         // the capacity of the pool
	running     int32         // the running number of workers
	workers     workerQueue   // idle workers   todo
	workerCache sync.Pool     // object pool of workers
	interval    time.Duration // the interval for cleaning up idle workers
	cond        *sync.Cond
	lock        sync.Locker
}

// new a pool with capacity
func NewPool(size int32, opts ...Option) (Pool, error) {
	if size < 1 {
		return nil, ErrInvalidPoolCapacity
	}

	gp := new(pool)

	// default config
	gp.capacity = size
	gp.workers = newWorkerList(size)
	gp.lock = &sync.Mutex{}
	gp.cond = sync.NewCond(gp.lock)
	gp.interval = defaultInterval

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

	// async monitor: periodically recycle idle workers
	go gp.cleanWorkerQueue(gp.interval)

	return gp, nil
}

func (gp *pool) cleanWorkerQueue(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		gp.lock.Lock()
		// get the worker
		clist := gp.workers.getCleanList()
		gp.lock.Unlock()
		// then stop the worker
		for i := range clist {
			clist[i].stop()
			clist[i] = nil // avoid memory leak
		}
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
	if w := gp.workers.get(); w != nil {
		gp.lock.Unlock()
		return w, nil
	}

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

func (gp *pool) retriveWorker(worker *worker) { // retrive worker to list
	gp.lock.Lock()
	gp.workers.put(worker)
	gp.lock.Unlock()

	gp.cond.Signal()
}

// atomic option
func (gp *pool) addRunning(delta int32) int32 {
	return atomic.AddInt32(&gp.running, delta)
}
