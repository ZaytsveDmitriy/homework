package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkersAllow      = errors.New("allow workers qnt less 1")
)

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) < 1 {
		return nil
	}
	if n < 1 {
		return ErrNoWorkersAllow
	}
	if m < 1 {
		return fmt.Errorf("%w: less then 1 error allow", ErrErrorsLimitExceeded)
	}

	errCnt := atomic.Int64{}
	chClose := make(chan struct{})
	chTasks := make(chan Task)

	wg := sync.WaitGroup{}

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			worker(chTasks, &errCnt, chClose)
			wg.Done()
		}()
	}

	for _, task := range tasks {
		if int(errCnt.Load()) >= m {
			close(chClose)
			break
		}

		chTasks <- task
	}

	close(chTasks)
	wg.Wait()

	if int(errCnt.Load()) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(chTasks <-chan Task, errCnt *atomic.Int64, chClose <-chan struct{}) {
	for {
		task, hasTask := <-chTasks
		if !hasTask {
			return
		}

		select {
		case <-chClose:
			return
		default:
		}

		err := task()
		if err != nil {
			errCnt.Add(1)
		}
	}
}
