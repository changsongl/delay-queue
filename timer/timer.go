package timer

import (
	"context"
	"sync"
	"time"

	"github.com/changsongl/delay-queue/pkg/log"
)

// TaskFunc only task function can be added to
// the timer.
type TaskFunc func() (hasMore bool, err error)

// Timer is for processing task. it checks buckets
// for popping jobs. it will put ready jobs to queue.
type Timer interface {
	AddTask(taskFunc TaskFunc)
	Run()
	Close()
}

// timer is Timer implementation struct.
type timer struct {
	wg           sync.WaitGroup // wait group for quit
	tasks        []taskStub     // task stub
	once         sync.Once      // once
	l            log.Logger     // logger
	taskInterval time.Duration  // fetch interval
	taskDelay    time.Duration  // fetch delay when bucket has more jobs after a fetching. Default no wait.
}

// taskStub task stub for function itself and context,
// and cancel function for this task.
type taskStub struct {
	f      TaskFunc
	ctx    context.Context
	cancel context.CancelFunc
	l      log.Logger
}

// New create a new timer for loading ready jobs from bucket
func New(l log.Logger, taskInterval, taskDelay time.Duration) Timer {
	return &timer{
		wg:           sync.WaitGroup{},
		l:            l.WithModule("timer"),
		taskInterval: taskInterval,
		taskDelay:    taskDelay,
	}
}

// AddTask add task to timer
func (t *timer) AddTask(taskFunc TaskFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	task := taskStub{
		f:      taskFunc,
		ctx:    ctx,
		cancel: cancelFunc,
		l:      t.l,
	}
	t.tasks = append(t.tasks, task)
}

// Run start all tasks, and wait all task is done
func (t *timer) Run() {
	t.wg.Add(len(t.tasks))

	for _, task := range t.tasks {
		go func(task taskStub) {
			defer t.wg.Done()
			task.run(t.taskInterval, t.taskDelay)
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
func (task taskStub) run(fetchInterval, fetchDelay time.Duration) {
	for {
		select {
		case <-task.ctx.Done():
			return
		default:
			hasMore, err := task.f()
			if err != nil {
				task.l.Error("task.f task run failed", log.Error(err))
				time.Sleep(fetchInterval)
				continue
			} else if !hasMore {
				time.Sleep(fetchInterval)
				continue
			}

			// have more jobs, wait delay time to fetch next time
			time.Sleep(fetchDelay)
		}
	}
}
