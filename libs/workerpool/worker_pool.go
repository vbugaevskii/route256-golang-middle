package workerpool

import (
	"context"
	"sync"
)

type TaskResult[I, O any] struct {
	Inp I
	Out O
	Err error
}

type WorkerPool[I, O any] struct {
	inp chan I
	out chan TaskResult[I, O]

	ctx context.Context
	wg  *sync.WaitGroup

	size int
	call func(i I) (O, error)
}

// create new worker pool
func NewWorkerPool[I, O any](ctx context.Context, size int, call func(i I) (O, error)) *WorkerPool[I, O] {
	// init pool
	wp := WorkerPool[I, O]{
		inp: make(chan I, size),
		out: make(chan TaskResult[I, O], size),

		ctx: ctx,
		wg:  &sync.WaitGroup{},

		size: size,
		call: call,
	}

	// create workers
	for i := 0; i < wp.size; i++ {
		wp.wg.Add(1)

		go func() {
			defer wp.wg.Done()

			for inp := range wp.inp {
				// try to execute the task or fall when pool is canceled
				select {
				case <-ctx.Done():
					return
				case wp.out <- wp.execTask(inp):

				}
			}
		}()
	}

	return &wp
}

// auxiliary function for executing the task
func (wp *WorkerPool[I, O]) execTask(inp I) TaskResult[I, O] {
	out, err := wp.call(inp)
	return TaskResult[I, O]{
		Inp: inp,
		Out: out,
		Err: err,
	}
}

// submit tasks to worker pool
func (wp *WorkerPool[I, O]) Submit(inp []I) <-chan TaskResult[I, O] {
	go func() {
		for _, i := range inp {
			wp.inp <- i
		}

		// no more tasks for pool, close input chanel
		close(wp.inp)
	}()

	return wp.out
}

// close worker pool
func (wp *WorkerPool[I, O]) Close() {
	// wait all the workers in group
	wp.wg.Wait()
	close(wp.out)

	// wait until reading is done
	<-wp.ctx.Done()

	// make pool impossible to use
	wp.inp = nil
	wp.out = nil
}
