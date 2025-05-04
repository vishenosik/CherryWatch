package app

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vishenosik/CherryWatch/internal/services/models"
	"github.com/vishenosik/concurrency"
)

type Pool struct {
	pool    *concurrency.Pool
	subChan <-chan models.Task
}

func MustNewPool(subscriptions ...chan models.Task) *Pool {
	pool, err := NewPool(subscriptions...)
	if err != nil {
		panic(err)
	}
	return pool
}

func NewPool(subscriptions ...chan models.Task) (*Pool, error) {
	return NewPoolContext(context.Background(), subscriptions...)
}

func NewPoolContext(ctx context.Context, subscriptions ...chan models.Task) (*Pool, error) {
	if len(subscriptions) == 0 {
		return nil, errors.New("no subscriptions provided")
	}
	return &Pool{
		pool:    concurrency.NewWorkerPoolContext(ctx, concurrency.WithWorkersControl(3, 256, 3)),
		subChan: concurrency.MergeChannels(ctx, uint16(1024), subscriptions...),
	}, nil
}

func (p *Pool) Start(_ context.Context) {
	p.pool.Start()

	go func() {
		for task := range p.subChan {
			_, err := p.pool.AddTask(
				concurrency.Task{
					ID:       task.ID,
					Func:     task.Func,
					Priority: concurrency.Priority(task.Priority),
				},
			)
			if err != nil {
				if errors.Is(err, concurrency.ErrPoolClosed) {
					// TODO: extra handling
					// TODO: log error?
					return
				}
			}
		}
	}()

}

func (p *Pool) Stop() {
	p.pool.Stop()
}
