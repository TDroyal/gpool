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
			w.gp.cond.Signal()
		}()

		for task := range w.taskChan {
			if task == nil {
				return
			}

			task()
			// recycle the worker
			w.gp.retriveWorker(w)
			return
		}
	}()
}

func (w *worker) inputTask(task Task) {
	w.taskChan <- task
}
