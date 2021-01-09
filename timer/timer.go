package timer

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TaskFunc func(num int) (int, error)

type Timer interface {
	AddTask(taskFunc TaskFunc)
	Run()
	Close()
}

type timer struct {
	num       int
	wg        sync.WaitGroup
	closeChan chan error
	tasks     []taskStub
	once      sync.Once
}

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

func (t *timer) Close() {
	t.once.Do(
		func() {
			for _, task := range t.tasks {
				task.cancel()
			}
		},
	)
}

func (task taskStub) run(num int) {
	for {
		select {
		case <-task.ctx.Done():
			return
		default:
			processNum, err := task.f(num)
			if err != nil {
				// do something
			} else if processNum == num {
				// do something
			}

			fmt.Println("process")
			time.Sleep(1 * time.Second)
		}
	}
}
