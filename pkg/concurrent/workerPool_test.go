package concurrent

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewWorkerPool(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		pool := NewWorkerPool()
		if pool.minWorkers != 3 || pool.maxWorkers != 256 || pool.initWorkers != 10 {
			t.Errorf("unexpected default configuration")
		}
	})

	t.Run("custom configuration", func(t *testing.T) {
		pool := NewWorkerPool(WithWorkersControl(5, 20, 10))
		if pool.minWorkers != 5 || pool.maxWorkers != 20 || pool.initWorkers != 10 {
			t.Errorf("custom configuration not applied")
		}
	})

	t.Run("auto-adjust init workers", func(t *testing.T) {
		pool := NewWorkerPool(WithWorkersControl(5, 10, 3)) // init < min
		if pool.initWorkers != 5 {
			t.Errorf("expected init workers to be adjusted to min")
		}

		pool = NewWorkerPool(WithWorkersControl(5, 10, 15)) // init > max
		if pool.initWorkers != 10 {
			t.Errorf("expected init workers to be adjusted to max")
		}
	})
}

func TestWorkerPool_Lifecycle(t *testing.T) {
	pool := NewWorkerPool(
		WithWorkersControl(2, 5, 2),
		WithTimeouts(time.Millisecond*500, time.Millisecond*200),
	)
	pool.Start()

	time.Sleep(time.Second)

	if pool.numWorkers.Load() != 2 {
		t.Errorf("expected 2 workers at start, got %d", pool.numWorkers.Load())
	}

	// Test worker scaling
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := pool.AddTask(func() { time.Sleep(3 * time.Second) })
			log.Println(err)
		}()
	}
	wg.Wait()

	// if pool.numWorkers.Load() <= 2 {
	// 	t.Error("workers did not scale up under load")
	// }

	// Wait for scale down
	time.Sleep(200 * time.Millisecond)
	if pool.numWorkers.Load() != 2 {
		t.Errorf("expected workers to scale down to 2, got %d", pool.numWorkers.Load())
	}

	pool.Stop()

	// Verify pool is closed
	err := pool.AddTask(func() {})
	if !errors.Is(err, ErrPoolClosed) {
		t.Errorf("expected ErrPoolClosed, got %v", err)
	}
}

func TestWorkerPool_TaskExecution(t *testing.T) {
	pool := NewWorkerPool(WithWorkersControl(1, 1, 1))
	pool.Start()
	defer pool.Stop()

	var executed atomic.Bool
	err := pool.AddTask(func() { executed.Store(true) })
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	if !executed.Load() {
		t.Error("task was not executed")
	}
}

func TestWorkerPool_ConcurrentTaskSubmission(t *testing.T) {
	pool := NewWorkerPool(WithWorkersControl(5, 20, 5))
	pool.Start()
	defer pool.Stop()

	var counter atomic.Int32
	const numTasks = 1000
	var wg sync.WaitGroup

	wg.Add(numTasks)
	for i := 0; i < numTasks; i++ {
		go func() {
			defer wg.Done()
			err := pool.AddTask(func() { counter.Add(1) })
			if err != nil {
				t.Errorf("failed to add task: %v", err)
			}
		}()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond) // Allow all tasks to complete

	if counter.Load() != numTasks {
		t.Errorf("expected %d tasks executed, got %d", numTasks, counter.Load())
	}
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewWorkerPoolContext(ctx, WithWorkersControl(1, 1, 1))
	pool.Start()

	cancel()
	// time.Sleep(50 * time.Millisecond) // Allow cancellation to propagate

	err := pool.AddTask(func() {})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestWorkerPool_PanicRecovery(t *testing.T) {
	pool := NewWorkerPool(WithWorkersControl(1, 1, 1))
	pool.Start()
	defer pool.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	err := pool.AddTask(func() {
		defer wg.Done()
		panic("test panic")
	})
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	wg.Wait()
	// Test should not crash
}

func TestWorkerPool_CloseBehavior(t *testing.T) {
	pool := NewWorkerPool(WithWorkersControl(1, 1, 1))
	pool.Start()

	// Add task before closing
	err := pool.AddTask(func() {})
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	pool.Stop()

	// Verify pool is closed
	if !pool.closed.Load() {
		t.Error("pool not marked as closed")
	}

	// Verify channels are closed
	select {
	case _, ok := <-pool.taskCH:
		if ok {
			t.Error("task channel not closed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for task channel")
	}

	select {
	case _, ok := <-pool.closeCH:
		if ok {
			t.Error("close channel not closed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for close channel")
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	pool := NewWorkerPool(WithWorkersControl(10, 100, 10))
	pool.Start()
	defer pool.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.AddTask(func() {
				// Simulate work
				time.Sleep(10 * time.Microsecond)
			})
		}
	})
}

func BenchmarkWorkerPoolHeavyWork(b *testing.B) {
	pool := NewWorkerPool(WithWorkersControl(10, 100, 10))
	pool.Start()
	defer pool.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.AddTask(func() {
				// Simulate heavier work
				time.Sleep(1 * time.Millisecond)
			})
		}
	})
}

func BenchmarkWorkerPoolScaleUp(b *testing.B) {
	pool := NewWorkerPool(
		WithWorkersControl(1, 100, 1),
		WithTimeouts(10*time.Millisecond, 0),
	)
	pool.Start()
	defer pool.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pool.AddTask(func() {
			time.Sleep(100 * time.Microsecond)
		})
	}
}
