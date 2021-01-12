package timer

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TaskFunc only task function can be added to
// the timer.
type TaskFunc func(num int) (int, error)

// Timer is for processing task. it checks buckets
// for popping jobs. it will put ready jobs to queue.
type Timer interface {
	AddTask(taskFunc TaskFunc)
	Run()
	Close()
}

// timer is Timer implementation struct.
type timer struct {
	num   int            // number of tasks
	wg    sync.WaitGroup // wait group for quit
	tasks []taskStub     // task stub
	once  sync.Once      // once
}

// taskStub task stub for function itself and context,
// and cancel function for this task.
type taskStub struct {
	f      TaskFunc
	ctx    context.Context
	cancel context.CancelFunc
}

func New() Timer {
	return &timer{num: 20, wg: sync.WaitGroup{}}
}

func (t *timer) AddTask(taskFunc TaskFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	task := taskStub{f: taskFunc, ctx: ctx, cancel: cancelFunc}
	t.tasks = append(t.tasks, task)
}

// Run start all tasks, and wait all task is done
func (t *timer) Run() {
	t.wg.Add(len(t.tasks))

	for _, task := range t.tasks {
		go func(task taskStub) {
			defer t.wg.Done()
			task.run(t.num)
		}(task)
	}

	t.wg.Wait()
}

// Close call all task cancel function to stop all tasks
func (t *timer) Close() {
	t.once.Do(
		func() {
			for _, task := range t.tasks {
				task.cancel()
			}
		},
	)
}

// run a task, and wait for context is done.
// this can be implement with more thinking.
func (task taskStub) run(num int) {
	for {
		select {
		case <-task.ctx.Done():
			return
		default:
			processNum, err := task.f(num)
			if err != nil {
				// do something
				// TODO: add logger
				fmt.Println(err)
				time.Sleep(1 * time.Second)
				continue
			} else if processNum != num {
				// do something
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}
}
