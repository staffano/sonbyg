package builder

import (
	"context"
	"sync"
)

// Worker can be assigned by the environment to work on a target
type Worker struct {
	free              bool
	workingOn         *Task
	workerDoneChannel chan bool
	ctxt              context.Context
	lock              sync.Mutex
	utilization       int64
}

// Free reports if the Worker is free to take up new works
func (w *Worker) Free() bool { return w.free }

// WorkingOn reports what the Worker currently is working on, if anything...
func (w *Worker) WorkingOn() *Task { return w.workingOn }

// Work tells the worker to start working on making the target.
func (w *Worker) Work(b *Builder, t *Task) {
	w.lock.Lock()
	w.workingOn = t
	w.free = false
	w.lock.Unlock()
	go func() {
		// We assume the target is attached to the same context as
		// we, so that a termination of the context will stop the
		// Make command and hence make us leave the go routine.
		b.Execute(w.workingOn)
		w.lock.Lock()
		w.utilization = w.utilization + 1
		w.workingOn = nil
		w.free = true
		w.lock.Unlock()
		if !t.RunAlways {
			b.Stamp(t)
		}
		w.workerDoneChannel <- true
	}()
}

// NewWorker creates a new worker in a context and atttached to
// the channel it should report to when work is done.
func NewWorker(ctxt context.Context, wdchan chan bool) *Worker {
	w := Worker{workerDoneChannel: wdchan, ctxt: ctxt, free: true}
	return &w
}
