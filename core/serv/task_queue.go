package serv

import (
	"fmt"
)

type TaskQueue struct {
	capacity int
	count    int
	pipe     chan Job
}

type Job interface {
	process(worker string)
}

func worker(pipe <-chan Job, name string) {
	for j := range pipe {
		j.process(name)
	}
}

type TaskQueueOption func(*TaskQueue)

func WithCount(count int) TaskQueueOption {
	return func(d *TaskQueue) {
		d.count = count
	}
}

func WithTaskQueueCapacity(capacity int) TaskQueueOption {
	return func(d *TaskQueue) {
		d.capacity = capacity
	}
}

func NewTaskQueue() *TaskQueue {
	t := TaskQueue{
		capacity: 10,
		count:    5,
	}
	t.pipe = make(chan Job, t.capacity)
	for i := 0; i < t.count; i++ {
		go worker(t.pipe, fmt.Sprintf("worker_%d", i))
	}
	return &t
}

func (t *TaskQueue) Add(j Job) {
	t.pipe <- j
}
