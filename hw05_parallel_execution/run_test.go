package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks with zerro allow errors", func(t *testing.T) {
		tasksCount := 2
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(10*(i+1)))

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				return nil
			})
		}

		workersCount := 2
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)

		require.ErrorIs(t, err, ErrErrorsLimitExceeded, "actual err - %v", err)
	})

	t.Run("tasks with zerro allow workers", func(t *testing.T) {
		tasksCount := 2
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(10*(i+1)))

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				return nil
			})
		}

		workersCount := 0
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)

		require.ErrorIs(t, err, ErrNoWorkersAllow, "actual err - %v", err)
	})

	t.Run("concurency with eventually", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			waitTime := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += waitTime

			tasks = append(tasks, func() error {
				t := time.NewTimer(waitTime)
				<-t.C
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Eventually(t, func() bool {
			return tasksCount == int(atomic.LoadInt32(&runTasksCount))
		},
			sumTime/2,
			10*time.Millisecond,
			"look like not concurency func")
	})
}
