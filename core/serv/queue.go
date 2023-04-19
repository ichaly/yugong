package serv

import (
	"context"
	"fmt"
	"go.uber.org/fx"
)

type Queue struct {
	capacity int
	count    int
	bus      chan Job
}

type Job interface {
	process(worker string)
}

type QueueOption func(*Queue)

func WithCount(count int) QueueOption {
	return func(d *Queue) {
		d.count = count
	}
}

func WithCapacity(capacity int) QueueOption {
	return func(d *Queue) {
		d.capacity = capacity
	}
}

func NewQueue(l fx.Lifecycle, opts ...QueueOption) *Queue {
	q := Queue{capacity: 10, count: 5}
	for _, o := range opts {
		o(&q)
	}
	q.bus = make(chan Job, q.capacity)
	l.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for i := 0; i < q.count; i++ {
				go q.worker(fmt.Sprintf("worker_%d", i))
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			q.Close()
			return nil
		},
	})
	return &q
}

func (my *Queue) worker(name string) {
	for j := range my.bus {
		j.process(name)
	}
}

func (my *Queue) Push(j Job) {
	go func() {
		my.bus <- j
	}()
}

func (my *Queue) Close() {
	close(my.bus)
}
