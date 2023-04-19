package serv

import (
	"fmt"
)

type TaskQueue struct {
	capacity int
	count    int
	bus      chan Job
}

type Job interface {
	process(worker string)
}

type TaskQueueOption func(*TaskQueue)

func WithCount(count int) TaskQueueOption {
	return func(d *TaskQueue) {
		d.count = count
	}
}

func WithCapacity(capacity int) TaskQueueOption {
	return func(d *TaskQueue) {
		d.capacity = capacity
	}
}

func NewTaskQueue(opts ...TaskQueueOption) *TaskQueue {
	t := TaskQueue{capacity: 10, count: 5}
	for _, o := range opts {
		o(&t)
	}
	t.bus = make(chan Job, t.capacity)
	for i := 0; i < t.count; i++ {
		go t.worker(fmt.Sprintf("worker_%d", i))
	}
	return &t
}

func (my *TaskQueue) worker(name string) {
	for j := range my.bus {
		j.process(name)
	}
}

func (my *TaskQueue) Push(j Job) {
	go func() {
		my.bus <- j
	}()
}

func (my *TaskQueue) Close() {
	close(my.bus)
}
