package gpool

// task
type Task func()

// worker
type worker struct {
	gp       *pool
	taskChan chan Task
}

func (w *worker) run() {
	w.gp.addRunning(1)
	go func() {
		defer func() {
			w.gp.addRunning(-1)
			// if there is a task waiting for an free worker, wake it up
			w.gp.workerCache.Put(w) // put worker to object pool
			w.gp.cond.Signal()
		}()

		for task := range w.taskChan {
			if task == nil {
				return
			}

			task()
			// recycle the worker to list
			w.gp.retriveWorker(w)
		}
	}()
}

func (w *worker) inputTask(task Task) {
	w.taskChan <- task
}

// stop and recycle the worker
func (w *worker) stop() {
	w.taskChan <- nil
}
